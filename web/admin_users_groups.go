package web

import (
	// "github.com/tacusci/berrycms/db"
	"net/http"
)

type AdminUserGroupsHandler struct {
	Router *MutableRouter
	route  string
}

func (ugh *AdminUserGroupsHandler) Get(w http.ResponseWriter, r *http.Request) {
	// groups := make([]db.Group, 0)

	// gt := db.GroupTable{}
}
