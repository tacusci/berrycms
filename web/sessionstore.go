package web

import (
	"fmt"
	"time"

	"github.com/tacusci/logging"

	"github.com/gorilla/sessions"
	"github.com/tacusci/berrycms/db"
)

var (
	sessionStoreSecretKey = []byte("83fdjuif49f4fjdim93490cvk4gkirv349")
	sessionsstore         = sessions.NewCookieStore(sessionStoreSecretKey)
)

func init() {
	sessionsstore.Options = &sessions.Options{
		HttpOnly: true,
	}
}

func ClearOldSessions(stop *chan bool) {
	startTime := time.Now()
	authSessionsTable := db.AuthSessionsTable{}
	for {
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
