package web

import (
	"time"

	"github.com/gorilla/sessions"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
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
				logging.Debug("Clearing old sessions, (older than 20 minutes)...")
				rows, err := authSessionsTable.Select(db.Conn, "*", "")
				if err != nil {
					for rows.Next() {
						authSession := db.AuthSession{}
						err := rows.Scan(&authSession.Authsessionid, &authSession.CreatedDateTime, &authSession.UserUUID, &authSession.SessionUUID)
						if err != nil {
							if time.Since(time.Unix(authSession.CreatedDateTime, 0)).Minutes() >= 20 {
								authSessionsTable.DeleteBySessionUUID(db.Conn, authSession.SessionUUID)
							}
						}
					}
				}
				startTime = time.Now()
			}
		}
	}
}
