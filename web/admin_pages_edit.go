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
func (aph *AdminPagesEditHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	RenderDefault(w, "admin.pages.edit.html", pctx)
}

func (aph *AdminPagesEditHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (aph *AdminPagesEditHandler) Route() string { return aph.route }

func (aph *AdminPagesEditHandler) HandlesGet() bool  { return true }
func (aph *AdminPagesEditHandler) HandlesPost() bool { return false }
