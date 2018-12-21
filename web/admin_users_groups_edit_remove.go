package web

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
	"net/http"
)

//AdminUserGroupsEditRemoveHandler contains response functions for pages admin page
type AdminUserGroupsEditRemoveHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (augerh *AdminUserGroupsEditRemoveHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (augerh *AdminUserGroupsEditRemoveHandler) Post(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	groupUUID := vars["uuid"]

	if groupUUID == "" {
		Error(w, errors.New("Missing group UUID"))
		return
	}

	var redirectURI = "/admin/users/groups/edit/" + groupUUID

	if augerh.Router.AdminHidden {
		redirectURI = fmt.Sprintf("/%s", augerh.Router.AdminHiddenPassword) + redirectURI
	}

	defer http.Redirect(w, r, redirectURI, http.StatusFound)

	err := r.ParseForm()

	if err != nil {
		Error(w, err)
		return
	}

	ut := db.UsersTable{}
	gmt := db.GroupMembershipTable{}
	amw := AuthMiddleware{}

	loggedInUser, err := amw.LoggedInUser(r)

	if loggedInUser != nil {
		for _, v := range r.PostForm {
			userToRemove, err := ut.SelectByUUID(db.Conn, v[0])
			if err != nil {
				logging.Error(err.Error())
				continue
			}

			if userToRemove == nil {
				continue
			}

			gmt.DeleteUserFromGroup(db.Conn, userToRemove, &db.Group{UUID: groupUUID})
		}
	}
}

//Route get URI route for handler
func (augerh *AdminUserGroupsEditRemoveHandler) Route() string { return augerh.route }

//HandlesGet retrieve whether this handler handles get requests
func (augerh *AdminUserGroupsEditRemoveHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (augerh *AdminUserGroupsEditRemoveHandler) HandlesPost() bool { return true }
