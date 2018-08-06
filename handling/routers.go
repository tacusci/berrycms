package handling

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/gobuffalo/plush"
	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

type MutableRouter struct {
	mu   sync.Mutex
	Root *mux.Router
}

func (mr *MutableRouter) Swap(root *mux.Router) {
	mr.mu.Lock()
	mr.Root = root
	mr.mu.Unlock()
}

func (mr *MutableRouter) Reload() {
	r := mux.NewRouter()

	loginHandler := &LoginHandler{}
	r.HandleFunc("/admin", loginHandler.Handle)
	usersHandler := &UsersHandler{}
	r.HandleFunc("/admin/users", usersHandler.Handle)
	pagesHandler := &PagesHandler{}
	r.HandleFunc("/admin/pages", pagesHandler.Handle)

	mr.Swap(r)
}

type Handler struct{}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {}

type LoginHandler struct{}

func (lh *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.html")
	if err != nil {
		logging.Error("Unable to find resources folder...")
		w.Write([]byte("<h1>500 Server Error</h1>"))
	}
	w.Write(content)
}

type UsersHandler struct{}

func (uh *UsersHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

type PagesHandler struct{}

func (ph *PagesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	pageroutes := make([]string, 0)

	pt := db.PagesTable{}
	row, err := pt.Select(db.Conn, "route", "")

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for row.Next() {
		p := &db.Page{}
		row.Scan(&p.Route)
		pageroutes = append(pageroutes, p.Route)
	}

	pctx := plush.NewContext()
	pctx.Set("names", pageroutes)

	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.pages.html")
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
