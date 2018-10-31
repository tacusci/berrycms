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
	"net/http/httptest"
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

	err := ut.Insert(db.Conn, userToCreate)

	if err != nil {
		t.Errorf("Error occurred inserting test root user %v", err)
	}

	rootUser, err := ut.SelectByUsername(db.Conn, "rootuser")

	if err != nil {
		t.Errorf("Error occurred inserting test root user %v", err)
	}

	audh := AdminUsersDeleteHandler{}
	req := httptest.NewRequest("POST", "/admin/users/delete", nil)
	responseRecorder := httptest.NewRecorder()

	//TODO: create and set PostForm to request here

	audh.Post(responseRecorder, req)
}
