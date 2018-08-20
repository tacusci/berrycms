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

func setLoggingLevel() {
	debugLevel := flag.Bool("d", false, "Set logging to debug")
	flag.Parse()

	loggingLevel := logging.InfoLevel

	if *debugLevel {
		logging.SetLevel(logging.DebugLevel)
		return
	}
	logging.SetLevel(loggingLevel)
}

func main() {
	setLoggingLevel()

	fmt.Printf("🍓 Berry CMS %s 🍓\n", VERSION)

	logging.InfoNnl("Connecting to mysql:localhost/berrycms schema...")

	db.Connect("mysql", "berryadmin:Password12345@/", "berrycms")
	db.Wipe()
	db.Setup()
	db.CreateTestData()
	go db.Heartbeat()

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
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
