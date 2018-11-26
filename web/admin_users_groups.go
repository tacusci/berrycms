package web

import (
	"github.com/gobuffalo/plush"
	"net/http"

	"github.com/tacusci/berrycms/db"
)

type AdminUserGroupsHandler struct {
	Router *MutableRouter
	route  string
}

func (ugh *AdminUserGroupsHandler) Get(w http.ResponseWriter, r *http.Request) {

	var groups = make([]db.Group, 0)
	var groupMemberships = make([]db.GroupMembership, 0)

	groupTable := db.GroupTable{}
	groupRows, err := groupTable.Select(db.Conn, "*", "")

	if err != nil {
		Error(w, err)
	}

	defer groupRows.Close()

	for groupRows.Next() {
		g := db.Group{}

		groupRows.Scan(&g.CreatedDateTime, &g.UUID, &g.Title)
		groups = append(groups, g)
		groupMembershipTable := db.GroupMembershipTable{}
		groupMembershipRows, err := groupMembershipTable.Select(db.Conn, "*", "groupuuid = "+g.UUID)

		if err != nil {
			Error(w, err)
		}

		defer groupMembershipRows.Close()

		for groupMembershipRows.Next() {
			gm := db.GroupMembership{}

			groupMembershipRows.Scan(&gm.CreatedDateTime, &gm.UUID, &gm.GroupUUID, &gm.UserUUID)
			groupMemberships = append(groupMemberships, gm)
			userTable := db.UsersTable{}
			userRows, err := userTable.Select(db.Conn, "*", "uuid = "+gm.UserUUID)

			if err != nil {
				Error(w, err)
			}

			defer userRows.Close()
		}
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
