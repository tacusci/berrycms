package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gobuffalo/plush"
	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

type AdminUserGroupsEditHandler struct {
	Router *MutableRouter
	route  string
}

func (augeh *AdminUserGroupsEditHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gmt := db.GroupMembershipTable{}

	//retrieve every membership for this group
	groupMembershipRows, err := gmt.Select(db.Conn, "createddatetime, groupuuid, useruuid", fmt.Sprintf("groupuuid = '%s'", vars["uuid"]))
	if err != nil {
		Error(w, errors.New("Group memberships not found"))
		return
	}

	usersInGroup := make([]db.User, 0)

	//read each membership into struct
	for groupMembershipRows.Next() {
		groupMembership := db.GroupMembership{}
		err := groupMembershipRows.Scan(&groupMembership.CreatedDateTime, &groupMembership.GroupUUID, &groupMembership.UserUUID)
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		ut := db.UsersTable{}

		//read each user of each membership into struct and add to users in group list
		user, err := ut.SelectByUUID(db.Conn, groupMembership.UserUUID)
		if err != nil {
			logging.Error(err.Error())
			continue
		}
		usersInGroup = append(usersInGroup, *user)
	}

	groupMembershipRows.Close()

	//retrieve list of all existing users
	ut := db.UsersTable{}
	userRows, err := ut.Select(db.Conn, "*", "")
	if err != nil {
		Error(w, err)
		return
	}

	usersNotInGroup := make([]db.User, 0)

	u := &db.User{}

	//read each existing user into struct and add to users not in group list if not already in the in group list
	for userRows.Next() {
		err = userRows.Scan(&u.UserId, &u.CreatedDateTime, &u.UserroleId, &u.UUID, &u.Username, &u.AuthHash, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		userInGroup := false
		for _, user := range usersInGroup {
			if u.UUID == user.UUID {
				userInGroup = true
				break
			}
		}

		if !userInGroup {
			usersNotInGroup = append(usersNotInGroup, *u)
		}
	}

	userRows.Close()

	gt := db.GroupTable{}
	groupRows, err := gt.Select(db.Conn, "title", fmt.Sprintf("uuid = '%s'", vars["uuid"]))
	if err != nil {
		Error(w, err)
		return
	}

	var groupTitle string
	if groupRows.Next() {
		groupRows.Scan(&groupTitle)
	}

	groupRows.Close()

	pctx := plush.NewContext()
	pctx.Set("title", "Edit Group")
	pctx.Set("submitroute", r.RequestURI)
	pctx.Set("grouptitle", groupTitle)
	pctx.Set("groupuuid", vars["uuid"])
	pctx.Set("usersInGroup", usersInGroup)
	pctx.Set("usersNotInGroup", usersNotInGroup)
	pctx.Set("unixtostring", UnixToTimeString)
	pctx.Set("quillenabled", false)
	pctx.Set("adminhiddenpassword", "")
	if augeh.Router.AdminHidden {
		pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", augeh.Router.AdminHiddenPassword))
	}
	RenderDefault(w, "admin.users.groups.edit.html", pctx)
}

func (augeh *AdminUserGroupsEditHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		http.Redirect(w, r, augeh.route, http.StatusFound)
	}

	//gmt := db.GroupMembershipTable{}
}

//Route get URI route for handler
func (augeh *AdminUserGroupsEditHandler) Route() string { return augeh.route }

//HandlesGet retrieve whether this handler handles get requests
func (augeh *AdminUserGroupsEditHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (augeh *AdminUserGroupsEditHandler) HandlesPost() bool { return true }
