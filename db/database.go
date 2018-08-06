package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tacusci/logging"

	//blank import to make sure right SQL driver is used to talk to DB
	_ "github.com/go-sql-driver/mysql"
)

var SchemaName string
var Conn *sql.DB
var AppName string

//Connect connects to database
func Connect(appName string, sqlDriver string, dbRoute string, schemaName string) {
	AppName = appName
	SchemaName = schemaName
	db, err := sql.Open(sqlDriver, dbRoute+SchemaName)
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
	usersTable := &UsersTable{}
	usersTable.Insert(Conn, User{
		FirstName:   "John",
		LastName:    "Doe",
		Username:    "jdoe",
		AuthHash:    util.HashAndSalt([]byte("iamjohndoe")),
		Email:       "person@place.com",
		PhoneNumber: "0449488484",
	})
}

func Heartbeat() {
	<-time.After(time.Second * 60)
	Conn.Ping()
}

//Wipe drops all database tables
func Wipe() error {
	tablesToDrop := getTables()
	for i := range tablesToDrop {
		tableToDrop := tablesToDrop[i]
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
	for i := range tablesToCreate {
		tableToCreate := tablesToCreate[i]
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
	return []Table{&UsersTable{}, &AuthTable{}, &UserRolesTable{}}
}
