package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
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
	gt := db.GroupTable{}
	amw := AuthMiddleware{}

	loggedInUser, err := amw.LoggedInUser(r)

	if loggedInUser != nil {

		rows, err := gt.Select(db.Conn, "uuid, title", fmt.Sprintf("uuid = '%s'", groupUUID))
		if err != nil {
			Error(w, err)
			return
		}

		var groupToRemoveFrom *db.Group

		if rows.Next() {
			groupToRemoveFrom = &db.Group{}
			rows.Scan(&groupToRemoveFrom.UUID, &groupToRemoveFrom.Title)
		}

		rows.Close()

		if groupToRemoveFrom == nil {
			Error(w, fmt.Errorf("Unable to read group of UUID %s from database", groupUUID))
			return
		}

		for _, v := range r.PostForm {
			userToRemove, err := ut.SelectByUUID(db.Conn, v[0])
			if err != nil {
				logging.Error(err.Error())
				continue
			}

			//don't allow the root user to be removed from the admins user group
			if userToRemove == nil || (groupToRemoveFrom.Title == "Admins" && userToRemove.UserroleId == int(db.ROOT_USER)) {
				continue
			}

			gmt.DeleteUserFromGroup(db.Conn, userToRemove, groupToRemoveFrom)
		}
	}
}

//Route get URI route for handler
func (augerh *AdminUserGroupsEditRemoveHandler) Route() string { return augerh.route }

//HandlesGet retrieve whether this handler handles get requests
func (augerh *AdminUserGroupsEditRemoveHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (augerh *AdminUserGroupsEditRemoveHandler) HandlesPost() bool { return true }
