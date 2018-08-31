package web

import (
	"fmt"
	"net/http"

	"github.com/tacusci/logging"

	"github.com/gorilla/mux"
	"github.com/tacusci/berrycms/db"

	"github.com/gobuffalo/plush"
)

//AdminPagesHandler contains response functions for pages admin page
type AdminPagesEditHandler struct {
	Router *MutableRouter
	route  string
}

//Get takes the web request and writes response to session
func (apeh *AdminPagesEditHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pt := db.PagesTable{}
	pageToEdit, err := pt.SelectByUUID(db.Conn, vars["uuid"])
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("Page to edit not found"))
		return
	}
	pctx := plush.NewContext()
	pctx.Set("title", fmt.Sprintf("Edit Page - %s", pageToEdit.Title))
	pctx.Set("submitroute", r.RequestURI)
	pctx.Set("pagetitle", pageToEdit.Title)
	pctx.Set("pageroute", pageToEdit.Route)
	pctx.Set("pagecontent", pageToEdit.Content)
	pctx.Set("quillenabled", true)
	RenderDefault(w, "admin.pages.edit.html", pctx)
}

func (apeh *AdminPagesEditHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, r.RequestURI, http.StatusFound)
	vars := mux.Vars(r)
	pt := db.PagesTable{}
	pageToEdit, err := pt.SelectByUUID(db.Conn, vars["uuid"])
	if err != nil {
		logging.Error(err.Error())
		return
	}

	err = r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		return
	}

	pageToEdit.Title = r.PostFormValue("title")
	pageToEdit.Route = r.PostFormValue("route")
	pageToEdit.Content = r.PostFormValue("pagecontent")

	err = pt.Update(db.Conn, pageToEdit)

	if err != nil {
		logging.Error(err.Error())
	}
}

func (apeh *AdminPagesEditHandler) Route() string { return apeh.route }

func (apeh *AdminPagesEditHandler) HandlesGet() bool  { return true }
func (apeh *AdminPagesEditHandler) HandlesPost() bool { return true }
