package web

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
	"net/http"
)

//AdminUserGroupsEditAddHandler contains response functions for pages admin page
type AdminUserGroupsEditAddHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Post(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	groupUUID := vars["uuid"]

	if groupUUID == "" {
		Error(w, errors.New("Missing group UUID"))
		return
	}

	var redirectURI = "/admin/users/groups/edit/" + groupUUID

	if augeah.Router.AdminHidden {
		redirectURI = fmt.Sprintf("/%s", augeah.Router.AdminHiddenPassword) + redirectURI
	}

	defer http.Redirect(w, r, redirectURI, http.StatusFound)

	err := r.ParseForm()

	if err != nil {
		Error(w, err)
		return
	}

	ut := db.UsersTable{}
	gmt := db.GroupMembershipTable{}
	gt := db.GroupTable{}
	amw := AuthMiddleware{}

	loggedInUser, err := amw.LoggedInUser(r)

	if loggedInUser != nil {

		groupTitle := ""
		rows, err := gt.Select(db.Conn, "title", fmt.Sprintf("uuid = '%s'", groupUUID))
		if err != nil {
			Error(w, err)
			return
		}

		for rows.Next() {
			rows.Scan(&groupTitle)
		}

		for _, v := range r.PostForm {
			userToAdd, err := ut.SelectByUUID(db.Conn, v[0])
			if err != nil {
				logging.Error(err.Error())
				continue
			}

			if userToAdd == nil {
				continue
			}

			gmt.AddUserToGroup(db.Conn, userToAdd, groupTitle)
		}
	}
}

//Route get URI route for handler
func (augeah *AdminUserGroupsEditAddHandler) Route() string { return augeah.route }

//HandlesGet retrieve whether this handler handles get requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesPost() bool { return true }
