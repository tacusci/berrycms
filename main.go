// Copyright (c) 2019 tacusci ltd
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
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/web"
	"github.com/tacusci/logging"
)

type options struct {
	testData            bool
	wipe                bool
	yesToAll            bool
	port                uint
	addr                string
	sql                 string
	sqlUsername         string
	sqlPassword         string
	sqlAddress          string
	activityLogLoc      string
	adminHiddenPassword string
	noRobots            bool
	logFileName         string
}

var shuttingDown bool

func parseCmdArgs() *options {
	opts := &options{}

	debugLevel := flag.Bool("dbg", false, "Set logging to debug")
	flag.BoolVar(&opts.testData, "testdb", false, "Creates testing data")
	flag.BoolVar(&opts.wipe, "wipe", false, "Completely wipes database")
	flag.BoolVar(&opts.yesToAll, "y", false, "Automatically agree to cli confirmation requests")
	flag.UintVar(&opts.port, "p", 8080, "Port to listen for HTTP requests on")
	flag.StringVar(&opts.addr, "a", "0.0.0.0", "IP address to listen against if multiple network adapters")
	flag.StringVar(&opts.sql, "db", "sqlite", "Database server type to try to connect to [sqlite/mysql]")
	flag.StringVar(&opts.sqlUsername, "dbuser", "berryadmin", "Database server username, ignored if using sqlite")
	flag.StringVar(&opts.sqlPassword, "dbpass", "", "Database server password, ignored if using sqlite")
	flag.StringVar(&opts.sqlAddress, "dbaddr", "/", "Database server location, ignored if using sqlite")
	flag.StringVar(&opts.activityLogLoc, "actlog", "", "Activity/access log file location")
	flag.StringVar(&opts.adminHiddenPassword, "ahp", "", "URI prefix to hide admin pages behind")
	flag.BoolVar(&opts.noRobots, "nrtxt", false, "Don't provide a robots.txt URI")
	flag.StringVar(&opts.logFileName, "log", "", "Server log file location")

	flag.Parse()

	loggingLevel := logging.WarnLevel
	logging.ColorLogLevelLabelOnly = true

	if *debugLevel {
		logging.SetLevel(logging.DebugLevel)
		return opts
	}
	logging.SetLevel(loggingLevel)

	return opts
}

func main() {
	opts := parseCmdArgs()

	flushInitialised := make(chan bool)
	if len(opts.logFileName) > 0 {
		go logging.FlushLogs(opts.logFileName, &flushInitialised)
		//halt main thread until creating file to flush logs to has initialised
		<-flushInitialised
	}

	logging.WhiteOutput(fmt.Sprintf("🍓 Berry CMS %s 🍓\n", db.VERSION))

	switch opts.sql {
	case "sqlite":
		db.Connect(db.SQLITE, "", "berrycms")
	case "mysql":
		db.Connect(db.MySQL, fmt.Sprintf("%s:%s@%s", opts.sqlUsername, opts.sqlPassword, opts.sqlAddress), "berrycms")
	default:
		logging.ErrorAndExit(fmt.Sprintf("Unknown database server type %s...", opts.sql))
	}

	var wipeOccurred bool

	if opts.wipe {
		//if yes to all flag is was used, user won't be prompted
		//to confirm as if statement won't continue evaluating
		//conditions eg., the bool val returned by 'askConfirmToWipe'
		//so it'll never be called
		if opts.yesToAll || askConfirmToWipe() {
			db.Wipe()
			wipeOccurred = true
		} else {
			logging.Info("Skipping wiping database...")
		}
	}

	db.Setup()

	//if wipe never happened but test data creation requested, display message/warning
	if !wipeOccurred && opts.testData {
		logging.Warn("Wipe not carried out, skipping creating test data...")
	}

	if wipeOccurred {
		if opts.testData {
			db.CreateTestData()
		}
	}

	go db.Heartbeat()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", opts.addr, opts.port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	rs := web.MutableRouter{
		Server:              srv,
		ActivityLogLoc:      opts.activityLogLoc,
		AdminHidden:         len(opts.adminHiddenPassword) > 0,
		AdminHiddenPassword: opts.adminHiddenPassword,
		NoRobots:            opts.noRobots,
	}
	rs.Reload()

	clearOldSessionsStop := make(chan bool)

	go web.ClearOldSessions(&clearOldSessionsStop)
	go listenForStopSig(srv, &clearOldSessionsStop)

	logging.Info(fmt.Sprintf("Starting http server @ %s 🌏 ...", srv.Addr))
	err := srv.ListenAndServe()

	if !shuttingDown {
		if err != nil {
			//only bother outputting error returned from listening server if we're not already trying to shutdown
			logging.ErrorAndExit(fmt.Sprintf("☠️  Error starting server (%s) ☠️", err.Error()))
		}
	}

	logging.Info("Closing DB connection...")
	db.Close()

	logging.Info("Shutting down... BYE! 👋")

	//stop writing log lines to file
	if logging.LoggingOutputReciever != nil {
		close(logging.LoggingOutputReciever)
	}
	close(flushInitialised)
}

func askConfirmToWipe() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		logging.YellowOutput("⚠ Wiping the database is irreversible, are you sure? ⚠ ")
		fmt.Printf(" [y/n]: ")

		response, err := reader.ReadString('\n')

		if err != nil {
			logging.ErrorAndExit(err.Error())
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
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
	shuttingDown = true
	logging.Error(fmt.Sprintf("☠️ Caught sig: %+v (Shutting down and cleaning up...) ☠️", sig))
	logging.Info("Stopping HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
}
