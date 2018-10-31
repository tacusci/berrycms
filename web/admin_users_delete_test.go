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
	"net/http/httptest"
	"testing"
)

func TestDeleteUsersPost(t *testing.T) {
	audh := AdminUsersDeleteHandler{}
	req := httptest.NewRequest("POST", "/admin/users/delete", nil)
	responseRecorder := httptest.NewRecorder()

	//TODO: create and set PostForm to request here

	audh.Post(responseRecorder, req)
}
