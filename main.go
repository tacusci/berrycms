package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gobuffalo/plush"

	"github.com/gorilla/mux"

	"github.com/tacusci/berrycms/db"
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
	db.Connect("mysql", "berryadmin:Password12345@/", "berrycms")
	db.Wipe()
	db.Setup()
	db.CreateTestData()
	go db.Heartbeat()

	r := mux.NewRouter()
	r.HandleFunc("/admin", AdminHandler)
	r.HandleFunc("/admin/users", AdminUsersHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go listenForStopSig(srv)
	srv.ListenAndServe()
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.html")
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
	}
	w.Write(content)
	// w.WriteString([]byte()))
}

func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {

	usernames := make([]string, 0)

	ut := db.UsersTable{}
	row, err := ut.Select(db.Conn, "username", "")

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for row.Next() {
		u := &db.User{}
		row.Scan(&u.Username)
		usernames = append(usernames, u.Username)
	}

	pctx := plush.NewContext()
	pctx.Set("names", usernames)

	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.users.html")
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
	}
	renderedContent, err := plush.Render(string(content), pctx)
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
	}
	w.Write([]byte(renderedContent))
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
