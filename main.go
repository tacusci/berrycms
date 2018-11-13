// Copyright (c) 2018, tacusci ltd
//
// Licensed under the GNU GENERAL PUBLIC LICENSE Version 3 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/gpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	sqlUsername = flag.String("dbuser", "berryadmin", "Database server username, ignored if using sqlite")
	sqlPassword = flag.String("dbpass", "", "Database server password, ignored if using sqlite")
	sqlAddress = flag.String("dbaddr", "/", "Database server location, ignored if using sqlite")

	flag.Parse()

	loggingLevel := logging.InfoLevel
	logging.ColorLogLevelLabelOnly = true

	if *debugLevel {
		logging.SetLevel(logging.DebugLevel)
		return
	}
	logging.SetLevel(loggingLevel)
}

func main() {
	parseCmdArgs()

	fmt.Printf("🍓 Berry CMS %s 🍓\n", VERSION)

	switch *sql {
	case "sqlite":
		db.Connect(db.SQLITE, "", "berrycms")
	case "mysql":
		db.Connect(db.MySQL, fmt.Sprintf("%s:%s@%s", *sqlUsername, *sqlPassword, *sqlAddress), "berrycms")
	default:
		logging.ErrorAndExit(fmt.Sprintf("Unknown database server type %s...", *sql))
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
