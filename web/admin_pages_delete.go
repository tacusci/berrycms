package web

import (
	"net/http"

	"github.com/tacusci/berrycms/db"

	"github.com/tacusci/logging"
)

//AdminPagesDeleteHandler handler to contain pointer to core router and the URI string
type AdminPagesDeleteHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (apdh *AdminPagesDeleteHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (apdh *AdminPagesDeleteHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logging.Error(err.Error())
		Error(w, err)
	}

	pt := db.PagesTable{}
	deletedPages := false
	for _, v := range r.PostForm {
		deletedPages = true
		pt.DeleteByUUID(db.Conn, v[0])
	}

	if deletedPages {
		apdh.Router.Reload()
	}

	http.Redirect(w, r, "/admin/pages", http.StatusFound)
}

//Route get URI route for handler
func (apdh *AdminPagesDeleteHandler) Route() string { return apdh.route }

//HandlesGet retrieve whether this handler handles get requests
func (apdh *AdminPagesDeleteHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (apdh *AdminPagesDeleteHandler) HandlesPost() bool { return true }
