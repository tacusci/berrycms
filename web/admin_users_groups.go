package web

import (
	"net/http"

	"github.com/tacusci/berrycms/db"
)

type AdminUserGroupsHandler struct {
	Router *MutableRouter
	route  string
}

func (ugh *AdminUserGroupsHandler) Get(w http.ResponseWriter, r *http.Request) {
	groups := make([]db.Group, 0)
	groupMemberships := make([]db.GroupMembership, 0)

	gt := db.GroupTable{}
	rows, err := gt.Select(db.Conn, "createddatetime, uuid, title", "")

	if err != nil {
		Error(w, err)
	}

	rows.Close()

	gmt := db.GroupMembershipTable{}
	rows, err = gmt.Select(db.Conn, "createddatetime", "uuid", "groupuuid", "useruuid")

	if err != nil {
		Error(w, err)
	}

	rows.Close()
}
