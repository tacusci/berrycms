// Copyright (c) 2018, tacusci ltd
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
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminUserGroupsEditAddHandler contains response functions for pages admin page
type AdminUserGroupsEditAddHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Post(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	groupUUID := vars["uuid"]

	if groupUUID == "" {
		Error(w, errors.New("Missing group UUID"))
		return
	}

	var redirectURI = "/admin/users/groups/edit/" + groupUUID

	if augeah.Router.AdminHidden {
		redirectURI = fmt.Sprintf("/%s", augeah.Router.AdminHiddenPassword) + redirectURI
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

		groupTitle := ""
		rows, err := gt.Select(db.Conn, "title", fmt.Sprintf("uuid = '%s'", groupUUID))
		if err != nil {
			Error(w, err)
			return
		}

		for rows.Next() {
			rows.Scan(&groupTitle)
		}

		rows.Close()

		for _, v := range r.PostForm {
			userToAdd, err := ut.SelectByUUID(db.Conn, v[0])
			if err != nil {
				logging.Error(err.Error())
				continue
			}

			if userToAdd == nil {
				continue
			}

			gmt.AddUserToGroup(db.Conn, userToAdd, groupTitle)
		}
	}
}

//Route get URI route for handler
func (augeah *AdminUserGroupsEditAddHandler) Route() string { return augeah.route }

//HandlesGet retrieve whether this handler handles get requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesPost() bool { return true }
