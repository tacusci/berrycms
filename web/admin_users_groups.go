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

	"github.com/gobuffalo/plush"
	"github.com/tacusci/logging"

	"github.com/tacusci/berrycms/db"
)

type AdminUserGroupsHandler struct {
	Router *MutableRouter
	route  string
}

func (ugh *AdminUserGroupsHandler) Get(w http.ResponseWriter, r *http.Request) {
	groups := make([]db.Group, 0)

	groupTable := db.GroupTable{}
	rows, err := groupTable.Select(db.Conn, "createddatetime, uuid, title", "")

	if err != nil {
		Error(w, err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		group := db.Group{}
		err := rows.Scan(&group.CreatedDateTime, &group.UUID, &group.Title)

		if err != nil {
			logging.Error(err.Error())
			continue
		}

		groups = append(groups, group)
	}

	pctx := plush.NewContext()
	pctx.Set("unixtostring", UnixToTimeString)
	pctx.Set("title", "Groups")
	pctx.Set("adminhiddenpassword", "")
	pctx.Set("quillenabled", false)
	pctx.Set("newgroupformaction", "/admin/users/groups/new")
	pctx.Set("groups", groups)
	if ugh.Router.AdminHidden {
		pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", ugh.Router.AdminHiddenPassword))
	}

	RenderDefault(w, "admin.users.groups.html", pctx)
}

func (ugh *AdminUserGroupsHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (ugh *AdminUserGroupsHandler) Route() string { return ugh.route }

func (ugh *AdminUserGroupsHandler) HandlesGet() bool { return true }

func (ugh *AdminUserGroupsHandler) HandlesPost() bool { return false }
