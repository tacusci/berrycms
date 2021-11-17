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

package cmd

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

	"golang.org/x/crypto/acme/autocert"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/web"
	"github.com/tacusci/berrycms/web/config"
	"github.com/tacusci/logging"
)

var shuttingDown bool

func parseCmdArgs() config.Options {
	opts := config.Options{}

	debugLevel := flag.Bool("dbg", false, "Set logging to debug")
	flag.BoolVar(&opts.TestData, "testdb", false, "Creates testing data")
	flag.BoolVar(&opts.Wipe, "wipe", false, "Completely wipes database")
	flag.BoolVar(&opts.YesToAll, "y", false, "Automatically agree to cli confirmation requests")
	flag.UintVar(&opts.Port, "p", 8080, "Port to listen for HTTP requests on")
	flag.StringVar(&opts.Addr, "a", "0.0.0.0", "IP address to listen against if multiple network adapters")
	flag.StringVar(&opts.Sql, "db", "sqlite", "Database server type to try to connect to [sqlite/mysql]")
	flag.StringVar(&opts.SqlUsername, "dbuser", "berryadmin", "Database server username, ignored if using sqlite")
	flag.StringVar(&opts.SqlPassword, "dbpass", "", "Database server password, ignored if using sqlite")
	flag.StringVar(&opts.SqlAddress, "dbaddr", "/", "Database server location, ignored if using sqlite")
	flag.StringVar(&opts.ActivityLogLoc, "actlog", "", "Activity/access log file location")
	flag.StringVar(&opts.AdminHiddenPassword, "ahp", "", "URI prefix to hide admin pages behind")
	flag.BoolVar(&opts.NoRobots, "nrtxt", false, "Don't provide a robots.txt URI")
	flag.BoolVar(&opts.NoSitemap, "nsxml", false, "Don't provide a sitemap.xml URI")
	flag.BoolVar(&opts.AdminPagesDisabled, "apd", false, "Admin interface pages disabled")
	flag.StringVar(&opts.LogFileName, "log", "", "Server log file location")
	flag.BoolVar(&opts.CpuProfile, "cpuprofile", false, "Enable CPU profiling")
	flag.StringVar(&opts.AutoCertDomain, "autocert", "", "Domain/web address to serve HTTPS against")

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
	if len(opts.LogFileName) > 0 {
		go logging.FlushLogs(opts.LogFileName, &flushInitialised)
		//halt main thread until creating file to flush logs to has initialised
		<-flushInitialised
	}

	logging.WhiteOutput(fmt.Sprintf("🍓 Berry CMS %s 🍓\n", db.VERSION))

	switch opts.Sql {
	case "sqlite":
		db.Connect(db.SQLITE, "", "berrycms")
	case "mysql":
		db.Connect(db.MySQL, fmt.Sprintf("%s:%s@%s", opts.SqlUsername, opts.SqlPassword, opts.SqlAddress), "berrycms")
	default:
		logging.ErrorAndExit(fmt.Sprintf("Unknown database server type %s...", opts.Sql))
	}

	var wipeOccurred bool

	if opts.Wipe {
		//if yes to all flag is was used, user won't be prompted
		//to confirm as if statement won't continue evaluating
		//conditions eg., the bool val returned by 'askConfirmToWipe'
		//so it'll never be called
		if opts.YesToAll || askConfirmToWipe() {
			db.Wipe()
			wipeOccurred = true
		} else {
			logging.Info("Skipping wiping database...")
		}
	}

	db.Setup()

	//if wipe never happened but test data creation requested, display message/warning
	if !wipeOccurred && opts.TestData {
		logging.Warn("Wipe not carried out, skipping creating test data...")
	}

	if wipeOccurred {
		if opts.TestData {
			db.CreateTestData()
		}
	}

	go db.Heartbeat()

	var certManager *autocert.Manager

	if opts.AutoCertDomain != "" {
		certManager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(opts.AutoCertDomain),
			Cache:      autocert.DirCache(cacheDir(opts.AutoCertDomain)),
		}
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", opts.Addr, opts.Port),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 120,
	}

	if certManager != nil {
		srv.Addr = ":https"
		srv.TLSConfig = certManager.TLSConfig()
	}

	rs := web.MutableRouter{
		Server:              srv,
		ActivityLogLoc:      opts.ActivityLogLoc,
		AdminOff:            opts.AdminPagesDisabled,
		AdminHidden:         len(opts.AdminHiddenPassword) > 0,
		AdminHiddenPassword: opts.AdminHiddenPassword,
		NoRobots:            opts.NoRobots,
		NoSitemap:           opts.NoSitemap,
		CpuProfile:          opts.CpuProfile,
	}
	rs.Reload()

	clearOldSessionsStop := make(chan bool)

	go web.ClearOldSessions(&clearOldSessionsStop)
	go listenForStopSig(srv, &clearOldSessionsStop)

	logging.Info(fmt.Sprintf("Starting http server @ %s 🌏 ...", srv.Addr))

	var err error
	if srv.TLSConfig == nil {
		err = srv.ListenAndServe()
	} else {
		err = srv.ListenAndServeTLS("", "")
	}

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

func cacheDir(domain string) (dir string) {
	if domain != "" {
		dir = fmt.Sprintf("%s%scache-autocert-%s", os.TempDir(), string(os.PathSeparator), domain)
		logging.Info(fmt.Sprintf("Saving acquired SSL cert to cache: %s", dir))
		var err error
		if err = os.MkdirAll(dir, 0700); err == nil {
			return dir
		}
		logging.Error(fmt.Sprintf("Error creating SSL cert cache folder: %s", err.Error()))
	}
	return ""
}

//fires on Ctrl+C/SIGTERM send to process
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
