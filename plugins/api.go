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

package plugins

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/robots"
	"github.com/tacusci/logging"
)

// ******** LOGGING FUNCS ********

type logapi struct{}

func (l *logapi) Info(call otto.FunctionCall) otto.Value {
	// unsafe, not confirming argument length
	if uuid, err := call.Otto.Get("UUID"); err == nil {
		if uuid.IsString() {
			logging.Info(fmt.Sprintf("PLUGIN {%s} -> %s", uuid.String(), call.Argument(0).String()))
		}
	} else {
		logging.Error(err.Error())
	}
	return otto.Value{}
}

func (l *logapi) Debug(call otto.FunctionCall) otto.Value {
	// unsafe, not confirming argument length
	if uuid, err := call.Otto.Get("UUID"); err == nil {
		if uuid.IsString() {
			logging.Debug(fmt.Sprintf("PLUGIN {%s} -> %s", uuid.String(), call.Argument(0).String()))
		}
	} else {
		logging.Error(err.Error())
	}
	return otto.Value{}
}

func (l *logapi) Error(call otto.FunctionCall) otto.Value {
	// unsafe, not confirming argument length
	if uuid, err := call.Otto.Get("UUID"); err == nil {
		if uuid.IsString() {
			logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", uuid.String(), call.Argument(0).String()))
		}
	} else {
		logging.Error(err.Error())
	}
	return otto.Value{}
}

// ******** END LOGGING FUNCS ********

// ******** CMS DATABASE FUNCS ********

type cmsdatabaseapi struct {
	Conn       *sql.DB
	PagesTable *db.PagesTable
	UsersTable *db.UsersTable
}

// ******** END CMS DATABASE FUNCS ********

// ******** DATABASE FUNCS ********

type databaseapi struct {
	Conn *sql.DB
}

func (d *databaseapi) Connect(user string, password string, addr string) {}

// ******** END DATABASE FUNCS ********

// ******** ROBOTS UTILS FUNCS ********

type robotsapi struct{}

func (r *robotsapi) Add(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		apiError(&call, "too many arguments to call 'robots.Add', want (string)")
		return otto.Value{}
	}
	var valPassed otto.Value = call.Argument(0)
	if !valPassed.IsString() {
		apiError(&call, "'robots.Add' function expected string")
		return otto.Value{}
	}
	val := []byte(valPassed.String())
	err := robots.Add(&val)

	if err != nil {
		apiError(&call, err.Error())
	}
	return otto.Value{}
}

func (r *robotsapi) Del(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		apiError(&call, "too many arguments to call 'robots.Del', want (string)")
		return otto.Value{}
	}
	var valPassed otto.Value = call.Argument(0)
	if !valPassed.IsString() {
		apiError(&call, "'robots.Del' function expected string")
		return otto.Value{}
	}
	val := []byte(valPassed.String())

	err := robots.Del(&val)
	if err != nil {
		apiError(&call, err.Error())
	}

	return otto.Value{}
}

// ******** END ROBOTS UTILS FUNCS ********

// ******** FILE FUNCS ********

type filesapi struct{}

func (f *filesapi) Read(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		return apiError(&call, "too many arguments to call 'files.Read', want (string)")
	}
	var valPassed otto.Value = call.Argument(0)
	if !valPassed.IsString() {
		return apiError(&call, "'files.Read' function expected string")
	}

	//get the current running/working directory abs path
	rootDir, err := os.Getwd()
	if err != nil {
		return apiError(&call, err.Error())
	}

	//append the running/working directory abs path to the passed path
	//we want to remove the leading path seperator to guarentee there's only the one
	//also want to remove '../' chars so reverse directory traversal is impossible
	absFilePath := fmt.Sprintf("%s%s%s", rootDir, string(os.PathSeparator), strings.Replace(strings.TrimPrefix(valPassed.String(), string(os.PathSeparator)), "../", "", -1))

	//get passed file info
	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		return apiError(&call, err.Error())
	}

	//if the file passed is a directory return list of files in directory
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(absFilePath)
		if err != nil {
			return apiError(&call, err.Error())
		}
		//convert file slice to otto value
		val, err := call.Otto.ToValue(files)
		if err != nil {
			return apiError(&call, err.Error())
		}
		return val
	}

	//read the byte data from the file passed and return it as a string
	data, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		return apiError(&call, err.Error())
	}
	//convert file data string to otto value
	val, err := call.Otto.ToValue(string(data))

	if err != nil {
		return apiError(&call, err.Error())
	}
	return val
}

// ******** END FILE FUNCS ********

// ******** MISC FUNCS ********

func apiError(call *otto.FunctionCall, outputMessage string) otto.Value {
	// unsafe, not confirming argument length
	if uuid, err := call.Otto.Get("UUID"); err == nil {
		if uuid.IsString() {
			logging.Error(fmt.Sprintf("PLUGIN {%s} -> %s", uuid.String(), outputMessage))
		}
	} else {
		logging.Error(err.Error())
	}
	return otto.Value{}
}

// ******** END MISC FUNCS ********
