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

package robots

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tacusci/berrycms/db"
)

var cache *bytes.Buffer

func Add(val *[]byte) error {
	if cache == nil {
		return errors.New("Robots cache unmutable... User has likely disabled robots.txt")
	}
	//add newline to uri to add so caller doesn't have to
	*val = append(*val, []byte("\n")...)
	_, err := cache.Write(*val)
	if err != nil {
		return err
	}
	return nil
}

func Del(val *[]byte) error {
	if cache == nil {
		return errors.New("Robots cache unmutable... User has likely disabled robots.txt")
	}

	//*OPTIMISATION* immediately return if there's nothing to delete from
	if cache.Len() == 0 {
		return nil
	}

	//add newline to uri to add so caller doesn't have to
	*val = append(*val, []byte("\n")...)
	existingVal := cache.Bytes()
	cache.Reset()
	cache.Write(bytes.Replace(existingVal, *val, []byte{}, -1))
	return nil
}

func Generate() error {
	Reset()
	_, err := cache.WriteString("User-agent: *\n")
	if err != nil {
		return err
	}
	//block indexing admin pages
	//NOTE: if the admin pages URI has been hidden, we're deliberately omitting this from robots.txt
	_, err = cache.WriteString("Disallow: /admin\n")
	if err != nil {
		return err
	}

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "route", "roleprotected = '1'")

	if err != nil {
		return err
	}

	var pageRouteToDisallow string

	for rows.Next() {
		err := rows.Scan(&pageRouteToDisallow)
		if err != nil {
			return err
		}

		_, err = cache.WriteString(fmt.Sprintf("Disallow: %s\n", pageRouteToDisallow))
		if err != nil {
			return err
		}
	}

	return nil
}

func CacheExists() bool {
	return cache != nil
}

func CacheBytes() []byte {
	return cache.Bytes()
}

func Reset() {
	//we don't want to allocate memory each reset
	if cache == nil {
		cache = &bytes.Buffer{}
	}
	cache.Reset()
}
