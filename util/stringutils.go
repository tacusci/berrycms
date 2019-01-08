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

package util

import (
	"regexp"

	"github.com/tacusci/logging"
	"golang.org/x/crypto/bcrypt"
)

type CompiledRegex struct {
	*regexp.Regexp
}

func (cr *CompiledRegex) GetMatchGroupContent(s string, gi int) string {
	result := cr.FindStringSubmatch(s)
	if len(result) >= gi {
		return result[gi]
	}
	return ""
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		logging.ErrorAndExit(err.Error())
	}
	return string(hash)
}
