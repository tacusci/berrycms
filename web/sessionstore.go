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
			if time.Since(startTime).Seconds() > 60 {
				authSessionsToRemove := make([]db.AuthSession, 0)
				rows, err := authSessionsTable.Select(db.Conn, "sessionuuid", fmt.Sprintf("createddatetime + %d =< %d", 60*20, time.Now().Unix()))

				if err != nil {
					logging.Error(err.Error())
					continue
				}

				for rows.Next() {
					authSessionToRemove := db.AuthSession{}
					err := rows.Scan(&authSessionToRemove.SessionUUID)
					if err != nil {
						logging.Error(err.Error())
						continue
					}
					authSessionsToRemove = append(authSessionsToRemove, authSessionToRemove)
				}

				rows.Close()

				for _, as := range authSessionsToRemove {
					authSessionsTable.DeleteBySessionUUID(db.Conn, as.SessionUUID)
				}

				startTime = time.Now()
			}
		}
	}
}
