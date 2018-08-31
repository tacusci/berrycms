package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

type AdminPagesNewHandler struct {
	Router *MutableRouter
	route  string
}

//Get takes the web request and writes response to session
func (apnh *AdminPagesNewHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "New Page")
	pctx.Set("submitroute", r.RequestURI)
	pctx.Set("pagetitle", "")
	pctx.Set("pageroute", "")
	pctx.Set("pagecontent", "")
	pctx.Set("quillenabled", true)
	RenderDefault(w, "admin.pages.new.html", pctx)
}

func (apnh *AdminPagesNewHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		return
	}

	pt := db.PagesTable{}
	pageToCreate := db.Page{}

	pageToCreate.CreatedDateTime = time.Now().Unix()
	pageToCreate.Title = r.PostFormValue("title")
	pageToCreate.Route = r.PostFormValue("route")
	pageToCreate.Content = r.PostFormValue("pagecontent")

	err = pt.Insert(db.Conn, pageToCreate)

	if err != nil {
		logging.Error(err.Error())
	}

	pageToCreate, err = pt.SelectByRoute(db.Conn, pageToCreate.Route)

	if err != nil {
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
	}

	apnh.Router.Reload()

	http.Redirect(w, r, fmt.Sprintf("/admin/pages/edit/%s", pageToCreate.UUID), http.StatusFound)
}

func (apnh *AdminPagesNewHandler) Route() string { return apnh.route }

func (apnh *AdminPagesNewHandler) HandlesGet() bool  { return true }
func (apnh *AdminPagesNewHandler) HandlesPost() bool { return true }
