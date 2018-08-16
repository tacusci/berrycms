package web

import "github.com/gorilla/sessions"

var (
	sessionStoreSecretKey = []byte("83fdjuif49f4fjdim93490cvk4gkirv349")
	store                 = sessions.NewCookieStore(sessionStoreSecretKey)
)

func init() {
	// store.Options = &sessions.Options{
	// 	HttpOnly: true,
	// 	MaxAge:   0,
	// 	Secure:   true,
	// 	Domain:   "localhost",
	// 	Path:     "/",
	// }
}
