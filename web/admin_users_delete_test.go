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

package web

import (
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/util"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func init() {
	//this is here to make sure the data contains only pages we're creating now
	db.Connect(db.SQLITE, "./berrycmstesting.db", "")
	db.Wipe()
	db.Setup()

}

func TestDeleteUsersPost(t *testing.T) {
	ut := db.UsersTable{}
	userToCreate := db.User{
		Username:        "rootuser",
		CreatedDateTime: time.Now().Unix(),
		Email:           "root@local.com",
		UserroleId:      int(db.ROOT_USER),
		FirstName:       "Root",
		LastName:        "User",
		AuthHash:        util.HashAndSalt([]byte("testingrootpass")),
	}

	adminUserNotRoot := db.User{
		Username:        "adminuser",
		CreatedDateTime: time.Now().Unix(),
		Email:           "admin@local.com",
		UserroleId:      int(db.MOD_USER),
		FirstName:       "Admin",
		LastName:        "User",
		AuthHash:        util.HashAndSalt([]byte("testuser")),
	}

	err := ut.Insert(db.Conn, userToCreate)

	if err != nil {
		t.Errorf("Error occurred inserting test root user %v", err)
	}

	err = ut.Insert(db.Conn, adminUserNotRoot)

	if err != nil {
		t.Errorf("Error occurred inserting test admin user %v", err)
	}

	rootUser, err := ut.SelectByUsername(db.Conn, "rootuser")

	if err != nil {
		t.Errorf("Error occurred trying to load root user from db %v", err)
	}

	if rootUser.UserroleId != int(db.ROOT_USER) {
		t.Error("Root user loaded from db is not set as a root user role type")
	}

	adminUserNotRoot, err = ut.SelectByUsername(db.Conn, "adminuser")

	if err != nil {
		t.Errorf("Error occurred trying to load admin user from db %v", err)
	}

	//above we've created the root user, next steps are to try and delete them using post form

	audh := AdminUsersDeleteHandler{}

	responseRecorder := httptest.NewRecorder()

	formValues := url.Values{}
	formValues["0"] = []string{rootUser.UUID}
	formValues["1"] = []string{adminUserNotRoot.UUID}

	req := httptest.NewRequest("POST", "/admin/users/delete", nil)
	req.PostForm = formValues

	audh.Post(responseRecorder, req)

	resp := responseRecorder.Result()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Test new user post didn't redirect request, STATUS CODE: %d", resp.StatusCode)
	}

	//location header will have been set on http server redirect
	if len(resp.Header["Location"]) > 0 && resp.Header["Location"][0] != "/admin/users" {
		t.Errorf("Test post new root user didn't set header to redirect to correct location")
	}

	rootUser, err = ut.SelectByUsername(db.Conn, "rootuser")

	if err != nil {
		t.Errorf("Error occurred trying to read root user from db %v", err)
	}

	if rootUser.Username != "rootuser" {
		t.Errorf("Root user no longer exists but shouldn't have been deleted")
	}

	adminUserNotRoot, err = ut.SelectByUsername(db.Conn, "adminuser")

	if err != nil {
		t.Errorf("Error occurred trying to read admin user from db %v", err)
	}

	if adminUserNotRoot.Username == "adminuser" {
		t.Error("Admin user still exists but should have been deleted")
	}
}
