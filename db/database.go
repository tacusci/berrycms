// Copyright (c) 2018, tacusci ltd
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

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tacusci/berrycms/util"

	"github.com/tacusci/logging"

	//blank import to make sure right SQL driver is used to talk to DB
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var SchemaName string
var Conn *sql.DB
var Type DBType

const (
	dbFileName string = "./berrycms.db"

	MySQL  DBType = iota
	SQLITE DBType = iota
)

type DBType int

func (dt *DBType) DriverName() string {
	if *dt == MySQL {
		return "mysql"
	} else if *dt == SQLITE {
		return "sqlite3"
	}
	return ""
}

//Connect connects to database
func Connect(dbType DBType, dbRoute string, schemaName string) {
	SchemaName = schemaName
	Type = dbType
	var dbLoc string
	switch dbType {
	case MySQL:
		dbLoc = dbRoute + SchemaName
	case SQLITE:
		dbLoc = dbRoute
		if dbRoute == "" {
			dbLoc = dbFileName
		}
	}
	logging.InfoNnl(fmt.Sprintf("Connecting to %s:%s schema...", Type.DriverName(), dbLoc))
	db, err := sql.Open(Type.DriverName(), dbLoc)
	if err != nil {
		logging.ErrorNnl(fmt.Sprintf(" DB error: %s\n", err.Error()))
	}
	err = db.Ping()
	if err != nil {
		logging.ErrorAndExit((fmt.Sprintf(" Error connecting to DB: %s", err.Error())))
		return
	}
	logging.GreenOutput(" Connected...\n")
	Conn = db
}

func Close() {
	if Conn != nil {
		Conn.Close()
	}
}

//CreateTestData fill database with known test data for development/testing purposes
func CreateTestData() {
	logging.Debug("Creating test user...")
	usersTable := UsersTable{}
	err := usersTable.Insert(Conn, &User{
		Username:        "jdoe",
		CreatedDateTime: time.Now().Unix(),
		FirstName:       "John",
		LastName:        "Doe",
		AuthHash:        util.HashAndSalt([]byte("iamjohndoe")),
		Email:           "person@place.com",
	})
	rootUser := User{}
	rows, err := usersTable.Select(Conn, "UUID", fmt.Sprintf("userid = 1 AND userroleid = %d", ROOT_USER))
	if err != nil {
		logging.Error("Unable to fetch root user...")
		return
	}
	for rows.Next() {
		rows.Scan(&rootUser.UUID)
	}
	logging.Debug("Creating test page...")
	pagesTable := PagesTable{}
	err = pagesTable.Insert(Conn, &Page{
		CreatedDateTime: time.Now().Unix(),
		AuthorUUID:      rootUser.UUID,
		Title:           "Add New",
		Route:           "/addnew",
		Content:         "<html><body><h2>Adding Carbonite page...</h2></body></html>",
	})
	if err != nil {
		logging.Error(err.Error())
	}
}

func Heartbeat() {
	for {
		time.Sleep(time.Second * 60)
		err := Conn.Ping()
		if err != nil {
			logging.Error(fmt.Sprintf("DB Ping error -> %s", err.Error()))
		}
	}
}

//Wipe drops all database tables
func Wipe() error {
	logging.Info("Wiping DB...")
	for _, tableToDrop := range getTables() {
		logging.Debug(fmt.Sprintf("Dropping %s table...", tableToDrop.Name()))
		dropSmt := fmt.Sprintf("DROP TABLE %s;", tableToDrop.Name())
		_, err := Conn.Exec(dropSmt)
		if err != nil {
			return err
		}
	}
	return nil
}

//Setup constructs all the tables etc.,
func Setup() {
	if Conn == nil {
		return
	}
	logging.Info("Setting up DB...")
	createTables(Conn)
}

func createTables(db *sql.DB) {
	logging.Debug("Creating all database tables...")
	tablesToCreate := getTables()
	for _, tableToCreate := range tablesToCreate {
		tableCreateStatement := createStatement(tableToCreate)

		logging.Debug(fmt.Sprintf("Creating table %s...", tableToCreate.Name()))
		logging.Debug(fmt.Sprintf("Running create statement: \"%s\"", tableCreateStatement))

		_, err := db.Exec(tableCreateStatement)
		tableToCreate.Init(db)

		if err != nil {
			logging.Error(err.Error())
		}
	}
}

func getTables() []Table {
	return []Table{&SystemInfoTable{}, &UsersTable{}, &PagesTable{}, &AuthSessionsTable{}}
}
