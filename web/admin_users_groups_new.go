package web

import (
	"fmt"
	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
	"net/http"
	"time"
)

//AdminUserGroupsNewHandler handler to contain pointer to core router and the URI string
type AdminUserGroupsNewHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (augnh *AdminUserGroupsNewHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "New Group")
	pctx.Set("submitroute", r.RequestURI)
	pctx.Set("pagetitle", "")
	pctx.Set("pageroute", "")
	pctx.Set("pagecontent", "")
	pctx.Set("quillenabled", false)
	RenderDefault(w, "admin.users.groups.new.html", pctx)
}

//Post handles post requests to URI
func (augnh *AdminUserGroupsNewHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		http.Redirect(w, r, augnh.route, http.StatusFound)
	}

	gt := db.GroupTable{}

	groupToCreate := &db.Group{
		CreatedDateTime: time.Now().Unix(),
		Title:           r.PostFormValue("title"),
	}

	err = gt.Insert(db.Conn, groupToCreate)

	if err != nil {
		logging.Error(err.Error())
	}

	groupToCreate, err = gt.SelectByTitle(db.Conn, groupToCreate.Title)

	if err != nil {
		logging.Error(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
	}

	augnh.Router.Reload()

	http.Redirect(w, r, fmt.Sprintf("/admin/users/groups/edit/%s", groupToCreate.UUID), http.StatusFound)
}

//Route get URI route for handler
func (augnh *AdminUserGroupsNewHandler) Route() string { return augnh.route }

//HandlesGet retrieve whether this handler handles get requests
func (augnh *AdminUserGroupsNewHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (augnh *AdminUserGroupsNewHandler) HandlesPost() bool { return true }
