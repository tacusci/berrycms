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

package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/plugins"
	"github.com/tacusci/berrycms/robots"
	"github.com/tacusci/berrycms/util"
	"github.com/tacusci/logging"
)

//MutableRouter is a mutex lock for the mux router
type MutableRouter struct {
	Server              *http.Server
	mu                  sync.Mutex
	Root                *mux.Router
	AdminOff            bool
	AdminHidden         bool
	AdminHiddenPassword string
	ActivityLogLoc      string
	NoRobots            bool
	NoSitemap           bool
	CpuProfile          bool
	staticwatcher       *watcher.Watcher
	pluginswatcher      *watcher.Watcher
	pm                  *plugins.Manager
}

//Swap takes a new mux router, locks accessing for old one, replaces it and then unlocks, keeps existing connections
func (mr *MutableRouter) Swap(root *mux.Router) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.Root = root
	mr.Server.Handler = mr.Root
}

//Reload map all admin/default page routes and load saved page routes from DB
func (mr *MutableRouter) Reload() {

	if !mr.NoRobots {
		//creates a robot string and loads into in-memory cache
		err := robots.Generate(mr.AdminOff)
		if err != nil {
			logging.Error(err.Error())
		}
	}

	if mr.staticwatcher != nil {
		mr.staticwatcher.Close()
	}
	mr.staticwatcher = watcher.New()

	if mr.pluginswatcher != nil {
		mr.pluginswatcher.Close()
	}
	mr.pluginswatcher = watcher.New()

	r := mux.NewRouter()

	if mr.CpuProfile {
		logging.Warn("Mapping CPU profiling URIs...")
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
		r.Handle("/debug/pprof/block", pprof.Handler("block"))
		r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
		r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	}

	if !mr.AdminOff {
		logging.Debug("Mapping default admin routes...")

		for _, handler := range GetDefaultHandlers(mr) {
			if handler.HandlesGet() {
				logging.Debug(fmt.Sprintf("Mapping default GET route %s", handler.Route()))
				r.HandleFunc(handler.Route(), handler.Get).Methods("GET")
			}

			if handler.HandlesPost() {
				logging.Debug(fmt.Sprintf("Mapping default POST route %s", handler.Route()))
				r.HandleFunc(handler.Route(), handler.Post).Methods("POST")
			}
		}

		ut := db.UsersTable{}
		if !ut.RootUserExists() {
			aunh := AdminUsersNewHandler{
				Router: mr,
			}
			//add explicit mapping of root user creation handler routes
			r.HandleFunc("/admin/users/root/new", aunh.Get).Methods("GET")
			r.HandleFunc("/admin/users/root/new", aunh.Post).Methods("POST")
		}
	} else {
		logging.Warn("ADMIN PAGES HAVE BEEN DISABLED!")
	}

	//do really need to map this even if the robots.txt cache hasn't been initialised
	robotsHandler := &RobotsHandler{
		route:  "/robots.txt",
		Router: mr,
	}

	logging.Debug(fmt.Sprintf("Mapping default GET route %s", robotsHandler.Route()))
	r.HandleFunc(robotsHandler.Route(), robotsHandler.Get).Methods("GET")

	//do really need to map this even if the sitemap.xml cache hasn't been initialised
	sitemapHandler := &SitemapHandler{
		route:  "/sitemap.xml",
		Router: mr,
	}

	logging.Debug(fmt.Sprintf("Mapping default GET route %s", sitemapHandler.Route()))
	r.HandleFunc(sitemapHandler.Route(), sitemapHandler.Get).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(fourOhFour)

	mr.mapSavedPageRoutes(r)

	pm := plugins.NewManager()

	if err := pm.Load(); err != nil {
		logging.Error(err.Error())
	}

	pm.Lock()
	for _, plugin := range *pm.Plugins() {
		if val, err := plugin.VM.Get("routesToRegister"); err == nil {
			//extract list value of 'routesToRegister' from plugin
			if valInterface, err := val.Export(); err == nil {
				//try and convert list to slice of strings exclusively, if this fails don't continue
				if routesToRegister, ok := valInterface.([]string); ok {
					routesToRegister = util.RemoveDuplicates(routesToRegister)
					//for each route/value from list, map the route to the db_pages handler
					for _, value := range routesToRegister {
						logging.Debug(fmt.Sprintf("PLUGIN {%s} -> Mapping page route '%s'", plugin.UUID(), value))
						mr.mapPluginCreatedRoute(r, value)
					}
				}
			}
		}
		if _, err := plugin.Call("main", nil, nil); err != nil {
			logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", plugin.UUID(), err.Error()))
		}
	}
	pm.Unlock()

	if err := mr.mapStaticDir(r, "static"); err == nil {
		go mr.monitorStatic("./static", mr.staticwatcher)
	} else {
		logging.Error("Unable to map static dir, not listening for directory changes...")
	}
	go mr.monitorPlugins("./plugins", mr.pluginswatcher)

	alm := ActivityLogMiddleware{
		Router: mr,
		LogLoc: mr.ActivityLogLoc,
	}
	r.Use(alm.Middleware)

	am := AuthMiddleware{Router: mr}
	r.Use(am.Middleware)

	mr.Swap(r)
}

func (mr *MutableRouter) mapSavedPageRoutes(r *mux.Router) {
	savedPageHandler := &SavedPageHandler{Router: mr}

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "route", "")
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		logging.Error(err.Error())
		return
	}
	for rows.Next() {
		p := db.Page{}
		rows.Scan(&p.Route)
		logging.Debug(fmt.Sprintf("Mapping database page route %s", p.Route))
		r.HandleFunc(p.Route, savedPageHandler.Get).Methods("GET")
		r.HandleFunc(p.Route, savedPageHandler.Post).Methods("POST")
	}
}

func (mr *MutableRouter) mapPluginCreatedRoute(r *mux.Router, route string) {
	r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ctx := plush.NewContext()
		ctx.Set("pagecontent", "")
		Render(w, r, &db.Page{Route: route, Content: ""}, ctx)
	}).Methods("GET")
}

func (mr *MutableRouter) mapStaticDir(r *mux.Router, sd string) error {
	fs, err := ioutil.ReadDir(sd)
	if os.IsNotExist(err) {
		logging.Error("Unable to find static folder... Creating...")
		err = os.Mkdir("./static", os.ModeDir)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	for _, f := range fs {
		pathPrefixLocation := fmt.Sprintf("%s%s%s", string(os.PathSeparator), f.Name(), string(os.PathSeparator))
		pathPrefixAddress := fmt.Sprintf("/%s/", f.Name())
		logging.Debug(fmt.Sprintf("Serving dir (%s)'s files...", f.Name()))
		r.PathPrefix(pathPrefixAddress).Handler(http.StripPrefix(pathPrefixAddress, http.FileServer(http.Dir(sd+pathPrefixLocation))))
	}
	return nil
}

func (mr *MutableRouter) monitorStatic(sd string, w *watcher.Watcher) {
	w.SetMaxEvents(1)

	w.FilterOps(watcher.Create, watcher.Remove, watcher.Rename)

	go func() {
		for {
			select {
			case <-w.Event:
				mr.Reload()
			case err := <-w.Error:
				logging.Error(err.Error())
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(sd); err != nil {
		logging.Error(err.Error())
	}

	if err := w.Start(time.Millisecond * 500); err != nil {
		logging.Error(err.Error())
	}
}

func (mr *MutableRouter) monitorPlugins(sd string, w *watcher.Watcher) {
	w.SetMaxEvents(1)

	w.FilterOps(watcher.Create, watcher.Remove, watcher.Rename)

	go func() {
		for {
			select {
			case <-w.Event:
				mr.pm = plugins.NewManager()
				mr.pm.Load()
			case err := <-w.Error:
				logging.Error(err.Error())
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(sd); err != nil {
		logging.Error(err.Error())
	}

	if err := w.Start(time.Millisecond * 500); err != nil {
		logging.Error(err.Error())
	}
}

type ActivityLogMiddleware struct {
	Router *MutableRouter
	LogLoc string
}

func (alm *ActivityLogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if f, err := os.OpenFile(alm.LogLoc, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660); err == nil {
			defer f.Close()
			var sb bytes.Buffer
			sb.WriteString(fmt.Sprintf("%s (%s -> %s %s)", logging.GetTimeString(), r.RemoteAddr, r.Method, r.RequestURI))
			userAgent := r.UserAgent()
			if len(userAgent) > 0 {
				sb.WriteString(fmt.Sprintf(" [User Agent: '%s']", userAgent))
			}
			sb.WriteString("\n")
			f.Write(sb.Bytes())
		}
		next.ServeHTTP(w, r)
	})
}

//AuthMiddleware authentication struct with auth helper functions
type AuthMiddleware struct {
	Router *MutableRouter
}

//Middleware attaches http handler to middleware
func (amw *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if amw.HasPermissionsForRoute(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Access denied", http.StatusForbidden)
		}
	})
}

//HasPermissionsForRoute checks that requesting client has permissions to access the requested URI
func (amw *AuthMiddleware) HasPermissionsForRoute(r *http.Request) bool {
	var routeIsProtected bool

	var hiddenAdminPrefixURI = ""

	if amw.Router.AdminHidden {
		hiddenAdminPrefixURI = fmt.Sprintf("/%s", amw.Router.AdminHiddenPassword)
	}

	routeIsProtected = strings.HasPrefix(r.RequestURI, hiddenAdminPrefixURI+"/admin")

	if routeIsProtected && strings.Compare(r.RequestURI, "/admin/users/root/new") == 0 {
		ut := db.UsersTable{}
		if !ut.RootUserExists() {
			routeIsProtected = false
		} else {
			//if this route is mapped then within the existing server runtime there was no root user
			//however if there now is one we want to unmap this route as well as preventing access
			amw.Router.Reload()
		}
	}

	if !routeIsProtected {
		pt := db.PagesTable{}
		page, err := pt.SelectByRoute(db.Conn, r.RequestURI)
		if err == nil {
			routeIsProtected = page.Roleprotected
		}
	}

	if routeIsProtected {
		return amw.IsLoggedIn(r)
	}
	return true
}

//IsLoggedIn checks if the requesting client is currently logged in
func (amw *AuthMiddleware) IsLoggedIn(r *http.Request) bool {
	var isLoggedIn bool

	authSessionStore, err := sessionsstore.Get(r, "auth")
	if err == nil {
		if authSessionUUID := authSessionStore.Values["sessionuuid"]; authSessionUUID != nil {
			authSessionsTable := db.AuthSessionsTable{}
			authSession, err := authSessionsTable.SelectBySessionUUID(db.Conn, authSessionUUID.(string))
			if err == nil {
				if len(authSession.UserUUID) > 0 {
					isLoggedIn = true
					authSession.LastActiveDateTime = time.Now().Unix()
					authSessionsTable.Update(db.Conn, authSession)
				} else {
					authSessionsTable.DeleteBySessionUUID(db.Conn, authSessionUUID.(string))
					authSessionStore.Options.MaxAge = -1
				}
			} else {
				errString := err.Error()
				if errString != "sql: no rows in result set" {
					logging.Error(err.Error())
				}
			}
		}
	} else {
		logging.Debug(fmt.Sprintf("Error trying to read existing session \"auth\" -> %s", err.Error()))
	}
	return isLoggedIn
}

//LoggedInUser get user of existing web session
func (amw *AuthMiddleware) LoggedInUser(r *http.Request) (*db.User, error) {
	authSessionStore, err := sessionsstore.Get(r, "auth")
	if err == nil {
		authSessionsTable := db.AuthSessionsTable{}
		if authSessionUUID := authSessionStore.Values["sessionuuid"]; authSessionUUID != nil {
			authSession, err := authSessionsTable.SelectBySessionUUID(db.Conn, authSessionUUID.(string))
			if err == nil {
				if len(authSession.UserUUID) > 0 {
					ut := db.UsersTable{}
					loggedInUser, err := ut.SelectByUUID(db.Conn, authSession.UserUUID)
					if err != nil {
						return nil, err
					}
					return loggedInUser, nil
				}
			}
		}
	}
	return nil, err
}

//Error writes HTTP error message to web response and add error message to log
func Error(w http.ResponseWriter, err error) {
	logging.Error(err.Error())

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "content", fmt.Sprintf("route = '%s'", "[500]"))

	if err != nil {
		//potential stack overflow, should change this
		Error(w, err)
	}

	defer rows.Close()

	p := &db.Page{}
	// set intial content here just in case no custom page exists
	p.Content = "<h1>500 Server Error</h1>"

	for rows.Next() {
		rows.Scan(&p.Content)
	}

	ctx := plush.NewContext()
	ctx.Set("pagecontent", template.HTML(p.Content))

	WriteHTMLAndStatus(w, RenderStr(ctx), http.StatusInternalServerError)
}

func renderFourOhFour() (*plush.Context, error) {
	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "content", fmt.Sprintf("route = '%s'", "[404]"))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	p := &db.Page{}
	p.Content = "<h1>404 page not found</h1>"

	for rows.Next() {
		rows.Scan(&p.Content)
	}

	ctx := plush.NewContext()
	ctx.Set("pagecontent", template.HTML(p.Content))

	return ctx, nil
}

func fourOhFour(w http.ResponseWriter, r *http.Request) {
	ctx, err := renderFourOhFour()
	if err != nil {
		Error(w, err)
		return
	}
	WriteHTMLAndStatus(w, RenderStr(ctx), http.StatusNotFound)
}

func WriteHTMLAndStatus(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
