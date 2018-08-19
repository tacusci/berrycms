package web

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminUsersHandler contains response functions for users admin page
type AdminUsersHandler struct {
	Router *MutableRouter
	route  string
}

//Get takes the web request and writes response to session
func (uh *AdminUsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	usernames := make([]string, 0)

	ut := db.UsersTable{}
	row, err := ut.Select(db.Conn, "createddatetime, username", "")

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for row.Next() {
		u := &db.User{}
		row.Scan(&u.CreatedDateTime, &u.Username)
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

func (uh *AdminUsersHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (uh *AdminUsersHandler) Route() string { return uh.route }

func (uh *AdminUsersHandler) HandlesGet() bool  { return true }
func (uh *AdminUsersHandler) HandlesPost() bool { return false }
