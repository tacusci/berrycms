package plugins

import (
	"fmt"

	"github.com/robertkrimen/otto"
	"github.com/tacusci/logging"
)

func PluginInfoLog(call otto.FunctionCall) otto.Value {
	// unsafe, not confirming argument length
	if uuid, err := call.Otto.Get("UUID"); err == nil {
		if uuid.IsString() {
			logging.Info(fmt.Sprintf("PLUGIN {%s}: %s", uuid.String(), call.Argument(0).String()))
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
			logging.Debug(fmt.Sprintf("%s", call.Argument(0).String()))
		}
	} else {
		logging.Error(err.Error())
	}
	return otto.Value{}
}

func PluginErrorLog(call otto.FunctionCall) otto.Value {
	// unsafe, not confirming argument length
	logging.ErrorNoColor(fmt.Sprintf("%s", call.Argument(0).String()))
	return otto.Value{}
}
