package web

import (
	"github.com/tacusci/berrycms/db"
	"net/http"
)

type AdminUserGroupsEditAddHandler struct {
	Router *MutableRouter
	route  string
}

func (augeah *AdminUserGroupsEditAddHandler) Get(w http.ResponseWriter, r *http.Request) {
	var users = make([]db.User, 0)
	vars := mux.Vars(r)

	ut := db.UsersTable{}
	rows, err := ut.Select(db.Conn, "uuid, username")
}
