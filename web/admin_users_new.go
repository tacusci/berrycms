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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/util"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/logging"
)

const usernameRegex = "^[A-Za-z0-9]+(?:[ _-][A-Za-z0-9]+)*$"
const emailRegex = "^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"

//AdminUsersNewHandler handler to contain pointer to core router and the URI string
type AdminUsersNewHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (aunh *AdminUsersNewHandler) Get(w http.ResponseWriter, r *http.Request) {
	var postRequestForNewRootUser = strings.Compare(r.RequestURI, "/admin/users/root/new") == 0

	pctx := plush.NewContext()
	if postRequestForNewRootUser {
		pctx.Set("title", "New Root User")
		pctx.Set("navBarEnabled", false)
		pctx.Set("newuserformaction", "/admin/users/root/new")
		pctx.Set("adminhiddenpassword", "")
		pctx.Set("createuserlabel", "Create Root User")
	} else {
		pctx.Set("title", "New User")
		pctx.Set("navBarEnabled", true)
		pctx.Set("newuserformaction", "/admin/users/new")
		pctx.Set("adminhiddenpassword", "")
		if aunh.Router.AdminHidden {
			pctx.Set("adminhiddenpassword", fmt.Sprintf("/%s", aunh.Router.AdminHiddenPassword))
		}
		pctx.Set("createuserlabel", "Create New User")
	}
	RenderDefault(w, "admin.users.new.html", pctx)
}

//Post handles post requests to URI
func (aunh *AdminUsersNewHandler) Post(w http.ResponseWriter, r *http.Request) {

	var postRequestForNewRootUser = strings.Compare(r.RequestURI, "/admin/users/root/new") == 0

	// the new user form can normally only be accessed by a logged in user, so just redirect to users man page
	if postRequestForNewRootUser {
		// if the new root user has just been created there won't be a login session nor other users
		var redirectURI = "/login"

		if aunh.Router.AdminHidden {
			redirectURI = fmt.Sprintf("/%s", aunh.Router.AdminHiddenPassword) + redirectURI
		}

		defer http.Redirect(w, r, redirectURI, http.StatusFound)
	} else {
		var redirectURI = "/admin/users"

		if aunh.Router.AdminHidden {
			redirectURI = fmt.Sprintf("/%s", aunh.Router.AdminHiddenPassword) + redirectURI
		}

		defer http.Redirect(w, r, redirectURI, http.StatusFound)
	}

	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		return
	}

	if validated, err := validatePostForm(r); err != nil || validated == false {
		logging.Error(err.Error())
		return
	}

	authHash := r.PostFormValue("authhash")
	repeatedAuthHash := r.PostFormValue("repeatedauthhash")

	if strings.Compare(authHash, repeatedAuthHash) == 0 {
		var userRoleID int
		if postRequestForNewRootUser {
			userRoleID = int(db.ROOT_USER)
		} else {
			userRoleID = int(db.REG_USER)
		}
		ut := db.UsersTable{}
		userToCreate := &db.User{
			Username:        r.PostFormValue("username"),
			CreatedDateTime: time.Now().Unix(),
			Email:           r.PostFormValue("email"),
			UserroleId:      userRoleID,
			FirstName:       r.PostFormValue("firstname"),
			LastName:        r.PostFormValue("lastname"),
			AuthHash:        util.HashAndSalt([]byte(authHash)),
		}

		if err := ut.Insert(db.Conn, userToCreate); err != nil {
			Error(w, err)
		}

		if postRequestForNewRootUser {
			logging.Debug("Root user POST form recieved")
			rootUser, err := ut.SelectRootUser(db.Conn)

			if err != nil {
				logging.Error("Unable to retrieve root user")
				Error(w, err)
			}

			gmt := db.GroupMembershipTable{}
			gmt.AddUserToGroup(db.Conn, rootUser, "Admins")
		}
	} else {
		//need to add setting error message on screen
	}
}

//Route get URI route for handler
func (aunh *AdminUsersNewHandler) Route() string { return aunh.route }

//HandlesGet retrieve whether this handler handles get requests
func (aunh *AdminUsersNewHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (aunh *AdminUsersNewHandler) HandlesPost() bool { return true }

// validate makes sure that the passwords match and that the username and email are in correct format
func validatePostForm(r *http.Request) (bool, error) {
	authHash := r.PostFormValue("authhash")
	repeatedAuthHash := r.PostFormValue("repeatedauthhash")

	if strings.Compare(authHash, repeatedAuthHash) != 0 {
		return false, errors.New("Password and repeated passwords don't match")
	}

	firstname := r.PostFormValue("firstname")
	lastname := r.PostFormValue("lastname")
	email := r.PostFormValue("email")
	username := r.PostFormValue("username")

	if len(firstname) <= 0 || len(lastname) <= 0 || len(email) <= 0 || len(username) <= 0 {
		return false, errors.New("One of required fields is blank")
	}

	if match, err := regexp.MatchString(usernameRegex, username); err != nil || match == false {
		return false, errors.New("Username does not match pattern regex")
	}

	if match, err := regexp.MatchString(emailRegex, email); err != nil || match == false {
		return false, errors.New("Email does not match pattern regex")
	}

	return true, nil
}
