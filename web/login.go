// Copyright (c) 2019 tacusci ltd
//
// Licensed under the GNU GENERAL PUBLIC LICENSE Version 3 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/gpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/gofrs/uuid"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//LoginHandler handler to contain pointer to core router and the URI string
type LoginHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (lh *LoginHandler) Get(w http.ResponseWriter, r *http.Request) {

	amw := AuthMiddleware{}

	if !amw.IsLoggedIn(r) {

		ut := db.UsersTable{}
		if !ut.RootUserExists() {
			http.Redirect(w, r, "/admin/users/root/new", http.StatusFound)
		}

		pctx := plush.NewContext()

		pctx.Set("formname", "loginform")
		pctx.Set("title", "Dashboard Login")
		pctx.Set("quillenabled", false)
		pctx.Set("formhash", lh.mapFormToHash(w, r, "loginform"))
		pctx.Set("loginerrormessage", "")
		pctx.Set("adminhiddenpassword", "")
		if lh.Router.AdminHidden {
			pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", lh.Router.AdminHiddenPassword))
		}

		loginErrorStore, err := sessionsstore.Get(r, "passerrmsg")

		if err != nil {
			Error(w, err)
			return
		}

		if loginErrorMessage := loginErrorStore.Values["errormessage"]; loginErrorMessage != nil && loginErrorMessage != "" {
			logging.Debug(fmt.Sprintf("Login attempt error message: '%s'", loginErrorMessage))
			pctx.Set("loginerrormessage", loginErrorMessage)
			loginErrorStore.Values["errormessage"] = ""
			loginErrorStore.Save(r, w)
		}

		RenderDefault(w, "login.html", pctx)
	} else {
		var redirectURI = "/admin"

		if lh.Router.AdminHidden {
			redirectURI = fmt.Sprintf("/%s", lh.Router.AdminHiddenPassword) + redirectURI
		}
		http.Redirect(w, r, redirectURI, http.StatusFound)
	}
}

//Post handles post requests to URI
func (lh *LoginHandler) Post(w http.ResponseWriter, r *http.Request) {

	logging.Debug("Recieved login form POST submission...")

	err := r.ParseForm()

	if err != nil {
		Error(w, err)
		return
	}

	if lh.fetchFormHash(w, r, r.PostFormValue("formname")) == r.PostFormValue("hashid") {

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

			if authSession, err := authSessionsTable.SelectByUserUUID(db.Conn, user.UUID); err != nil {
				logging.Debug(fmt.Sprintf("There's no existing session uuid for user: %s of UUID: %s, creating session of UUID: %s...", user.Username, user.UUID, sessionUUID))
				err := authSessionsTable.Insert(db.Conn, &db.AuthSession{
					CreatedDateTime:    time.Now().Unix(),
					LastActiveDateTime: time.Now().Unix(),
					SessionUUID:        sessionUUID,
					UserUUID:           user.UUID,
				})

				if err != nil {
					Error(w, err)
				}
			} else {
				logging.Debug(fmt.Sprintf("Existing session for uuid for user: %s of UUID: %s, updating...", user.Username, user.UUID))
				err := authSessionsTable.Update(db.Conn, &db.AuthSession{
					CreatedDateTime:    authSession.CreatedDateTime,
					LastActiveDateTime: time.Now().Unix(),
					SessionUUID:        sessionUUID,
					UserUUID:           user.UUID,
				})
				if err != nil {
					Error(w, err)
				}
			}

			authSessionStore, err := sessionsstore.Get(r, "auth")

			if err != nil {
				Error(w, err)
			}

			authSessionStore.Values["sessionuuid"] = sessionUUID
			authSessionStore.Save(r, w)

			logging.Debug("Updated session store with new session UUID and added created date/timestamp")
		} else {
			authSessionStore, err := sessionsstore.Get(r, "auth")

			if err != nil {
				Error(w, err)
			}

			logging.Debug("Login unsuccessful...")
			authSessionStore.Values["sessionuuid"] = ""
			authSessionStore.Options.MaxAge = -1

			authSessionStore.Save(r, w)

			loginErrorStore, err := sessionsstore.Get(r, "passerrmsg")

			if err != nil {
				Error(w, err)
			}

			loginErrorStore.Values["errormessage"] = "Username or password incorrect..."
			loginErrorStore.Save(r, w)
		}
	} else {
		logging.Error("Login form submitted with invalid uuid hash")
	}

	http.Redirect(w, r, lh.route, http.StatusFound)
}

func (lh *LoginHandler) mapFormToHash(w http.ResponseWriter, r *http.Request, formName string) string {
	formSessionStore, err := sessionsstore.Get(r, "forms")
	defer formSessionStore.Save(r, w)

	if err != nil {
		Error(w, err)
	}

	newUUID, err := uuid.NewV4()

	if err != nil {
		Error(w, err)
	}

	formUUID := newUUID.String()
	formSessionStore.Values[formName] = formUUID

	return formUUID
}

func (lh *LoginHandler) fetchFormHash(w http.ResponseWriter, r *http.Request, formName string) string {
	formSessionStore, err := sessionsstore.Get(r, "forms")
	defer formSessionStore.Save(r, w)

	if err != nil {
		Error(w, err)
	}

	var formUUID string
	if formSessionStore.Values[formName] != nil {
		formUUID = formSessionStore.Values[formName].(string)
	}

	return formUUID
}

//Route get URI route for handler
func (lh *LoginHandler) Route() string { return lh.route }

//HandlesGet retrieve whether this handler handles get requests
func (lh *LoginHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (lh *LoginHandler) HandlesPost() bool { return true }
