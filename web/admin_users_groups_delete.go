package web

import (
	"fmt"
	"net/http"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminUserGroupsDeleteHandler handler to contain pointer to core router and the URI string
type AdminUserGroupsDeleteHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (augdh *AdminUserGroupsDeleteHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (augdh *AdminUserGroupsDeleteHandler) Post(w http.ResponseWriter, r *http.Request) {

	var redirectURI = "/admin/users/groups"

	if augdh.Router.AdminHidden {
		redirectURI = fmt.Sprintf("/%s", augdh.Router.AdminHiddenPassword) + redirectURI
	}

	defer http.Redirect(w, r, redirectURI, http.StatusFound)

	err := r.ParseForm()

	if err != nil {
		Error(w, err)
	}

	gt := db.GroupTable{}
	amw := AuthMiddleware{}

	loggedInUser, err := amw.LoggedInUser(r)

	for _, v := range r.PostForm {
		groupToDelete, err := gt.SelectByUUID(db.Conn, v[0])

		if err != nil {
			logging.Error(err.Error())
			continue
		}

		if groupToDelete == nil {
			continue
		}

		if groupToDelete.Title != "Admins" && groupToDelete.Title != "Moderators" && groupToDelete.Title != "Users" {
			if loggedInUser != nil {
				gt.DeleteByUUID(db.Conn, groupToDelete.UUID)
			}
		}
	}
}

//Route get URI route for handler
func (augdh *AdminUserGroupsDeleteHandler) Route() string { return augdh.route }

//HandlesGet retrieve whether this handler handles get requests
func (augdh *AdminUserGroupsDeleteHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (augdh *AdminUserGroupsDeleteHandler) HandlesPost() bool { return true }
