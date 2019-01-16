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
	"net/http"

	"github.com/tacusci/berrycms/robots"
)

type RobotsHandler struct {
	Router *MutableRouter
	route  string
}

func (rh *RobotsHandler) Get(w http.ResponseWriter, r *http.Request) {
	//if the robots .txt file cache hasn't been created then technically there is no robots page
	if !robots.CacheExists() {
		fourOhFour(w, r)
		return
	}

	w.Write(robots.CacheBytes())
}

func (rh *RobotsHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (rh *RobotsHandler) Route() string { return rh.route }

//HandlesGet retrieve whether this handler handles get requests
func (rh *RobotsHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (rh *RobotsHandler) HandlesPost() bool { return false }
