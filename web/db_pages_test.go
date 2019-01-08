// Copyright (c) 2019, tacusci ltd
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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	db.Connect(db.SQLITE, "./berrycmstesting.db", "")
	db.Wipe()
	db.Setup()
}

func TestSavedPageGet(t *testing.T) {
	sph := SavedPageHandler{}
	req := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()

	sph.Get(responseRecorder, req)

	resp := responseRecorder.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Test get retrieved a response which is not 404 for page which doesn't exist...")
	}

	pt := db.PagesTable{}

	pt.Insert(db.Conn, &db.Page{
		CreatedDateTime: time.Now().Unix(),
		Roleprotected:   false,
		AuthorUUID:      "",
		Title:           "Test Page",
		Route:           "/testpage",
		//page content is never saved as HTML but instead as QuillJS delta JSON objects
		Content: "[{\"insert\":\"This is a test page!\\n\"}]",
	})

	req = httptest.NewRequest("GET", "/testpage", nil)
	responseRecorder = httptest.NewRecorder()

	sph.Get(responseRecorder, req)

	resp = responseRecorder.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test get retrieved a respone which is not 200, STATUS: %d...", resp.StatusCode)
	}

	if bodyText, err := ioutil.ReadAll(resp.Body); err == nil {
		if "<html><head><link rel=\"stylesheet\" href=\"/css/berry-default.css\"><link rel=\"stylesheet\" href=\"/css/font.css\"></head><body><p>This is a test page!</p></body></html>" != string(bodyText) {
			t.Errorf("Fetched page content does not match expected content")
		}
	}
}
