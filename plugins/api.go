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
	"fmt"

	"github.com/robertkrimen/otto"
	"github.com/tacusci/berrycms/robots"
	"github.com/tacusci/logging"
)

// ******** LOGGING FUNCS ********

func PluginInfoLog(call otto.FunctionCall) otto.Value {
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

func PluginDebugLog(call otto.FunctionCall) otto.Value {
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

func PluginErrorLog(call otto.FunctionCall) otto.Value {
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

// ******** ROBOTS UTILS FUNCS ********

func AddRobotsEntry(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		apiError(&call, "'AddRobotsEntry' function takes only one argument")
		return otto.Value{}
	}
	var valPassed otto.Value = call.Argument(0)
	if !valPassed.IsString() {
		apiError(&call, "'AddRobotsEntry' function expected string")
		return otto.Value{}
	}
	val := []byte(valPassed.String())
	robots.Add(&val)
	return otto.Value{}
}

// ******** END ROBOTS UTILS FUNCS ********

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
