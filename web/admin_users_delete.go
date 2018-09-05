package web

import (
	"fmt"
	"net/http"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminUsersDeleteHandler handler to contain pointer to core router and the URI string
type AdminUsersDeleteHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (audh *AdminUsersDeleteHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (audh *AdminUsersDeleteHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/admin/users", http.StatusFound)

	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		Error(w, err)
	}

	ut := db.UsersTable{}
	st := db.AuthSessionsTable{}

	for _, v := range r.PostForm {
		userToDelete, err := ut.SelectByUUID(db.Conn, v[0])

		if err != nil {
			logging.Error(err.Error())
			Error(w, err)
		}

		//don't allow deletion of the root user account
		if db.UsersRoleFlag(userToDelete.UserroleId) != db.ROOT_USER {
			amw := AuthMiddleware{}
			loggedInUser, err := amw.LoggedInUser(r)

			if err != nil {
				logging.Error(err.Error())
				Error(w, err)
			}

			//make sure that the logged in user is not the same as user to delete
			if loggedInUser.UUID != userToDelete.UUID {
				st.Delete(db.Conn, fmt.Sprintf("uuid = '%s'", userToDelete.UUID))
				ut.DeleteByUUID(db.Conn, userToDelete.UUID)
			}
		}
	}
}

//Route get URI route for handler
func (audh *AdminUsersDeleteHandler) Route() string { return audh.route }

//HandlesGet retrieve whether this handler handles get requests
func (audh *AdminUsersDeleteHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (audh *AdminUsersDeleteHandler) HandlesPost() bool { return true }
