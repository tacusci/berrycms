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
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
	"net/http"
)

//AdminPagesDeleteHandler handler to contain pointer to core router and the URI string
type AdminPagesDeleteHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (apdh *AdminPagesDeleteHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (apdh *AdminPagesDeleteHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		Error(w, err)
	}

	pt := db.PagesTable{}
	deletedPages := false
	for _, v := range r.PostForm {
		deletedPages = true
		pt.DeleteByUUID(db.Conn, v[0])
	}

	if deletedPages {
		apdh.Router.Reload()
	}

	var redirectURI = "/admin/pages"

	if apdh.Router.AdminHidden {
		redirectURI = fmt.Sprintf("/%s", apdh.Router.AdminHiddenPassword) + redirectURI
	}

	http.Redirect(w, r, redirectURI, http.StatusFound)
}

//Route get URI route for handler
func (apdh *AdminPagesDeleteHandler) Route() string { return apdh.route }

//HandlesGet retrieve whether this handler handles get requests
func (apdh *AdminPagesDeleteHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (apdh *AdminPagesDeleteHandler) HandlesPost() bool { return true }
