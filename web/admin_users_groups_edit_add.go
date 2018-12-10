package web

import (
	"github.com/gorilla/mux"
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
	rows, err := ut.Select(db.Conn, "uuid, username", vars["uuid"])

	if err != nil {
		Error(w, err)
	}

	defer rows.Close()

	for rows.Next() {
		user := db.User{}
		rows.Scan(&user.UUID, &user.Username)

		users = append(users, user)
	}
}
