package web

import (
	"errors"
	"net/http"

	"github.com/tacusci/berrycms/db"
)

//LoginHandler contains response functions for admin login
type LogoutHandler struct {
	Router *MutableRouter
	route  string
}

func (lh *LogoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := logout(w, r)
	if err != nil {
		Error(w, err)
	}
}

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

	return nil
}

func (lh *LogoutHandler) Route() string { return lh.route }

func (lh *LogoutHandler) HandlesGet() bool  { return false }
func (lh *LogoutHandler) HandlesPost() bool { return true }
