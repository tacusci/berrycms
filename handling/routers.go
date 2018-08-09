package handling

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gobuffalo/plush"
	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/util"
	"github.com/tacusci/logging"
)

//MutableRouter is a mutex lock for the mux router
type MutableRouter struct {
	Server *http.Server
	mu     sync.Mutex
	Root   *mux.Router
	dw     *util.RecursiveDirWatch
}

//Swap takes a new mux router, locks accessing for old one, replaces it and then unlocks, keeps existing connections
func (mr *MutableRouter) Swap(root *mux.Router) {
	mr.mu.Lock()
	mr.Root = root
	mr.Server.Handler = mr.Root
	mr.mu.Unlock()
}

//Reload map all admin/default page routes and load saved page routes from DB
func (mr *MutableRouter) Reload() {
	r := mux.NewRouter()

	logging.Debug("Mapping default admin routes...")

	loginHandler := &LoginHandler{Router: mr}
	r.HandleFunc("/login", loginHandler.Get).Methods("GET")
	r.HandleFunc("/login", loginHandler.Post).Methods("POST")
	adminHandler := &AdminHandler{Router: mr}
	r.HandleFunc("/admin", adminHandler.Get).Methods("GET")
	usersHandler := &UsersHandler{Router: mr}
	r.HandleFunc("/admin/users", usersHandler.Get).Methods("GET")
	pagesHandler := &PagesHandler{Router: mr}
	r.HandleFunc("/admin/pages", pagesHandler.Get).Methods("GET")

	mr.mapSavedPageRoutes(r)
	mr.mapStaticDir(r, "static")
	go mr.monitorStatic("static")

	amw := authMiddleware{}
	amw.Populate()

	r.Use(amw.Middleware)

	mr.Swap(r)
}

func (mr *MutableRouter) mapSavedPageRoutes(r *mux.Router) {
	savedPageHandler := &SavedPageHandler{Router: mr}

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "route", "")
	if err != nil {
		logging.Error(err.Error())
		return
	}
	for rows.Next() {
		p := db.Page{}
		rows.Scan(&p.Route)
		logging.Debug(fmt.Sprintf("Mapping database page route %s", p.Route))
		r.HandleFunc(p.Route, savedPageHandler.Get)
	}
}

func (mr *MutableRouter) mapStaticDir(r *mux.Router, sd string) {
	fs, err := ioutil.ReadDir(sd)
	if err != nil {
		logging.Error("Unable to find static folder...")
		return
	}
	for _, f := range fs {
		pathPrefixLocation := fmt.Sprintf("%s%s%s", string(os.PathSeparator), f.Name(), string(os.PathSeparator))
		pathPrefixAddress := fmt.Sprintf("/%s/", f.Name())
		logging.Debug(fmt.Sprintf("Serving dir (%s)'s files...", f.Name()))
		r.PathPrefix(pathPrefixAddress).Handler(http.StripPrefix(pathPrefixAddress, http.FileServer(http.Dir(sd+pathPrefixLocation))))
	}
}

func (mr *MutableRouter) monitorStatic(sd string) {
	mr.dw = &util.RecursiveDirWatch{Change: make(chan bool), Stop: make(chan bool)}
	go mr.dw.WatchDir(sd)
	for {
		if <-mr.dw.Change {
			logging.Debug("Change detected in static dir...")
			break
		}
	}
	mr.dw.Stop <- true
	mr.Reload()
}

type authMiddleware struct {
	userTokens map[string]string
}

func (amw *authMiddleware) Populate() {
	amw.userTokens = make(map[string]string)
	amw.userTokens["000000"] = "root"
}

func (amw *authMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if amw.HasPermissions(w, r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func (amw *authMiddleware) HasPermissions(w http.ResponseWriter, r *http.Request) bool {
	sessionToken := r.Header.Get("X-Session-Token")
	if strings.HasPrefix(r.RequestURI, "/admin") {
		if _, found := amw.userTokens[sessionToken]; found {
			return true
		}
	} else {
		pt := db.PagesTable{}
		requestedPage, err := pt.SelectByRoute(db.Conn, r.RequestURI)
		if err != nil {
			return true
		}
		if requestedPage.Roleprotected {
			if _, found := amw.userTokens[sessionToken]; found {
				return true
			} else {
				return false
			}
		} else {
			return true
		}
	}
	return false
}

//LoginHandler contains response functions for admin login
type LoginHandler struct {
	Router *MutableRouter
}

//Get takes the web request and writes response to session
func (lh *LoginHandler) Get(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "login.html")
	if err != nil {
		logging.Error("Unable to find resources folder...")
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	pctx := plush.NewContext()
	pctx.Set("formname", "loginform")
	pctx.Set("formhash", "fjei4ijiorgrig4ijio34ofj4ig34i4j3i")

	renderedContent, err := plush.Render(string(content), pctx)
	if err != nil {
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write([]byte(renderedContent))
}

func (lh *LoginHandler) Post(w http.ResponseWriter, r *http.Request) {
	logging.Info("Recieved post request for login...")
	r.Method = "GET"
	lh.Get(w, r)
}

type AdminHandler struct {
	Router *MutableRouter
}

func (ah *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.html")
	if err != nil {
		logging.Error("Unable to find resources folder...")
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write(content)
}

//UsersHandler contains response functions for users admin page
type UsersHandler struct {
	Router *MutableRouter
}

//Get takes the web request and writes response to session
func (uh *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	renderedContent, err := plush.Render(string(content), pctx)
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write([]byte(renderedContent))
}

//PagesHandler contains response functions for pages admin page
type PagesHandler struct {
	Router *MutableRouter
}

//Get takes the web request and writes response to session
func (ph *PagesHandler) Get(w http.ResponseWriter, r *http.Request) {
	pageroutes := make([]string, 0)

	pt := db.PagesTable{}
	row, err := pt.Select(db.Conn, "route", "")

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for row.Next() {
		p := db.Page{}
		row.Scan(&p.Route)
		pageroutes = append(pageroutes, p.Route)
	}

	pctx := plush.NewContext()
	pctx.Set("names", pageroutes)

	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.pages.html")
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	renderedContent, err := plush.Render(string(content), pctx)
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write([]byte(renderedContent))
}

type SavedPageHandler struct {
	Router *MutableRouter
}

func (sph *SavedPageHandler) Get(w http.ResponseWriter, r *http.Request) {
	pt := db.PagesTable{}
	//JUST FOR LIVE/HOT ROUTE REMAPPING TESTING
	if r.RequestURI == "/addnew" {
		pt.Insert(db.Conn, db.Page{
			Title:         "Carbon",
			Route:         "/carbonite",
			Content:       "<h2>Carbonite</h2>",
			Roleprotected: true,
		})
		sph.Router.Reload()
	}
	row, err := pt.Select(db.Conn, "content", fmt.Sprintf("route = '%s'", r.RequestURI))
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	p := db.Page{}
	for row.Next() {
		row.Scan(&p.Content)
	}

	ctx := plush.NewContext()
	html, err := plush.Render(p.Content, ctx)
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write([]byte(html))
}
