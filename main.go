package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/web"
	"github.com/tacusci/logging"
)

const (
	VERSION = "v0.0.1a"
)

var (
	devMode     *bool
	port        *int
	addr        *string
	sql         *string
	sqlUsername *string
	sqlPassword *string
	sqlAddress  *string
)

func parseCmdArgs() {
	debugLevel := flag.Bool("dbg", false, "Set logging to debug")
	devMode = flag.Bool("dev", false, "Turn on development mode")
	port = flag.Int("p", 8080, "Port to listen for HTTP requests on")
	addr = flag.String("a", "0.0.0.0", "IP address to listen against if multiple network adapters")
	sql = flag.String("db", "sqlite", "Database server type to try to connect to [sqlite/mysql]")
	sqlUsername = flag.String("dbu", "berryadmin", "Database server username, ignored if using sqlite")
	sqlPassword = flag.String("dbp", "", "Database server password, ignored if using sqlite")
	sqlAddress = flag.String("dbaddr", "/", "Database server location, ignored if using sqlite")

	flag.Parse()

	loggingLevel := logging.InfoLevel

	if *debugLevel {
		logging.SetLevel(logging.DebugLevel)
		return
	}
	logging.SetLevel(loggingLevel)
}

func main() {
	parseCmdArgs()

	fmt.Printf("🍓 Berry CMS %s 🍓\n", VERSION)

	if *sql == "sqlite" {
		db.Connect(db.SQLITE, "", "berrycms")
	} else if *sql == "mysql" {
		db.Connect(db.MySQL, fmt.Sprintf("%s:%s@%s", *sqlUsername, *sqlPassword, *sqlAddress), "berrycms")
	}

	if *devMode {
		db.Wipe()
	}

	db.Setup()

	if *devMode {
		db.CreateTestData()
	}

	go db.Heartbeat()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", *addr, *port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	rs := web.MutableRouter{
		Server: srv,
	}
	rs.Reload()

	clearOldSessionsStop := make(chan bool)

	go web.ClearOldSessions(&clearOldSessionsStop)
	go listenForStopSig(srv, &clearOldSessionsStop)

	logging.Info(fmt.Sprintf("Starting http server @ %s 🌏 ...", srv.Addr))
	err := srv.ListenAndServe()

	if err != nil {
		logging.ErrorAndExit(fmt.Sprintf("☠️  Error starting server (%s) ☠️", err.Error()))
	}
}

func listenForStopSig(srv *http.Server, wc *chan bool) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	logging.Debug("Stopping clearing old sessions...")
	//send a terminate command to the session clearing goroutine's channel
	*wc <- true
	logging.Error(fmt.Sprintf("☠️  Caught sig: %+v (Shutting down and cleaning up...) ☠️", sig))
	logging.Info("Closing DB connection...")
	db.Close()
	logging.Info("Stopping HTTP server...")
	logging.Info("Shutting down... BYE! 👋")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}
