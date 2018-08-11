package web

import (
	"net/http"
)

type Handler interface {
	Route() string
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	HandlesGet() bool
	HandlesPost() bool
}

func GetDefaultHandlers(router *MutableRouter) []Handler {
	return []Handler{
		&LoginHandler{
			route:  "/login",
			Router: router,
		},
		&AdminHandler{
			route:  "/admin",
			Router: router,
		},
		&AdminUsersHandler{
			route:  "/admin/users",
			Router: router,
		},
		&AdminPagesHandler{
			route:  "/admin/pages",
			Router: router,
		},
	}
}
