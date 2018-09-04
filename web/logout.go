package web

import (
	"errors"
	"net/http"

	"github.com/tacusci/berrycms/db"
)

//LogoutHandler handler to contain pointer to core router and the URI string
type LogoutHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (lh *LogoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := logout(w, r)
	if err != nil {
		Error(w, err)
	}
}

//Post handles post requests to URI
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

//Route get URI route for handler
func (lh *LogoutHandler) Route() string { return lh.route }

//HandlesGet retrieve whether this handler handles get requests
func (lh *LogoutHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (lh *LogoutHandler) HandlesPost() bool { return true }
