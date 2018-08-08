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
