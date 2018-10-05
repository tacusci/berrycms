package plugins

import (
	"fmt"

	"github.com/robertkrimen/otto"
	"github.com/tacusci/logging"
)

func Log(call otto.FunctionCall) otto.Value {
	logging.Info(fmt.Sprintf("%s", call.Argument(0).String()))
	return otto.Value{}
}
