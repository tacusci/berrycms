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

package util

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/tacusci/logging"
)

type RecursiveDirWatch struct {
	Change chan bool
	Stop   chan bool
}

func (w *RecursiveDirWatch) WatchDir(sd string) {
	startChecksum := ""
	startTime := time.Now()
	w.genChecksum(sd, &startChecksum)
	for {
		select {
		case <-w.Stop:
			logging.Debug("Stopping directory structure monitoring...")
			return
		default:
			if time.Since(startTime).Seconds() > 5 {
				currentChecksum := ""
				w.genChecksum(sd, &currentChecksum)
				if startChecksum != currentChecksum {
					startChecksum = currentChecksum
					w.Change <- true
				} else {
					w.Change <- false
				}
				startTime = time.Now()
			}
		}
	}
}

func (w *RecursiveDirWatch) genChecksum(sd string, checksum *string) {
	fileList := []string{}
	w.getFileList(sd, &fileList)

	bytesForHash := []byte{}
	for _, fn := range fileList {
		for _, b := range []byte(fn) {
			bytesForHash = append(bytesForHash, b)
		}
	}
	h := md5.New()
	*checksum = fmt.Sprintf("%x", h.Sum(bytesForHash))
}

func (w *RecursiveDirWatch) getFileList(sd string, fl *[]string) {
	fs, err := ioutil.ReadDir(sd)
	if err != nil {
		logging.Error("Unable to find folder...")
		return
	}
	for _, f := range fs {
		if f.IsDir() {
			*fl = append(*fl, f.Name())
			w.getFileList(sd+string(os.PathSeparator)+f.Name(), fl)
		}
	}
}
