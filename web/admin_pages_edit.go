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
	"net/http"
	"strings"

	"github.com/dchenk/go-render-quill"

	"github.com/tacusci/logging"

	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"

	"github.com/gobuffalo/plush"
)

//AdminPagesEditHandler contains response functions for pages admin page
type AdminPagesEditHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (apeh *AdminPagesEditHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pt := db.PagesTable{}
	pageToEdit, err := pt.SelectByUUID(db.Conn, vars["uuid"])
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("Page to edit not found"))
		return
	}

	if html, err := quill.Render([]byte(pageToEdit.Content)); err == nil {
		pctx := plush.NewContext()
		pctx.Set("title", fmt.Sprintf("Edit Page - %s", pageToEdit.Title))
		pctx.Set("submitroute", r.RequestURI)
		pctx.Set("pagetitle", pageToEdit.Title)
		pctx.Set("pageroute", pageToEdit.Route)
		pctx.Set("pagecontent", template.HTML(string(html)))
		pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", apeh.Router.AdminHiddenPassword))
		pctx.Set("quillenabled", true)
		RenderDefault(w, "admin.pages.edit.html", pctx)
	} else {
		Error(w, err)
	}
}

//Post handles post requests to URI
func (apeh *AdminPagesEditHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, r.RequestURI, http.StatusFound)
	vars := mux.Vars(r)
	pt := db.PagesTable{}
	pageToEdit, err := pt.SelectByUUID(db.Conn, vars["uuid"])
	if err != nil {
		logging.Error(err.Error())
		return
	}

	if pageToEdit == nil {
		logging.Error(fmt.Sprintf("%s", "Page to edit doesn't exist, stopping..."))
		return
	}

	err = r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		return
	}

	pageToEdit.Title = r.PostFormValue("title")
	oldPageRoute := pageToEdit.Route
	pageToEdit.Route = r.PostFormValue("route")
	pageToEdit.Content = r.PostFormValue("pagecontent")

	err = pt.Update(db.Conn, pageToEdit)

	if err != nil {
		logging.Error(err.Error())
	}

	//reloading all page routes is potentially really intensive, so only do this if the route has actually changed
	if strings.Compare(oldPageRoute, pageToEdit.Route) != 0 {
		apeh.Router.Reload()
	}
}

//Route get URI route for handler
func (apeh *AdminPagesEditHandler) Route() string { return apeh.route }

//HandlesGet retrieve whether this handler handles get requests
func (apeh *AdminPagesEditHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (apeh *AdminPagesEditHandler) HandlesPost() bool { return true }
