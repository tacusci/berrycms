package web

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tacusci/logging"
)

type AdminHandler struct {
	Router *MutableRouter
	route  string
}

func (ah *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.html")
	if err != nil {
		logging.Error("Unable to find resources folder...")
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write(content)
}

func (ah *AdminHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (ah *AdminHandler) Route() string { return ah.route }

func (ah *AdminHandler) HandlesGet() bool  { return true }
func (ah *AdminHandler) HandlesPost() bool { return false }
