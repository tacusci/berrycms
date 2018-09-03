package web

import (
	"net/http"

	"github.com/tacusci/berrycms/db"

	"github.com/tacusci/logging"
)

type AdminPagesDeleteHandler struct {
	Router *MutableRouter
	route  string
}

func (apdh *AdminPagesDeleteHandler) Get(w http.ResponseWriter, r *http.Request) {}
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

func (apdh *AdminPagesDeleteHandler) Route() string { return apdh.route }

func (apdh *AdminPagesDeleteHandler) HandlesGet() bool  { return false }
func (apdh *AdminPagesDeleteHandler) HandlesPost() bool { return true }
