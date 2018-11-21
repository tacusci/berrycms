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
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/tacusci/logging"
)

type UsersRoleFlag int

const (
	ROOT_USER UsersRoleFlag = 2
	MOD_USER  UsersRoleFlag = 3
	REG_USER  UsersRoleFlag = 4
)

//Field interface to describe a table field and all of its attributes
type Field struct {
	fieldTag      reflect.StructTag
	kind          reflect.Kind
	AutoIncrement bool
	PrimaryKey    bool
	UniqueIndex   bool
	IsDateTime    bool
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

	if strings.Contains(fieldTagString, "DT") {
		f.IsDateTime = true
	}
}

func (f *Field) translateTypes() {
	switch f.Type {
	case "string":
		f.Type = "VARCHAR(125)"
	case "bool":
		//f.Type = "BOOLEAN"
		f.Type = "BIT(1)"
	case "int":
		if Type == MySQL {
			f.Type = "INT"
		} else if Type == SQLITE {
			f.Type = "INTEGER"
		}
	case "uint32":
		if Type == MySQL {
			f.Type = "INT"
		} else if Type == SQLITE {
			f.Type = "INTEGER"
		}
	case "uint64":
		if Type == MySQL {
			f.Type = "BIGINT"
		} else if Type == SQLITE {
			f.Type = "INTEGER"
		}
	}

	if f.IsDateTime {
		f.Type = "BIGINT"
	}

	f.Type = strings.ToUpper(f.Type)
}

func (f *Field) getFormatString() string {
	switch f.kind {
	case reflect.Bool:
		return "%t"
	case reflect.Int:
		return "%d"
	case reflect.Int64:
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
	Userid          int    `tbl:"PKNNAIUI"`
	CreatedDateTime int64  `tbl:"NNDT"`
	Userroleid      int    `tbl:"NN"`
	UUID            string `tbl:"NNUI"`
	Username        string `tbl:"NNUI"`
	Authhash        string `tbl:"NN"`
	Firstname       string `tbl:"NN"`
	Lastname        string `tbl:"NN"`
	Email           string `tbl:"NNUI"`
}

//Init carries out default data entry
func (ut *UsersTable) Init(db *sql.DB) {}

//Name gets the table name, have to implement to make UsersTable inherit Table
func (ut *UsersTable) Name() string { return "users" }

//RootUserExists checks if at least one root user exists
func (ut *UsersTable) RootUserExists() bool {
	rows, err := ut.Select(Conn, "userid", fmt.Sprintf("userroleid = %d", ROOT_USER))

	if err != nil {
		logging.Error(err.Error())
		return false
	}

	defer rows.Close()

	var i = 0
	for rows.Next() {
		i++
		if i > 0 {
			break
		}
	}

	return i > 0
}

//InsertMultiple takes a slice of user structs and passes them all to 'Insert'
func (ut *UsersTable) InsertMultiple(db *sql.DB, us []*User) error {
	var err error
	for i := range us {
		err = ut.Insert(db, us[i])
	}
	return err
}

//Insert adds user struct to users table, it also sets default values
func (ut *UsersTable) Insert(db *sql.DB, u *User) error {
	//TODO: change this to simply call validate()
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
		insertStatement := ut.buildPreparedInsertStatement(u)
		_, err = db.Exec(insertStatement, u.CreatedDateTime, u.UserroleId, u.UUID, u.Username, u.AuthHash, u.FirstName, u.LastName, u.Email)
		if err != nil {
			return err
		}
	}

	return nil
}

//Select returns table rows from a select using the passed where condition
func (ut *UsersTable) Select(db *sql.DB, whatToSelect string, whereClause string) (*sql.Rows, error) {
	if len(whereClause) > 0 {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s", whatToSelect, ut.Name(), whereClause))
	} else {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s", whatToSelect, ut.Name()))
	}
}

func (ut *UsersTable) SelectRootUser(db *sql.DB) (*User, error) {
	u := &User{}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE userroleid = %d", ut.Name(), int(ROOT_USER)))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&u.UserId, &u.CreatedDateTime, &u.UserroleId, &u.UUID, &u.Username, &u.AuthHash, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (ut *UsersTable) SelectByUsername(db *sql.DB, username string) (*User, error) {
	u := &User{}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE username = '%s'", ut.Name(), username))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&u.UserId, &u.CreatedDateTime, &u.UserroleId, &u.UUID, &u.Username, &u.AuthHash, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (ut *UsersTable) SelectByUUID(db *sql.DB, uuid string) (*User, error) {
	u := &User{}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE uuid = '%s'", ut.Name(), uuid))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&u.UserId, &u.CreatedDateTime, &u.UserroleId, &u.UUID, &u.Username, &u.AuthHash, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (ut *UsersTable) DeleteByUUID(db *sql.DB, uuid string) (int64, error) {
	res, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", ut.Name()), uuid)

	if err != nil {
		return 0, err
	}

	numDeleted, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return numDeleted, nil
}

//BuildFields takes the table struct and maps all of the struct fields to their own struct
func (ut *UsersTable) buildFields() []Field {
	return buildFieldsFromTable(ut)
}

func (ut *UsersTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(ut, m)
}

func (ut *UsersTable) buildPreparedInsertStatement(m Model) string {
	return buildPreparedInsertStatementFromTable(ut, m)
}

// ******** End UserTable ********

// ******** Start GroupTable ********

type GroupTable struct {
	Groupid         int    `tbl:"PKNNAIUI"`
	CreatedDateTime int64  `tbl:"NNDT"`
	UUID            string `tbl:"NNUI"`
	Title           string `tbl:"NNUI"`
}

func (gt *GroupTable) Init(db *sql.DB) {
	adminGroup := &Group{
		CreatedDateTime: time.Now().Unix(),
		Title:           "Admins",
	}
	gt.Insert(db, adminGroup)

	modGroup := &Group{
		CreatedDateTime: time.Now().Unix(),
		Title:           "Moderators",
	}
	gt.Insert(db, modGroup)

	userGroup := &Group{
		CreatedDateTime: time.Now().Unix(),
		Title:           "Users",
	}
	gt.Insert(db, userGroup)
}

func (gt *GroupTable) Name() string {
	return "groups"
}

func (gt *GroupTable) Insert(db *sql.DB, g *Group) error {
	//TODO: change this to simply call validate()
	if g.UUID != "" {
		return fmt.Errorf("Page to insert already has UUID %s", g.UUID)
	}

	if g.UUID == "" {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		g.UUID = newUUID.String()
		insertStatement := gt.buildPreparedInsertStatement(g)
		_, err = db.Exec(insertStatement, g.CreatedDateTime, g.UUID, g.Title)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gt *GroupTable) Update(db *sql.DB, g *Group) error {
	if g.Validate() {
		updateStatement := fmt.Sprintf("UPDATE %s SET createddatetime = ?, uuid = ?, title = ? WHERE uuid = ?", gt.Name())
		_, err := db.Exec(updateStatement, g.CreatedDateTime, g.Title, g.UUID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Group to insert already has UUID")
}

func (gt *GroupTable) Select(db *sql.DB, whatToSelect string, whereClause string) (*sql.Rows, error) {
	if len(whereClause) > 0 {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s", whatToSelect, gt.Name(), whereClause))
	} else {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s", whatToSelect, gt.Name()))
	}
}

func (gt *GroupTable) SelectByTitle(db *sql.DB, groupTitle string) (*Group, error) {
	g := &Group{}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE title = '%s'", gt.Name(), groupTitle))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&g.Groupid, &g.CreatedDateTime, &g.UUID, &g.Title)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (gt *GroupTable) buildFields() []Field {
	return buildFieldsFromTable(gt)
}

func (gt *GroupTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(gt, m)
}

func (gt *GroupTable) buildPreparedInsertStatement(m Model) string {
	return buildPreparedInsertStatementFromTable(gt, m)
}

// ******** End Group Table ********

// ******** Start Group Membership Table ********

type GroupMembershipTable struct {
	GroupMembershipid int    `tbl:"PKNNAIUI"`
	GroupUUID         string `tbl:"NN"`
	UserUUID          string `tbl:"NN"`
}

func (gmt *GroupMembershipTable) Init(db *sql.DB) {}

func (gmt *GroupMembershipTable) Name() string {
	return "groupmemberships"
}

func (gmt *GroupMembershipTable) Insert(db *sql.DB, gm *GroupMembership) error {
	//TODO: change this to simply call validate()
	if gm.UUID != "" {
		return fmt.Errorf("Group membership to insert already has UUID %s", gm.UUID)
	}

	if gm.UUID == "" {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		gm.UUID = newUUID.String()
		insertStatement := gmt.buildPreparedInsertStatement(gm)
		_, err = db.Exec(insertStatement, gm.CreatedDateTime, gm.UUID, gm.GroupUUID, gm.UserUUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gmt *GroupMembershipTable) Update(db *sql.DB, gm *GroupMembership) error {
	if gm.Validate() {
		updateStatement := fmt.Sprintf("UPDATE %s SET createddatetime = ?, groupuuid = ?, useruuid = ? WHERE uuid = ?", gmt.Name())
	}
}

func (gmt *GroupMembershipTable) buildFields() []Field {
	return buildFieldsFromTable(gmt)
}

func (gmt *GroupMembershipTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(gmt, m)
}

func (gmt *GroupMembershipTable) buildPreparedInsertStatement(m Model) string {
	return buildPreparedInsertStatementFromTable(gmt, m)
}

// ******** Start Pages Table ********

type PagesTable struct {
	Pageid          int    `tbl:"PKNNAIUI"`
	CreatedDateTime int64  `tbl:"NNDT"`
	UUID            string `tbl:"NNUI"`
	Roleprotected   bool   `tbl:"NN"`
	AuthorUUID      string `tbl:"NN"`
	Title           string `tbl:"NNUI"`
	Route           string `tbl:"NNUI"`
	Content         string `tbl:"NN"`
}

func (pt *PagesTable) Init(db *sql.DB) {}

func (pt *PagesTable) Name() string {
	return "pages"
}

func (pt *PagesTable) Insert(db *sql.DB, p *Page) error {
	//TODO: change this to simply call validate()
	if p.UUID != "" {
		return fmt.Errorf("Page to insert already has UUID %s", p.UUID)
	}

	if p.UUID == "" {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		p.UUID = newUUID.String()
		insertStatement := pt.buildPreparedInsertStatement(p)
		_, err = db.Exec(insertStatement, p.CreatedDateTime, p.UUID, p.Roleprotected, p.AuthorUUID, p.Title, p.Route, p.Content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pt *PagesTable) Update(db *sql.DB, p *Page) error {
	updateStatement := fmt.Sprintf("UPDATE %s SET createddatetime = ?, uuid = ?, roleprotected = ?, authoruuid = ?, title = ?, route = ?, content = ? WHERE uuid = ?", pt.Name())
	_, err := db.Exec(updateStatement, p.CreatedDateTime, p.UUID, p.Roleprotected, p.AuthorUUID, p.Title, p.Route, p.Content, p.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (pt *PagesTable) Select(db *sql.DB, whatToSelect string, whereClause string) (*sql.Rows, error) {
	if len(whereClause) > 0 {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s", whatToSelect, pt.Name(), whereClause))
	}
	return db.Query(fmt.Sprintf("SELECT %s FROM %s", whatToSelect, pt.Name()))
}

func (pt *PagesTable) SelectByRoute(db *sql.DB, route string) (*Page, error) {
	p := &Page{}

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE route = '%s'", pt.Name(), route))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&p.PageId, &p.CreatedDateTime, &p.UUID, &p.Roleprotected, &p.AuthorUUID, &p.Title, &p.Route, &p.Content)
		if err != nil {
			return nil, err
		}
	}

	//page not found, therefore don't return blank page struct
	if p.UUID == "" {
		return nil, fmt.Errorf("Page not found in table %s", pt.Name())
	}

	return p, nil
}

func (pt *PagesTable) SelectByUUID(db *sql.DB, uuid string) (*Page, error) {
	p := &Page{}
	row := db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE uuid = '%s'", pt.Name(), uuid))
	err := row.Scan(&p.PageId, &p.CreatedDateTime, &p.UUID, &p.Roleprotected, &p.AuthorUUID, &p.Title, &p.Route, &p.Content)
	if err != nil {
		return nil, err
	}
	//page not found, therefore don't return blank page struct
	if p.UUID == "" {
		return nil, fmt.Errorf("Page not found in table %s", pt.Name())
	}
	return p, nil
}

func (pt *PagesTable) DeleteByUUID(db *sql.DB, uuid string) (int64, error) {
	res, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", pt.Name()), uuid)

	if err != nil {
		return 0, err
	}

	numDeleted, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return numDeleted, nil
}

func (pt *PagesTable) buildFields() []Field {
	return buildFieldsFromTable(pt)
}

func (pt *PagesTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(pt, m)
}

func (pt *PagesTable) buildPreparedInsertStatement(m Model) string {
	return buildPreparedInsertStatementFromTable(pt, m)
}

// ******** End Pages Table ********

// ******** Start Auth Table ********

type AuthSessionsTable struct {
	Authsessionid      int    `tbl:"PKNNAIUI"`
	CreatedDateTime    int64  `tbl:"NNDT"`
	LastActiveDateTime int64  `tbl:"NNDT"`
	UserUUID           string `tbl:"NNUI"`
	SessionUUID        string `tbl:"NNUI"`
}

func (ast *AuthSessionsTable) Init(db *sql.DB) {}

func (ast *AuthSessionsTable) Name() string { return "authsessions" }

func (ast *AuthSessionsTable) Insert(db *sql.DB, as *AuthSession) error {
	if as.Validate() {
		insertStatement := ast.buildPreparedInsertStatement(as)
		_, err := db.Exec(insertStatement, as.CreatedDateTime, as.LastActiveDateTime, as.UserUUID, as.SessionUUID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("AuthSession doesn't have a user UUID and/or a session UUID")
}

//Update - Takes auth session to update existing user session entry session UUID
func (ast *AuthSessionsTable) Update(db *sql.DB, as *AuthSession) error {
	if as.Validate() {
		updateStatement := fmt.Sprintf("UPDATE %s SET createddatetime = ?, lastactivedatetime = ?, sessionuuid = ? WHERE useruuid = ?", ast.Name())
		_, err := db.Exec(updateStatement, as.CreatedDateTime, as.LastActiveDateTime, as.SessionUUID, as.UserUUID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("AuthSession doesn't have a user UUID and/or a session UUID")
}

func (ast *AuthSessionsTable) Select(db *sql.DB, whatToSelect string, whereClause string) (*sql.Rows, error) {
	if len(whereClause) > 0 {
		return db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s", whatToSelect, ast.Name(), whereClause))
	}
	return db.Query(fmt.Sprintf("SELECT %s FROM %s", whatToSelect, ast.Name()))
}

func (ast *AuthSessionsTable) SelectBySessionUUID(db *sql.DB, sessionUUID string) (*AuthSession, error) {
	as := &AuthSession{}
	row := db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE sessionuuid = '%s'", ast.Name(), sessionUUID))
	err := row.Scan(&as.Authsessionid, &as.CreatedDateTime, &as.LastActiveDateTime, &as.UserUUID, &as.SessionUUID)
	if err != nil {
		return nil, err
	}
	return as, nil
}

func (ast *AuthSessionsTable) SelectByUserUUID(db *sql.DB, userUUID string) (*AuthSession, error) {
	as := &AuthSession{}
	row := db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE useruuid = '%s'", ast.Name(), userUUID))
	err := row.Scan(&as.Authsessionid, &as.CreatedDateTime, &as.LastActiveDateTime, &as.UserUUID, &as.SessionUUID)
	if err != nil {
		return nil, err
	}
	return as, nil
}

func (ast *AuthSessionsTable) Delete(db *sql.DB, whereClause string) error {
	if len(whereClause) > 0 {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", ast.Name(), whereClause))
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Where to delete clause is blank")
}

func (ast *AuthSessionsTable) DeleteBySessionUUID(db *sql.DB, sessionUUID string) error {
	if len(sessionUUID) > 0 {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE sessionuuid = ?", ast.Name()), sessionUUID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Session UUID to delete by is blank")
}

//BuildFields takes the table struct and maps all of the struct fields to their own struct
func (ast *AuthSessionsTable) buildFields() []Field {
	return buildFieldsFromTable(ast)
}

func (ast *AuthSessionsTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(ast, m)
}

func (ast *AuthSessionsTable) buildPreparedInsertStatement(m Model) string {
	return buildPreparedInsertStatementFromTable(ast, m)
}

// ******** End Auth Table ********

// ******** Start SystemInfo Table ********

type SystemInfoTable struct {
	Version string `tbl:NN`
}

func (sit *SystemInfoTable) Init(db *sql.DB) {
	err := sit.Insert(db, &SystemInfo{Version: VERSION})
	if err != nil {
		logging.ErrorAndExit(fmt.Sprintf("Issue creating version string record: %s", err.Error()))
	}
}

func (sit *SystemInfoTable) Name() string { return "systeminfo" }

func (sit *SystemInfoTable) Insert(db *sql.DB, systemInfo *SystemInfo) error {
	insertStatement := sit.buildPreparedInsertStatement(systemInfo)
	_, err := db.Exec(insertStatement, systemInfo.Version)
	if err != nil {
		return err
	}
	return nil
}

func (sit *SystemInfoTable) Update(db *sql.DB, as *SystemInfo) error {
	updateStatement := fmt.Sprintf("UPDATE %s SET version = ?", sit.Name())
	_, err := db.Exec(updateStatement, sit.Version)
	if err != nil {
		return err
	}
	return nil
}

func (sit *SystemInfoTable) buildFields() []Field {
	return buildFieldsFromTable(sit)
}

func (sit *SystemInfoTable) buildInsertStatement(m Model) string {
	return buildInsertStatementFromTable(sit, m)
}

func (sit *SystemInfoTable) buildPreparedInsertStatement(m Model) string {
	return buildPreparedInsertStatementFromTable(sit, m)
}

// ******** End SystemInfo Table ********

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
	UserId          int    `tbl:"AI" json:"userid"`
	CreatedDateTime int64  `json:"createddatetime"`
	UserroleId      int    `json:"userroleid"`
	UUID            string `json:"UUID"`
	Username        string `json:"username"`
	AuthHash        string `json:"authhash"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	Email           string `json:"email"`
}

//Login takes the current username and authhash values of self and tries
//using them to authenticate/login. A successful login will return/generate
//a JWT token for further use in any subsequent API request
func (u *User) Login() bool {
	ut := &UsersTable{}
	user, err := ut.SelectByUsername(Conn, u.Username)

	if err != nil {
		logging.Error(err.Error())
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.AuthHash), []byte(u.AuthHash))
	if err == nil {
		return true
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

type Group struct {
	Groupid         int    `tbl:"AI" json:"groupid"`
	CreatedDateTime int64  `json:"createddatetime"`
	UUID            string `json:"UUID"`
	Title           string `json:"title"`
}

func (g *Group) TableName() string {
	return "groups"
}

func (g *Group) BuildFields() []Field {
	return buildFieldsFromModel(g)
}

func (g *Group) Validate() bool {
	return len(g.UUID) > 0
}

type GroupMembership struct {
	Groupmembershipid int    `tbl:"AI" json:"groupmembershipid"`
	CreatedDateTime   int64  `json:"createddatetime"`
	UUID              string `json:"UUID"`
	GroupUUID         string `json:"GroupUUID"`
	UserUUID          string `json:"UserUUID"`
}

func (gm *GroupMembership) TableName() string {
	return "groupmemberships"
}

func (gm *GroupMembership) BuildFields() []Field {
	return buildFieldsFromModel(gm)
}

func (gm *GroupMembership) Validate() bool {
	return len(gm.UUID) > 0
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
	PageId          int    `tbl:"AI" json:"pageid"`
	CreatedDateTime int64  `json:"createddatetime"`
	UUID            string `json:"UUID"`
	Roleprotected   bool   `json:"roleprotected"`
	AuthorUUID      string `json:"authoruuid"`
	Title           string `json:"title"`
	Route           string `json:"route"`
	Content         string `json:"content"`
}

func (p *Page) TableName() string {
	return "pages"
}

func (p *Page) BuildFields() []Field {
	return buildFieldsFromModel(p)
}

type AuthSession struct {
	Authsessionid      int    `tbl:"AI" json:"authsessionid"`
	CreatedDateTime    int64  `json:"createddatetime"`
	LastActiveDateTime int64  `json:"lastactivedatetime"`
	UserUUID           string `json:"userUUID"`
	SessionUUID        string `json:"sessionUUID"`
}

func (as *AuthSession) TableName() string {
	return "authsessions"
}

func (as *AuthSession) BuildFields() []Field {
	return buildFieldsFromModel(as)
}

func (as *AuthSession) Validate() bool {
	return len(as.UserUUID) > 0 && len(as.SessionUUID) > 0
}

type SystemInfo struct {
	Version string `json:"version"`
}

func (si *SystemInfo) TableName() string {
	return "systeminfo"
}

func (si *SystemInfo) BuildFields() []Field {
	return buildFieldsFromModel(si)
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
			if modelField.Type != "boolean" && modelField.Type != "BIT(1)" && modelField.Type != "BIGINT" {
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

func buildPreparedInsertStatementFromTable(t Table, m Model) string {
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
			insertStatementBuilder.WriteString("?")
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
	stringBulder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (", t.Name()))

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
			if Type == MySQL {
				pkFieldCount++
				if pkFieldCount > 1 {
					logging.ErrorAndExit(fmt.Sprintf("Error creating %s table: More than one PK field found...", t.Name()))
				}
			} else if Type == SQLITE {
				stringBulder.WriteString(" PRIMARY KEY")
			}
			pkField = field
		}
		if field.AutoIncrement {
			if Type == MySQL {
				stringBulder.WriteString(" AUTO_INCREMENT")
			} else if Type == SQLITE {
				stringBulder.WriteString(" AUTOINCREMENT")
			}
		}
		if field.NotNull {
			stringBulder.WriteString(" NOT NULL")
		}
		if field.UniqueIndex {
			if Type == MySQL {
				uniqueIndexFields = append(uniqueIndexFields, field)
				uniqueIndexFieldsCount++
			} else if Type == SQLITE {
				stringBulder.WriteString(" UNIQUE")
			}
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
