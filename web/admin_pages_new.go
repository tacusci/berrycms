package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminPagesNewHandler handler to contain pointer to core router and the URI string
type AdminPagesNewHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
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

//Post handles post requests to URI
func (apnh *AdminPagesNewHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		http.Redirect(w, r, "/admin/pages/new", http.StatusFound)
	}

	pt := db.PagesTable{}

	amw := AuthMiddleware{}
	loggedInUser, err := amw.LoggedInUser(r)

	if err != nil {
		logging.Error(err.Error())
		http.Redirect(w, r, "/admin/pages/new", http.StatusFound)
	}

	pageToCreate := db.Page{
		CreatedDateTime: time.Now().Unix(),
		Title:           r.PostFormValue("title"),
		AuthorUUID:      loggedInUser.UUID,
		Route:           r.PostFormValue("route"),
		Content:         r.PostFormValue("pagecontent"),
	}

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

//Route get URI route for handler
func (apnh *AdminPagesNewHandler) Route() string { return apnh.route }

//HandlesGet retrieve whether this handler handles get requests
func (apnh *AdminPagesNewHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (apnh *AdminPagesNewHandler) HandlesPost() bool { return true }
