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
)

//AdminHandler handler to contain pointer to core router and the URI string
type AdminHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (ah *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "Dashboard")
	pctx.Set("quillenabled", false)
	pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", ah.Router.AdminHiddenPassword))
	RenderDefault(w, "admin.html", pctx)
}

//Post handles post requests to URI
func (ah *AdminHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (ah *AdminHandler) Route() string { return ah.route }

//HandlesGet retrieve whether this handler handles get requests
func (ah *AdminHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (ah *AdminHandler) HandlesPost() bool { return false }
