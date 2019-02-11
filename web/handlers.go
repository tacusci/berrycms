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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	var adminHiddenPrefix = ""

	if router.AdminHidden {
		adminHiddenPrefix = fmt.Sprintf("/%s", router.AdminHiddenPassword)
	}

	return []Handler{
		&RobotsHandler{
			route:  "/robots.txt",
			Router: router,
		},
		&LoginHandler{
			route:  adminHiddenPrefix + "/login",
			Router: router,
		},
		&LogoutHandler{
			route:  adminHiddenPrefix + "/logout",
			Router: router,
		},
		&AdminHandler{
			route:  adminHiddenPrefix + "/admin",
			Router: router,
		},
		&AdminUsersHandler{
			route:  adminHiddenPrefix + "/admin/users",
			Router: router,
		},
		&AdminUsersNewHandler{
			route:  adminHiddenPrefix + "/admin/users/new",
			Router: router,
		},
		&AdminUsersDeleteHandler{
			route:  adminHiddenPrefix + "/admin/users/delete",
			Router: router,
		},
		&AdminPagesHandler{
			route:  adminHiddenPrefix + "/admin/pages",
			Router: router,
		},
		&AdminPagesNewHandler{
			route:  adminHiddenPrefix + "/admin/pages/new",
			Router: router,
		},
		&AdminPagesEditHandler{
			route:  adminHiddenPrefix + "/admin/pages/edit/{uuid}",
			Router: router,
		},
		&AdminPagesDeleteHandler{
			route:  adminHiddenPrefix + "/admin/pages/delete",
			Router: router,
		},
		&AdminUserGroupsHandler{
			route:  adminHiddenPrefix + "/admin/users/groups",
			Router: router,
		},
		&AdminUserGroupsNewHandler{
			route:  adminHiddenPrefix + "/admin/users/groups/new",
			Router: router,
		},
		&AdminUserGroupsEditHandler{
			route:  adminHiddenPrefix + "/admin/users/groups/edit/{uuid}",
			Router: router,
		},
		&AdminUserGroupsEditAddHandler{
			route:  adminHiddenPrefix + "/admin/users/groups/edit/{uuid}/add",
			Router: router,
		},
		&AdminUserGroupsEditRemoveHandler{
			route:  adminHiddenPrefix + "/admin/users/groups/edit/{uuid}/remove",
			Router: router,
		},
		&AdminUserGroupsDeleteHandler{
			route:  adminHiddenPrefix + "/admin/users/groups/delete",
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
	var respCode = http.StatusOK
	var htmlHead = "<head><link rel=\"stylesheet\" href=\"/css/berry-default.css\"><link rel=\"stylesheet\" href=\"/css/font.css\"></head>"
	var respBytesData []byte

	//render page from plush template
	html, err := plush.Render("<html>"+htmlHead+"<body><%= pagecontent %></body></html>", ctx)
	if err != nil {
		Error(w, err)
		return err
	}

	redirectRequested := false

	pm := plugins.NewManager()

	//have to lock as unfortunately do not support calling any plugin function twice at the exact same time
	pm.Lock()
	for _, plugin := range *pm.Plugins() {
		plugin.Document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			logging.Error(err.Error())
			break
		}
		plugin.VM.Set("document", plugin.Document)
		val, err := plugin.Call("on_get_render", nil, &p.Route)
		//val, err := plugin.Call("onGetRender", nil, &p.Route)
		if err != nil {
			logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", plugin.UUID(), err.Error()))
			continue
		}

		if &val != nil && val.IsObject() {
			editedPage := val.Object()

			editedPageRoute, err := editedPage.Get("route")
			//don't want to respond with 500 to user
			if err != nil {
				logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", plugin.UUID(), err.Error()))
				continue
			}

			modifiedStatusCode, err := editedPage.Get("code")
			if err != nil {
				logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", plugin.UUID(), err.Error()))
				continue
			}

			data, err := editedPage.Get("data")
			if err != nil {
				logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", plugin.UUID(), err.Error()))
				continue
			}

			if editedPageRoute.IsString() {
				//plugin has modified current page route, registering pending redirect
				if editedPageRoute.String() != p.Route {
					redirectRequested = true
					//by default use status found
					respCode = http.StatusFound
					p.Route = editedPageRoute.String()
				}
			}

			if modifiedStatusCode.IsNumber() {
				modifiedStatusCodeInt, err := modifiedStatusCode.ToInteger()
				if err != nil {
					logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", plugin.UUID(), err.Error()))
					continue
				}
				respCode = int(modifiedStatusCodeInt)
			}

			//if it's a golang slice it comes back as an object
			if data.IsObject() {
				if val_interface, err := data.Export(); err == nil {
					if dataBytes, ok := val_interface.([]byte); ok {
						respBytesData = dataBytes
						break
					}
				}
			}
		}

		//no point in running other plugins
		if redirectRequested {
			break
		}

		htmlFromDocument, err := plugin.Document.Html()
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		//only set the returned HTML if retrieving it hasn't errored
		html = htmlFromDocument
	}
	pm.Unlock()

	if redirectRequested {
		http.Redirect(w, r, p.Route, respCode)
		return nil
	}

	if len(respBytesData) > 0 {
		http.ServeContent(w, r, "", time.Now(), bytes.NewReader(respBytesData))
		return nil
	}

	if respCode == http.StatusNotFound {
		ctx, err := renderFourOhFour()
		if err != nil {
			return err
		}
		html = RenderStr(ctx)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(respCode)
	w.Write([]byte(html))
	return nil
}

//RenderStr uses plush rendering engine to read page content from the DB and create HTML content as string
func RenderStr(ctx *plush.Context) string {
	html, err := plush.Render("<html><head><link rel=\"stylesheet\" href=\"/css/berry-default.css\"><link rel=\"stylesheet\" href=\"/css/font.css\"></head><%= pagecontent %></html>", ctx)
	if err != nil {
		logging.Error(err.Error())
		return "<h1>500 Server Error</h1>"
	}
	return html
}
