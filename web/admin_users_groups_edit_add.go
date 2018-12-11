package web

import (
	"fmt"
	"net/http"

	"github.com/tacusci/logging"

	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
)

//AdminUserGroupsEditAddHandler contains response functions for pages admin page
type AdminUserGroupsEditAddHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Get(w http.ResponseWriter, r *http.Request) {
	var usersInGroup = make([]db.User, 0)
	vars := mux.Vars(r)
	groupUUID := vars["uuid"]

	gt := db.GroupMembershipTable{}
	groupMembershipRows, err := gt.Select(db.Conn, "useruuid", fmt.Sprintf("groupuuid = '%s'", groupUUID))

	if err != nil {
		Error(w, err)
		return
	}

	defer groupMembershipRows.Close()

	for groupMembershipRows.Next() {
		gm := db.GroupMembership{}
		err := groupMembershipRows.Scan(&gm.UserUUID)

		if err != nil {
			Error(w, err)
			continue
		}

		ut := db.UsersTable{}
		userRows, err := ut.Select(db.Conn, "createddatetime, userroleid, uuid, username, authhash, firstname, lastname, email", fmt.Sprintf("uuid = '%s'", gm.UserUUID))

		if err != nil {
			Error(w, err)
			continue
		}

		for userRows.Next() {
			user := db.User{}
			err = userRows.Scan(&user.CreatedDateTime, &user.UserroleId, &user.UUID, &user.Username, &user.AuthHash, &user.FirstName, &user.LastName, &user.Email)

			if err != nil {
				Error(w, err)
				continue
			}
			usersInGroup = append(usersInGroup, user)
		}

		//close at end of for
		userRows.Close()
	}

	logging.Debug(fmt.Sprintf("Found %d users in group -> %s", len(usersInGroup), groupUUID))
}

//Post handles post requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Post(w http.ResponseWriter, r *http.Request) {

}

//Route get URI route for handler
func (augeah *AdminUserGroupsEditAddHandler) Route() string { return augeah.route }

//HandlesGet retrieve whether this handler handles get requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesPost() bool { return true }
