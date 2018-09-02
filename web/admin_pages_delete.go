package web

import (
	"net/http"

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
}

func (apdh *AdminPagesNewHandler) Route() string { return apdh.route }

func (apdh *AdminPagesNewHandler) HandlesGet() bool  { return false }
func (apdh *AdminPagesNewHandler) HandlesPost() bool { return true }
