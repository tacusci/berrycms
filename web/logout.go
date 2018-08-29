package web

import (
	"errors"
	"net/http"

	"github.com/tacusci/berrycms/db"
)

//LogoutHandler contains response functions for logout
type LogoutHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles HTTP get requests for logout route
func (lh *LogoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := logout(w, r)
	if err != nil {
		Error(w, err)
	}
}

//Post handles HTTP post requests for logout route
func (lh *LogoutHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := logout(w, r)
	if err != nil {
		Error(w, err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) error {
	authSessionsTable := db.AuthSessionsTable{}

	authSessionStore, err := sessionsstore.Get(r, "auth")

	if err != nil {
		return err
	}

	sessionUUID := authSessionStore.Values["sessionuuid"]

	if sessionUUID != nil && sessionUUID.(string) != "" {
		if err := authSessionsTable.DeleteBySessionUUID(db.Conn, sessionUUID.(string)); err == nil {
			authSessionStore.Values["sessionuuid"] = ""
			authSessionStore.Options.MaxAge = -1
			authSessionStore.Save(r, w)
		} else {
			return err
		}
	} else {
		return errors.New("Unable to read existing session UUID from cookie store")
	}

	http.Redirect(w, r, "/", http.StatusFound)

	return nil
}

func (lh *LogoutHandler) Route() string { return lh.route }

func (lh *LogoutHandler) HandlesGet() bool  { return false }
func (lh *LogoutHandler) HandlesPost() bool { return true }
