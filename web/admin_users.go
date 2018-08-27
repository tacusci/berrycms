package web

import (
	"net/http"

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
	pctx.Set("title", "Users")
	pctx.Set("unixtostring", UnixToTimeString)

	RenderDefault(w, "admin.users.html", pctx)
}

func (uh *AdminUsersHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (uh *AdminUsersHandler) Route() string { return uh.route }

func (uh *AdminUsersHandler) HandlesGet() bool  { return true }
func (uh *AdminUsersHandler) HandlesPost() bool { return false }
