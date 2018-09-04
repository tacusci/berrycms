package web

import (
	"net/http"

	"github.com/gobuffalo/plush"
)

//AdminHandler handler to contain pointer to core router and the URI string
type AdminHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (ah *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "Dashboard")
	pctx.Set("quillenabled", false)
	RenderDefault(w, "admin.html", pctx)
}

//Post handles post requests to URI
func (ah *AdminHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (ah *AdminHandler) Route() string { return ah.route }

//HandlesGet retrieve whether this handler handles get requests
func (ah *AdminHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (ah *AdminHandler) HandlesPost() bool { return false }
