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
	"html/template"
	"net/http"

	"github.com/tacusci/logging"

	quill "github.com/dchenk/go-render-quill"
	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/plugins"
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
	// assume response is fine/OK
	var respCode = http.StatusFound

	pt := db.PagesTable{}

	p, err := pt.SelectByRoute(db.Conn, r.RequestURI)

	if err != nil {
		Error(w, err)
		return
	}

	if p == nil {
		fourOhFour(w, r)
		return
	}

	err = r.ParseForm()

	if err != nil {
		Error(w, err)
		return
	}

	redirectRequested := false

	pm := plugins.NewManager()

	//have to lock as unfortunately do not support calling any plugin function twice at the exact same time
	pm.Lock()
	for _, plugin := range *pm.Plugins() {
		if err != nil {
			logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", plugin.UUID(), err.Error()))
			continue
		}
		val, err := plugin.Call("on_post_recieve", nil, &p.Route, r.PostForm)
		if err != nil {
			plugin.Error(err)
			continue
		}

		if &val != nil && val.IsObject() {
			editedPage := val.Object()

			editedPageRoute, err := editedPage.Get("route")
			//don't want to respond with 500 to user
			if err != nil {
				plugin.Error(err)
				continue
			}

			if editedPageRoute.IsString() {
				//plugin has modified current page route, registering pending redirect
				if editedPageRoute.String() != "" {
					redirectRequested = true
					//by default use status found
					respCode = http.StatusFound
					modifiedStatusCode, err := editedPage.Get("code")
					if err != nil {
						plugin.Error(err)
						continue
					}
					if modifiedStatusCode.IsNumber() {
						modifiedStatusCodeInt, err := modifiedStatusCode.ToInteger()
						if err != nil {
							plugin.Error(err)
							continue
						}
						respCode = int(modifiedStatusCodeInt)
					}
					p.Route = editedPageRoute.String()
				}
			}
		}

		//no point in running other plugins
		if redirectRequested {
			break
		}
	}
	pm.Unlock()

	if redirectRequested {
		http.Redirect(w, r, p.Route, respCode)
		return
	}
}

//Route get URI route for handler
func (sph *SavedPageHandler) Route() string { return sph.route }

//HandlesGet retrieve whether this handler handles get requests
func (sph *SavedPageHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (sph *SavedPageHandler) HandlesPost() bool { return true }
