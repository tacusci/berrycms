package web

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tacusci/logging"
)

type IndexHandler struct {
	Router *MutableRouter
	route  string
}

func (ih *IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "index.html")
	if err != nil {
		logging.Error("Unable to find resources folder...")
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write(content)
}

func (ih *IndexHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (ih *IndexHandler) Route() string { return ih.route }

func (ih *IndexHandler) HandlesGet() bool  { return true }
func (ih *IndexHandler) HandlesPost() bool { return false }
