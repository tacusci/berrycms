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
	"fmt"
	"net/http"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminPagesHandler handler to contain pointer to core router and the URI string
type AdminPagesHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (aph *AdminPagesHandler) Get(w http.ResponseWriter, r *http.Request) {
	pages := make([]db.Page, 0)
	authors := make([]string, 0)

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "createddatetime, uuid, title, route, authoruuid", "")

	if err != nil {
		Error(w, err)
	}

	defer rows.Close()

	ut := db.UsersTable{}

	for rows.Next() {
		p := db.Page{}
		rows.Scan(&p.CreatedDateTime, &p.UUID, &p.Title, &p.Route, &p.AuthorUUID)
		pages = append(pages, p)

		authorUser, err := ut.SelectByUUID(db.Conn, p.AuthorUUID)

		if err != nil {
			logging.Error(err.Error())
		} else {
			authors = append(authors, fmt.Sprintf("%s %s", authorUser.FirstName, authorUser.LastName))
		}
	}

	pctx := plush.NewContext()
	pctx.Set("unixtostring", UnixToTimeString)
	pctx.Set("title", "Pages")
	pctx.Set("quillenabled", false)
	pctx.Set("pages", pages)
	pctx.Set("authors", authors)
	pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", aph.Router.AdminHiddenPassword))

	RenderDefault(w, "admin.pages.html", pctx)
}

//Post handles post requests to URI
func (aph *AdminPagesHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (aph *AdminPagesHandler) Route() string { return aph.route }

//HandlesGet retrieve whether this handler handles get requests
func (aph *AdminPagesHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (aph *AdminPagesHandler) HandlesPost() bool { return false }
