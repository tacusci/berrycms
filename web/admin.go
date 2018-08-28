package web

import (
	"net/http"

	"github.com/gobuffalo/plush"
)

type AdminHandler struct {
	Router *MutableRouter
	route  string
}

func (ah *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "Dashboard")
	pctx.Set("quillenabled", false)
	RenderDefault(w, "admin.html", pctx)
}

func (ah *AdminHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (ah *AdminHandler) Route() string { return ah.route }

func (ah *AdminHandler) HandlesGet() bool  { return true }
func (ah *AdminHandler) HandlesPost() bool { return false }
