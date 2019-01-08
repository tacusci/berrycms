// Copyright (c) 2019, tacusci ltd
//
// Licensed under the GNU GENERAL PUBLIC LICENSE Version 3 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/gpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
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
	pctx.Set("adminhiddenpassword", "")
	if augnh.Router.AdminHidden {
		pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", augnh.Router.AdminHiddenPassword))
	}
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

	var redirectURI = "/admin/users/groups/edit/%s"

	if augnh.Router.AdminHidden {
		redirectURI = fmt.Sprintf("/%s", augnh.Router.AdminHiddenPassword) + redirectURI
	}

	http.Redirect(w, r, fmt.Sprintf(redirectURI, groupToCreate.UUID), http.StatusFound)
}

//Route get URI route for handler
func (augnh *AdminUserGroupsNewHandler) Route() string { return augnh.route }

//HandlesGet retrieve whether this handler handles get requests
func (augnh *AdminUserGroupsNewHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (augnh *AdminUserGroupsNewHandler) HandlesPost() bool { return true }
