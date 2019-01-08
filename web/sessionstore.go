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
	"time"

	"github.com/tacusci/logging"

	"github.com/gorilla/sessions"
	"github.com/tacusci/berrycms/db"
)

var (
	// need to change this to be random per db somehow
	sessionStoreSecretKey = []byte("83fdjuif49f4fjdim93490cvk4gkirv349")
	sessionsstore         = sessions.NewCookieStore(sessionStoreSecretKey)
)

func init() {
	sessionsstore.Options = &sessions.Options{
		HttpOnly: true,
	}
}

//ClearOldSessions start checking every 10 seconds for existing sessions older than 20 minutes
func ClearOldSessions(stop *chan bool) {
	startTime := time.Now()
	authSessionsTable := db.AuthSessionsTable{}
	for {
		time.Sleep(5 * time.Millisecond)
		select {
		case <-*stop:
			return
		default:
			if time.Since(startTime).Seconds() > 10 {
				err := authSessionsTable.Delete(db.Conn, fmt.Sprintf("lastactivedatetime + %d <= %d", 60*20, time.Now().Unix()))

				if err != nil {
					logging.Error(err.Error())
				}

				startTime = time.Now()
			}
		}
	}
}
