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
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/tacusci/berrycms/plugins"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//Handler root interface describing handler structure
type Handler interface {
	Route() string
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	HandlesGet() bool
	HandlesPost() bool
}

//GetDefaultHandlers get fixed list of all default handlers
func GetDefaultHandlers(router *MutableRouter) []Handler {
	return []Handler{
		&LoginHandler{
			route:  "/login",
			Router: router,
		},
		&LogoutHandler{
			route:  "/logout",
			Router: router,
		},
		&AdminHandler{
			route:  "/admin",
			Router: router,
		},
		&AdminUsersHandler{
			route:  "/admin/users",
			Router: router,
		},
		&AdminUsersNewHandler{
			route:  "/admin/users/new",
			Router: router,
		},
		&AdminUsersDeleteHandler{
			route:  "/admin/users/delete",
			Router: router,
		},
		&AdminPagesHandler{
			route:  "/admin/pages",
			Router: router,
		},
		&AdminPagesNewHandler{
			route:  "/admin/pages/new",
			Router: router,
		},
		&AdminPagesEditHandler{
			route:  "/admin/pages/edit/{uuid}",
			Router: router,
		},
		&AdminPagesDeleteHandler{
			route:  "/admin/pages/delete",
			Router: router,
		},
		&AdminUserGroupsHandler{
			route:  "/admin/users/groups",
			Router: router,
		},
		&AdminUserGroupsNewHandler{
			route:  "/admin/users/groups/new",
			Router: router,
		},
	}
}

//UnixToTimeString take unix time and convert to string of to European time format
func UnixToTimeString(unix int64) string {
	return time.Unix(unix, 0).Format("15:04:05 02-01-2006")
}

//RenderDefault uses plush rendering engine to take default page template and create HTML content
func RenderDefault(w http.ResponseWriter, template string, pctx *plush.Context) error {
	header, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "header.snip")

	if err != nil {
		Error(w, err)
		return err
	}

	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + template)
	if err != nil {
		Error(w, err)
		return err
	}

	renderedContent, err := plush.Render(string(append(append(header, []byte("\n")...), content...))+"\n</html>", pctx)

	if err != nil {
		Error(w, err)
		return err
	}

	_, err = w.Write([]byte(renderedContent))
	return err
}

//Render uses plush rendering engine to read page content from the DB and create HTML content
func Render(w http.ResponseWriter, r *http.Request, p *db.Page, ctx *plush.Context) error {

	// assume response is fine/OK
	var code = http.StatusOK
	var htmlHead = "<head><link rel=\"stylesheet\" href=\"/css/berry-default.css\"><link rel=\"stylesheet\" href=\"/css/font.css\"></head>"

	pm := plugins.NewManager()

	pm.Lock()
	for _, plugin := range *pm.Plugins() {
		val, _ := plugin.Call("onPreRender", nil, &p.Route, &htmlHead, &p.Content)
		if &val != nil && val.IsObject() {
			editedPage := val.Object()

			if editedPageRoute, err := editedPage.Get("route"); err == nil {
				if editedPageRoute.IsString() {
					if editedPageRoute.String() != p.Route {
						http.Redirect(w, r, editedPageRoute.String(), http.StatusFound)
						return nil
					}
				}
			}

			if editedPageHeader, err := editedPage.Get("header"); err == nil {
				if editedPageHeader.IsString() {
					htmlHead = editedPageHeader.String()
				}
			} else {
				logging.Error(fmt.Sprintf("Error from plugin {%s} -> %s", plugin.UUID(), err.Error()))
			}

			if editedPageContent, err := editedPage.Get("body"); err == nil {
				if editedPageContent.IsString() {
					p.Content = editedPageContent.String()
					ctx.Set("pagecontent", template.HTML(p.Content))
				}
			} else {
				logging.Error(fmt.Sprintf("Error from plugin {%s} -> %s", plugin.UUID(), err.Error()))
			}

			if editedPageResponseCode, err := editedPage.Get("code"); err == nil {
				if editedPageResponseCode.IsNumber() {
					if responseCode, err := editedPageResponseCode.ToInteger(); err == nil {
						code = int(responseCode)
					} else {
						logging.Error(fmt.Sprintf("Error from plugin {%s} -> %s", plugin.UUID(), err.Error()))
					}
				}
			}
		}
	}
	pm.Unlock()

	html, err := plush.Render("<html>"+htmlHead+"<body><%= pagecontent %></body></html>", ctx)
	if err != nil {
		Error(w, err)
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write([]byte(html))
	return err
}

//RenderStr uses plush rendering engine to read page content from the DB and create HTML content as string
func RenderStr(p *db.Page, ctx *plush.Context) string {
	html, err := plush.Render("<html><head><link rel=\"stylesheet\" href=\"/css/berry-default.css\"><link rel=\"stylesheet\" href=\"/css/font.css\"></head><%= pagecontent %></html>", ctx)
	if err != nil {
		logging.Error(err.Error())
		return "<h1>500 Server Error</h1>"
	}
	return html
}
