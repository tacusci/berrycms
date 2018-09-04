package web

import (
	"net/http"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/logging"
)

//AdminUsersNewHandler handler to contain pointer to core router and the URI string
type AdminUsersNewHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (aunh *AdminUsersNewHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "New User")
	RenderDefault(w, "admin.users.new.html", pctx)
}

//Post handles post requests to URI
func (aunh *AdminUsersNewHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/admin/users", http.StatusFound)
	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		return
	}
}

//Route get URI route for handler
func (aunh *AdminUsersNewHandler) Route() string { return aunh.route }

//HandlesGet retrieve whether this handler handles get requests
func (aunh *AdminUsersNewHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (aunh *AdminUsersNewHandler) HandlesPost() bool { return true }
