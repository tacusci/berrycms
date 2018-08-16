package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/uuid"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//LoginHandler contains response functions for admin login
type LoginHandler struct {
	Router *MutableRouter
	route  string
}

//Get takes the web request and writes response to session
func (lh *LoginHandler) Get(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "login.html")
	if err != nil {
		logging.Error("Unable to find resources folder...")
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	pctx := plush.NewContext()
	pctx.Set("formname", "loginform")
	pctx.Set("formhash", "fjei4ijiorgrig4ijio34ofj4ig34i4j3i")

	renderedContent, err := plush.Render(string(content), pctx)
	if err != nil {
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write([]byte(renderedContent))
}

func (lh *LoginHandler) Post(w http.ResponseWriter, r *http.Request) {

	logging.Debug("Recieved login form POST submission...")

	err := r.ParseForm()

	if err != nil {
		Error(w, err)
		return
	}

	ut := db.UsersTable{}
	user, err := ut.SelectByUsername(db.Conn, r.PostFormValue("username"))

	if err != nil {
		Error(w, err)
		return
	}

	user.AuthHash = r.PostFormValue("authhash")

	if user.Login() {
		logging.Debug("Login successful...")

		v4UUID, err := uuid.NewV4()

		if err != nil {
			Error(w, err)
			return
		}

		sessionUUID := v4UUID.String()

		authSessionsTable := db.AuthSessionsTable{}

		if _, err := authSessionsTable.SelectByUserUUID(db.Conn, user.UUID); err == nil {
			logging.Debug(fmt.Sprintf("There's no existing session uuid for user: %s of UUID: %s, creating session of UUID: %s...", user.Username, user.UUID, sessionUUID))
			err := authSessionsTable.Insert(db.Conn, db.AuthSession{
				SessionUUID: sessionUUID,
				UserUUID:    user.UUID,
			})

			if err != nil {
				Error(w, err)
				return
			}
		} else {
			logging.Debug(fmt.Sprintf("Existing session for uuid for user: %s of UUID: %s, updating...", user.Username, user.UUID))
			err := authSessionsTable.Update(db.Conn, db.AuthSession{
				SessionUUID: sessionUUID,
				UserUUID:    user.UUID,
			})
			if err != nil {
				Error(w, err)
				return
			}
		}

		authSessionStore, err := sessionsstore.Get(r, "auth")

		if err != nil {
			Error(w, err)
			return
		}

		authSessionStore.Values["sessionuuid"] = sessionUUID
		authSessionStore.Save(r, w)

		logging.Debug("Updated session store with new session UUID and added created date/timestamp")
	} else {
		authSessionStore, err := sessionsstore.Get(r, "auth")

		if err != nil {
			Error(w, err)
			return
		}

		logging.Debug("Login unsuccessful...")
		authSessionStore.Values["sessionuuid"] = ""

		authSessionStore.Save(r, w)
	}

	http.Redirect(w, r, lh.route, http.StatusFound)
}

func (lh *LoginHandler) Route() string { return lh.route }

func (lh *LoginHandler) HandlesGet() bool  { return true }
func (lh *LoginHandler) HandlesPost() bool { return true }
