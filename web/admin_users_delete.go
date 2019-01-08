// Copyright (c) 2019 tacusci ltd
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

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminUsersDeleteHandler handler to contain pointer to core router and the URI string
type AdminUsersDeleteHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (audh *AdminUsersDeleteHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (audh *AdminUsersDeleteHandler) Post(w http.ResponseWriter, r *http.Request) {

	var redirectURI = "/admin/users"

	if audh.Router.AdminHidden {
		redirectURI = fmt.Sprintf("/%s", audh.Router.AdminHiddenPassword) + redirectURI
	}

	defer http.Redirect(w, r, redirectURI, http.StatusFound)

	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		Error(w, err)
	}

	ut := db.UsersTable{}
	st := db.AuthSessionsTable{}
	pt := db.PagesTable{}
	amw := AuthMiddleware{}

	loggedInUser, err := amw.LoggedInUser(r)

	for _, v := range r.PostForm {
		userToDelete, err := ut.SelectByUUID(db.Conn, v[0])

		if err != nil {
			logging.Error(err.Error())
			continue
		}

		if userToDelete == nil {
			continue
		}

		//don't allow deletion of the root user account
		if db.UsersRoleFlag(userToDelete.UserroleId) != db.ROOT_USER {

			if err != nil {
				logging.Error(err.Error())
				Error(w, err)
			}

			//make sure that the logged in user is not the same as user to delete
			//the first condition evals before the second, that way no nil pointer exception occurs
			if (loggedInUser != nil) && (loggedInUser.UUID != userToDelete.UUID) {
				rows, err := pt.Select(db.Conn, "uuid", fmt.Sprintf("authoruuid = '%s'", userToDelete.UUID))

				if err != nil {
					logging.Error(err.Error())
					Error(w, err)
				}

				defer rows.Close()

				rowCount := 0
				for rows.Next() {
					if rowCount > 0 {
						break
					}
					rowCount++
				}

				//make sure that the user to delete isn't the author of any pages (should probably do something different to this in future)
				if rowCount == 0 {
					st.Delete(db.Conn, fmt.Sprintf("uuid = '%s'", userToDelete.UUID))
					ut.DeleteByUUID(db.Conn, userToDelete.UUID)
					gmt := db.GroupMembershipTable{}
					//will delete user from all groups, maybe this should be a different function?
					gmt.DeleteUserFromGroup(db.Conn, userToDelete, &db.Group{UUID: "*"})
				}
			}
		}
	}
}

//Route get URI route for handler
func (audh *AdminUsersDeleteHandler) Route() string { return audh.route }

//HandlesGet retrieve whether this handler handles get requests
func (audh *AdminUsersDeleteHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (audh *AdminUsersDeleteHandler) HandlesPost() bool { return true }
