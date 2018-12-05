package web

import (
	"fmt"
	"github.com/gobuffalo/plush"
	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
	"net/http"
)

type AdminUserGroupsEditHandler struct {
	Router *MutableRouter
	route  string
}

func (augeh *AdminUserGroupsEditHandler) Get(w http.ResponseWriter, r *http.Request) {
	var users = make([]db.User, 0)
	vars := mux.Vars(r)
	gmt := db.GroupMembershipTable{}

	rows, err := gmt.Select(db.Conn, "createddatetime, groupuuid, useruuid", fmt.Sprintf("groupuuid = '%s'", vars["uuid"]))
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("Group memberships not found"))
		return
	}

	for rows.Next() {
		groupMembership := db.GroupMembership{}
		err := rows.Scan(&groupMembership.CreatedDateTime, &groupMembership.GroupUUID, &groupMembership.UserUUID)
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		ut := db.UsersTable{}

		user, err := ut.SelectByUUID(db.Conn, groupMembership.UserUUID)
		if err != nil {
			logging.Error(err.Error())
			continue
		}
		users = append(users, *user)
	}

	pctx := plush.NewContext()
	pctx.Set("title", "Edit Group")
	pctx.Set("submitroute", r.RequestURI)
	pctx.Set("users", users)
	pctx.Set("unixtostring", UnixToTimeString)
	pctx.Set("quillenabled", false)
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
