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
	"html/template"
	"net/http"

	"github.com/tacusci/logging"

	quill "github.com/dchenk/go-render-quill"
	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
)

//SavedPageHandler handler to contain pointer to core router and the URI string
type SavedPageHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (sph *SavedPageHandler) Get(w http.ResponseWriter, r *http.Request) {
	pt := db.PagesTable{}

	p, err := pt.SelectByRoute(db.Conn, r.RequestURI)

	if err != nil {
		logging.Error(err.Error())
	}

	if p == nil {
		fourOhFour(w, r)
		return
	}

	ctx := plush.NewContext()
	ctx.Set("pagecontent", template.HTML(p.Content))

	// if trying to render the page content from delta fails, then it just won't replace previous context pagecontent value
	if html, _ := quill.Render([]byte(p.Content)); err == nil {
		ctx.Set("pagecontent", template.HTML(html))
	}

	Render(w, r, p, ctx)
}

//Post handles post requests to URI
func (sph *SavedPageHandler) Post(w http.ResponseWriter, r *http.Request) {
	pt := db.PagesTable{}

	p, err := pt.SelectByRoute(db.Conn, r.RequestURI)

	if err != nil {
		logging.Error(err.Error())
	}

	if p == nil {
		fourOhFour(w, r)
		return
	}
}

//Route get URI route for handler
func (sph *SavedPageHandler) Route() string { return sph.route }

//HandlesGet retrieve whether this handler handles get requests
func (sph *SavedPageHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (sph *SavedPageHandler) HandlesPost() bool { return true }
