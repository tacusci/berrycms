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
	"github.com/tacusci/berrycms/handling"
	"github.com/tacusci/logging"
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

	logging.InfoNnl("Connecting to mysql:localhost/berrycms schema...")

	db.Connect("mysql", "berryadmin:Password12345@/", "berrycms")
	db.Wipe()
	db.Setup()
	db.CreateTestData()
	go db.Heartbeat()

	rs := &handling.MutableRouter{}
	rs.Reload()

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      rs.Root,
	}

	go listenForStopSig(srv)
	srv.ListenAndServe()
}

func listenForStopSig(srv *http.Server) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	logging.Error(fmt.Sprintf("Caught sig: %+v (Shutting down and cleaning up...)", sig))
	logging.Info("Closing DB connection...")
	db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	logging.Info("Stopping HTTP server...")
	srv.Shutdown(ctx)
	logging.Info("Shutting down... BYE!")
	os.Exit(0)
}
