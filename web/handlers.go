package web

import (
	"net/http"
	"time"
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
		&IndexHandler{
			route:  "/",
			Router: router,
		},
		&LoginHandler{
			route:  "/login",
			Router: router,
		},
		&LogoutHandler{
			route:  "/logout",
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

func UnixToTimeString(unix int64) string {
	return time.Unix(unix, 0).Format("15:04:05 02-01-2006")
}
