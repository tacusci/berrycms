package web

import (
	"net/http"

	"github.com/gobuffalo/plush"
)

type AdminUsersNewHandler struct {
	Router *MutableRouter
	route  string
}

func (aunh *AdminUsersNewHandler) Get(w http.ResponseWriter, r *http.Request) {
	pctx := plush.NewContext()
	pctx.Set("title", "New User")
	RenderDefault(w, "admin.users.new.html", pctx)
}
func (aunh *AdminUsersNewHandler) Post(w http.ResponseWriter, r *http.Request) {

}

func (aunh *AdminUsersNewHandler) Route() string { return aunh.route }

func (aunh *AdminUsersNewHandler) HandlesGet() bool  { return true }
func (aunh *AdminUsersNewHandler) HandlesPost() bool { return true }
