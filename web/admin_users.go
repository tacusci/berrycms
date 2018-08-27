package web

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

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
	users := make([]db.User, 0)

	ut := db.UsersTable{}
	rows, err := ut.Select(db.Conn, "createddatetime, firstname, lastname, username, email", "")
	defer rows.Close()

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for rows.Next() {
		u := db.User{}
		rows.Scan(&u.CreatedDateTime, &u.FirstName, &u.LastName, &u.Username, &u.Email)
		users = append(users, u)
	}

	pctx := plush.NewContext()
	pctx.Set("users", users)
	pctx.Set("unixtostring", func(unix int64) string {
		return time.Unix(unix, 0).Format("15:04:05 02-01-2006")
	})

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
