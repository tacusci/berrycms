package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/tacusci/berrycms/util"
	"github.com/tacusci/logging"
)

type UsersRoleFlag int

const (
	ROOT UsersRoleFlag = 2
)

//Field interface to describe a table field and all of its attributes
type Field struct {
	fieldTag      reflect.StructTag
	kind          reflect.Kind
	AutoIncrement bool
	PrimaryKey    bool
	UniqueIndex   bool
	NotNull       bool
	Name          string
	Type          string
	Value         interface{}
}

func (f *Field) parseFlagTags() {
	fieldTagString := f.fieldTag.Get("tbl")

	if strings.Contains(fieldTagString, "PK") {
		f.PrimaryKey = true
	}

	if strings.Contains(fieldTagString, "NN") {
		f.NotNull = true
	}

	if strings.Contains(fieldTagString, "AI") {
		f.AutoIncrement = true
	}

	if strings.Contains(fieldTagString, "UI") {
		f.UniqueIndex = true
	}
}

func (f *Field) translateTypes() {
	switch f.Type {
	case "string":
		f.Type = "VARCHAR(125)"
	case "bool":
		//f.Type = "BOOLEAN"
		f.Type = "BIT(1)"
	case "uint32":
		f.Type = "LONG"
	}
	f.Type = strings.ToUpper(f.Type)
}

func (f *Field) getFormatString() string {
	switch f.kind {
	case reflect.Bool:
		return "%t"
	case reflect.Int:
		return "%d"
	default:
		return "%s"
	}
}

// ****************************************** TABLES ******************************************

//Table interface to inherit from all table structs
type Table interface {
	Init(db *sql.DB)
	Name() string
	buildFields() []Field
	buildInsertStatement(m Model) string
}

// ******** UserTable ********

//UsersTable describes the table structure for UsersTable in db
type UsersTable struct {
	Userid     int    `tbl:"PKNNAIUI"`
	Userroleid int    `tbl:"NN"`
	UUID       string `tbl:"NNUI"`
	Username   string `tbl:"NNUI"`
	Authhash   string `tbl:"NN"`
	Firstname  string `tbl:"NN"`
	Lastname   string `tbl:"NN"`
	Email      string `tbl:"NNUI"`
}

//Init carries out default data entry
func (ut *UsersTable) Init(db *sql.DB) {
	resultRows, err := ut.Select(db, "*", "userid = 1")
	if err == nil {
		if !resultRows.Next() {
			logging.Debug("Creating default root user account...")
			err = ut.Insert(db, User{
				Username:   "root",
				UserroleId: int(ROOT),
				AuthHash:   util.HashAndSalt([]byte("iamroot")),
				FirstName:  "Root",
				LastName:   "User",
				Email:      "none",
			})
			if err != nil {
				logging.ErrorAndExit(err.Error())
			}
		} else {
			logging.Debug("Root user already exists... Cannot re-create.")
		}
	} else {
		logging.Error(err.Error())
	}
}

//Name gets the table name, have to implement to make UsersTable inherit Table
func (ut *UsersTable) Name() string {
	return "users"
}

//InsertMultiple takes a slice of user structs and passes them all to 'Insert'
func (ut *UsersTable) InsertMultiple(db *sql.DB, us []User) error {
	var err error
	for i := range us {
		err = ut.Insert(db, us[i])
	}
	return err
}

//Insert adds user struct to users table, it also sets default values
func (ut *UsersTable) Insert(db *sql.DB, u User) error {

	if u.UUID != "" {
		return fmt.Errorf("User to insert already has UUID %s", u.UUID)
	}

	err := u.Validate()

	if err != nil {
		return err
	}

	if u.UUID == "" {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		u.UUID = newUUID.String()
		if u.UserroleId == 0 {
			u.UserroleId = 3
		}
		insertStatement := ut.buildInsertStatement(&u)
		_, err = db.Exec(insertStatement)
		if err != nil {
			return err
		}
	}

	return nil
}

//Select returns table rows from a select using the passed where condition
func (ut *UsersTable) Select(db *sql.DB, whatToSelect string, whereClause string) (*sql.Rows, error) {
	if len(whereClause) > 0 {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s", whatToSelect, SchemaName, ut.Name(), whereClause))
	} else {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s.%s", whatToSelect, SchemaName, ut.Name()))
	}
}

//BuildFields takes the table struct and maps all of the struct fields to their own struct
func (ut *UsersTable) buildFields() []Field {
	return buildFieldsFromTable(ut)
}

func (ut *UsersTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(ut, m)
}

// ******** End UserTable ********

// ******** Start User Roles Table ********

//UserRolesTable describes the table structure for UserRolesTable in db
type UserRolesTable struct {
	Userrolesid int    `tbl:"PKNNAIUI"`
	Rolename    string `tbl:"NN"`
}

//Init initialises the UserRolesTable table with default data
func (urt *UserRolesTable) Init(db *sql.DB) {
	urt.Insert(db, UserRole{
		Rolename: "super",
	})
	urt.Insert(db, UserRole{
		Rolename: "admin",
	})
	urt.Insert(db, UserRole{
		Rolename: "moderator",
	})
	urt.Insert(db, UserRole{
		Rolename: "guest",
	})
}

//Name gets the table name, have to implement to make UserRolesTable inherit Table
func (urt *UserRolesTable) Name() string {
	return "userroles"
}

//Insert inserts parsed UserRole model into the user role table
func (urt *UserRolesTable) Insert(db *sql.DB, ur UserRole) error {
	if ur.Rolename != "" {
		insertStatement := urt.buildInsertStatement(&ur)
		_, err := db.Exec(insertStatement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (urt *UserRolesTable) buildFields() []Field {
	return buildFieldsFromTable(urt)
}

func (urt *UserRolesTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(urt, m)
}

// ******** End User Roles Table ********

// ******** Start Pages Table ********
type PagesTable struct {
	Pageid  int    `tbl:"PKNNAIUI"`
	UUID    string `tbl:"NNUI"`
	Title   string `tbl:"NNUI"`
	Route   string `tbl:"NNUI"`
	Content string `tbl:"NN"`
}

func (pt *PagesTable) Init(db *sql.DB) {}

func (pt *PagesTable) Name() string {
	return "pages"
}

func (pt *PagesTable) Insert(db *sql.DB, p Page) error {
	if p.UUID != "" {
		return fmt.Errorf("Page to insert already has UUID %s", p.UUID)
	}

	if p.UUID == "" {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		p.UUID = newUUID.String()
		insertStatement := pt.buildInsertStatement(&p)
		_, err = db.Exec(insertStatement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pt *PagesTable) Select(db *sql.DB, whatToSelect string, whereClause string) (*sql.Rows, error) {
	if len(whereClause) > 0 {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s", whatToSelect, SchemaName, pt.Name(), whereClause))
	} else {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s.%s", whatToSelect, SchemaName, pt.Name()))
	}
}

func (pt *PagesTable) buildFields() []Field {
	return buildFieldsFromTable(pt)
}

func (pt *PagesTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(pt, m)
}

// ******** End Pages Table ********

// ****************************************** END TABLES ******************************************
/////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////
//////////////////////////////////////////////
//////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// ****************************************** MODELS ******************************************

//Model describes the structure of a model
type Model interface {
	TableName() string
	BuildFields() []Field
}

//User describes the content of a user, it should match the columns present in the users table
type User struct {
	UserId     int    `tbl:"AI" json:"userid"`
	UserroleId int    `json:"userroleid"`
	UUID       string `json:"UUID"`
	Username   string `json:"username"`
	AuthHash   string `json:"authhash"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
}

//Login takes the current username and authhash values of self and tries
//using them to authenticate/login. A successful login will return/generate
//a JWT token for further use in any subsequent API request
func (u *User) Login() bool {
	p := u.AuthHash
	ut := &UsersTable{}
	row, err := ut.Select(Conn, "userid, userroleid, uuid, username, authhash", fmt.Sprintf("username = '%s'", u.Username))
	if err != nil {
		logging.Error(err.Error())
		return false
	}
	defer row.Close()
	for row.Next() {
		u := &User{}
		err := row.Scan(&u.UserId, &u.UserroleId,
			&u.UUID, &u.Username, &u.AuthHash)
		if err != nil {
			logging.Error(err.Error())
			return false
		}
		logging.Info(fmt.Sprintf("UserID: %d", u.UserId))
		err = bcrypt.CompareHashAndPassword([]byte(u.AuthHash), []byte(p))
		if err == nil {
			return true
		}
	}
	return false
}

//TableName gets the name of the users table
func (u *User) TableName() string {
	return "users"
}

//BuildFields generates a list of fields generated from the fields of the user struct
func (u *User) BuildFields() []Field {
	return buildFieldsFromModel(u)
}

//Validate makes sure that required fields have not been left blank
func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("Missing username")
	}
	if u.AuthHash == "" {
		return fmt.Errorf("Missing password")
	}
	return nil
}

//UserRole describes the content of a userrole entry, it should match the columns present in the userrole table
type UserRole struct {
	Userroleid int `tbl:"AI"`
	Rolename   string
}

//TableName gets the name of the userrole table
func (ur *UserRole) TableName() string {
	return "userroles"
}

//BuildFields generates a list of fields generated from the fields of the userrole struct
func (ur *UserRole) BuildFields() []Field {
	return buildFieldsFromModel(ur)
}

type Page struct {
	PageId  int    `tbl:"AI" json:"pageid"`
	UUID    string `json:"UUID"`
	Title   string `json:"title"`
	Route   string `json:"route"`
	Content string `json:"content"`
}

func (p *Page) TableName() string {
	return "pages"
}

func (p *Page) BuildFields() []Field {
	return buildFieldsFromModel(p)
}

// ****************************************** END MODELS ******************************************

func buildInsertStatementFromTable(t Table, m Model) string {
	var insertStatementBuilder bytes.Buffer
	insertStatementBuilder.WriteString(fmt.Sprintf("INSERT INTO %s (", t.Name()))

	tableFields := t.buildFields()
	tableFieldsCount := len(tableFields)

	for i := 0; i < tableFieldsCount; i++ {
		tableField := tableFields[i]
		if !tableField.AutoIncrement && tableField.Name != "" {
			insertStatementBuilder.WriteString(fmt.Sprintf("%s", tableField.Name))
			if i+1 < tableFieldsCount {
				insertStatementBuilder.WriteString(", ")
			}
		}
	}
	insertStatementBuilder.WriteString(") VALUES (")

	modelFields := m.BuildFields()
	modelFieldsCount := len(modelFields)

	for i := 0; i < modelFieldsCount; i++ {
		modelField := modelFields[i]
		if !modelField.AutoIncrement {
			formatString := modelField.getFormatString()
			if modelField.Type != "boolean" && modelField.Type != "BIT(1)" {
				formatString = fmt.Sprintf("'%s'", formatString)
			}
			insertStatementBuilder.WriteString(fmt.Sprintf(formatString, modelField.Value))
			if i+1 < modelFieldsCount {
				insertStatementBuilder.WriteString(", ")
			}
		}
	}
	insertStatementBuilder.WriteString(")")

	return insertStatementBuilder.String()
}

func buildFieldsFromModel(m Model) []Field {
	fields := make([]Field, 0)

	val := reflect.ValueOf(m).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		newField := Field{
			kind:     valueField.Kind(),
			fieldTag: tag,
			Name:     typeField.Name,
			Type:     typeField.Type.String(),
			Value:    reflect.Value(valueField),
		}
		newField.parseFlagTags()
		newField.translateTypes()
		fields = append(fields, newField)
	}

	return fields
}

//using reflection to map the model struct to a create statement
func buildFieldsFromTable(t Table) []Field {
	fields := make([]Field, 0)

	tableStructValue := reflect.ValueOf(t).Elem()
	tableStructType := tableStructValue.Type()

	for i := 0; i < tableStructValue.NumField(); i++ {
		//get the field
		tableStructField := tableStructValue.Field(i)
		newField := Field{
			kind:     tableStructField.Kind(),
			fieldTag: tableStructType.Field(i).Tag,
			Name:     strings.ToLower(tableStructType.Field(i).Name),
			Type:     tableStructField.Type().String(),
		}
		newField.parseFlagTags()
		newField.translateTypes()
		fields = append(fields, newField)
	}
	return fields
}

func createStatement(t Table) string {
	var stringBulder bytes.Buffer
	stringBulder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`.`%s` (", SchemaName, t.Name()))

	tableFields := t.buildFields()
	tableFieldsCount := len(tableFields)

	var pkField Field
	pkFieldCount := 0

	var uniqueIndexFields []Field
	uniqueIndexFieldsCount := 0

	for j := 0; j < tableFieldsCount; j++ {
		field := tableFields[j]
		stringBulder.WriteString(fmt.Sprintf("`%s` %s", field.Name, field.Type))
		if field.PrimaryKey {
			pkFieldCount++
			if pkFieldCount > 1 {
				logging.ErrorAndExit(fmt.Sprintf("Error creating %s table: More than one PK field found...", t.Name()))
			}
			pkField = field
		}
		if field.AutoIncrement {
			stringBulder.WriteString(" AUTO_INCREMENT")
		}
		if field.NotNull {
			stringBulder.WriteString(" NOT NULL")
		}
		if field.UniqueIndex {
			uniqueIndexFields = append(uniqueIndexFields, field)
			uniqueIndexFieldsCount++
		}
		if j+1 < tableFieldsCount || pkFieldCount > 0 {
			stringBulder.WriteString(",")
		}
	}

	if pkFieldCount == 1 {
		stringBulder.WriteString(fmt.Sprintf(" PRIMARY KEY (`%s`)", pkField.Name))
	}

	if len(uniqueIndexFields) > 0 {
		stringBulder.WriteString(",")
	}

	for i := 0; i < uniqueIndexFieldsCount; i++ {
		uniqueIndexField := uniqueIndexFields[i]
		stringBulder.WriteString(fmt.Sprintf(" UNIQUE INDEX `%s_UNIQUE` (`%s` ASC)", uniqueIndexField.Name, uniqueIndexField.Name))
		if i+1 < uniqueIndexFieldsCount {
			stringBulder.WriteString(",")
		}
	}

	stringBulder.WriteString(");")
	return stringBulder.String()
}
