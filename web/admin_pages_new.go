package web

import (
	"net/http"

	"github.com/gobuffalo/plush"
)

type AdminPagesNewHandler struct {
	Router *MutableRouter
	route  string
}

//Get takes the web request and writes response to session
func (apnh *AdminPagesNewHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "New Page")
	pctx.Set("quillenabled", true)
	RenderDefault(w, "admin.pages.new.html", pctx)
}

func (apnh *AdminPagesNewHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (apnh *AdminPagesNewHandler) Route() string { return apnh.route }

func (apnh *AdminPagesNewHandler) HandlesGet() bool  { return true }
func (apnh *AdminPagesNewHandler) HandlesPost() bool { return false }
