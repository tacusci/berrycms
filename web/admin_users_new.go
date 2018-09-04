package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/util"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/logging"
)

//AdminUsersNewHandler handler to contain pointer to core router and the URI string
type AdminUsersNewHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (aunh *AdminUsersNewHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "New User")
	RenderDefault(w, "admin.users.new.html", pctx)
}

//Post handles post requests to URI
func (aunh *AdminUsersNewHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/admin/users", http.StatusFound)
	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		http.Redirect(w, r, "/admin/users", http.StatusFound)
		return
	}

	authHash := r.PostFormValue("authhash")
	repeatedAuthHash := r.PostFormValue("repeatedauthhash")

	if strings.Compare(authHash, repeatedAuthHash) == 0 {
		ut := db.UsersTable{}
		userToCreate := db.User{
			Username:        r.PostFormValue("username"),
			CreatedDateTime: time.Now().Unix(),
			Email:           r.PostFormValue("email"),
			UserroleId:      int(db.REG_USER),
			FirstName:       r.PostFormValue("firstname"),
			LastName:        r.PostFormValue("lastname"),
			AuthHash:        util.HashAndSalt([]byte(authHash)),
		}

		err := ut.Insert(db.Conn, userToCreate)

		if err != nil {
			logging.Error(err.Error())
		}
	} else {
		//need to add setting error message on screen
	}
}

//Route get URI route for handler
func (aunh *AdminUsersNewHandler) Route() string { return aunh.route }

//HandlesGet retrieve whether this handler handles get requests
func (aunh *AdminUsersNewHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (aunh *AdminUsersNewHandler) HandlesPost() bool { return true }
