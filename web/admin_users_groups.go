package web

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/logging"

	"github.com/tacusci/berrycms/db"
)

type AdminUserGroupsHandler struct {
	Router *MutableRouter
	route  string
}

func (ugh *AdminUserGroupsHandler) Get(w http.ResponseWriter, r *http.Request) {

	gt := db.GroupTable{}

	groupRows, err := gt.Select(db.Conn, "createddatetime, uuid, title", "")
	if err != nil {
		Error(w, err)
	}

	defer groupRows.Close()

	for groupRows.Next() {
		group := db.Group{}
		groupRows.Scan(&group.CreatedDateTime, &group.UUID, &group.Title)

		gmt := db.GroupMembershipTable{}

		groupMembershipRows, err := gmt.Select(db.Conn, "groupuuid, useruuid", fmt.Sprintf("groupuuid = '%s'", group.UUID))
		if err != nil {
			Error(w, err)
		}

		for groupMembershipRows.Next() {
			groupMembership := db.GroupMembership{}
			groupMembershipRows.Scan(&groupMembership.CreatedDateTime, &groupMembership.UUID, &groupMembership.GroupUUID, &groupMembership.UserUUID)

			ut := db.UsersTable{}

			userRows, err := ut.Select(db.Conn, "createddatetime, userroleid, uuid, username, authhash, firstname, lastname, email", fmt.Sprintf("uuid = '%s'", groupMembership.UserUUID))
			if err != nil {
				Error(w, err)
			}

			for userRows.Next() {
				user := db.User{}
				userRows.Scan(&user.CreatedDateTime, &user.UserroleId, &user.UUID, &user.Username, &user.AuthHash, &user.FirstName, &user.LastName, &user.Email)
				logging.Debug(fmt.Sprintf("Found user of UUID %s in group %s", user.UUID, group.UUID))
			}

			userRows.Close()
		}

		groupMembershipRows.Close()
	}

	pctx := plush.NewContext()
	pctx.Set("unixtostring", UnixToTimeString)
	pctx.Set("title", "Groups")
	pctx.Set("quillenabled", false)

	RenderDefault(w, "admin.users.groups.html", pctx)
}

func (ugh *AdminUserGroupsHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (ugh *AdminUserGroupsHandler) Route() string { return ugh.route }

func (ugh *AdminUserGroupsHandler) HandlesGet() bool { return true }

func (ugh *AdminUserGroupsHandler) HandlesPost() bool { return false }
