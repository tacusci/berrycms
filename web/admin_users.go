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

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
)

//AdminUsersHandler handler to contain pointer to core router and the URI string
type AdminUsersHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (uh *AdminUsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	users := make([]db.User, 0)

	ut := db.UsersTable{}
	rows, err := ut.Select(db.Conn, "createddatetime, uuid, firstname, lastname, username, email", "")
	defer rows.Close()

	if err != nil {
		Error(w, err)
	}

	for rows.Next() {
		u := db.User{}
		rows.Scan(&u.CreatedDateTime, &u.UUID, &u.FirstName, &u.LastName, &u.Username, &u.Email)
		users = append(users, u)
	}

	pctx := plush.NewContext()
	pctx.Set("users", users)
	pctx.Set("title", "Users")
	pctx.Set("quillenabled", false)
	pctx.Set("adminhiddenpassword", "")
	if uh.Router.AdminHidden {
		pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", uh.Router.AdminHiddenPassword))
	}
	pctx.Set("unixtostring", UnixToTimeString)

	RenderDefault(w, "admin.users.html", pctx)
}

//Post handles post requests to URI
func (uh *AdminUsersHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (uh *AdminUsersHandler) Route() string { return uh.route }

//HandlesGet retrieve whether this handler handles get requests
func (uh *AdminUsersHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (uh *AdminUsersHandler) HandlesPost() bool { return false }
