// Copyright (c) 2019, tacusci ltd
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
