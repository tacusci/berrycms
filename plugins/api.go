package plugins

import (
	"fmt"

	"github.com/robertkrimen/otto"
	"github.com/tacusci/logging"
)

func PluginInfoLog(call otto.FunctionCall) otto.Value {
	logging.InfoNoColor(fmt.Sprintf("%s", call.Argument(0).String()))
	return otto.Value{}
}
