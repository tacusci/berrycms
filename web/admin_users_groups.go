package web

import (
	"fmt"
	"github.com/gobuffalo/plush"
	"github.com/tacusci/logging"
	"net/http"

	"github.com/tacusci/berrycms/db"
)

type AdminUserGroupsHandler struct {
	Router *MutableRouter
	route  string
}

func (ugh *AdminUserGroupsHandler) Get(w http.ResponseWriter, r *http.Request) {

	gt := db.GroupTable{}

	rows, err := gt.Select(db.Conn, "createddatetime, uuid, title", "")
	if err != nil {
		Error(w, err)
	}

	for rows.Next() {
		group := db.Group{}
		rows.Scan(&group.CreatedDateTime, &group.UUID, &group.Title)
		logging.Debug(fmt.Sprintf("Loaded group {CT: %s, UUID: %s, Title: %s}", UnixToTimeString(group.CreatedDateTime), group.UUID, group.Title))
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
