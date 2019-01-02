package web

import (
	"net/http"

	"github.com/tacusci/berrycms/robots"
)

type RobotsHandler struct {
	Router *MutableRouter
	route  string
}

func (rh *RobotsHandler) Get(w http.ResponseWriter, r *http.Request) {
	robotsCacheContent, err := robots.RobotsCache.Get([]byte("robots"))
	if err != nil {
		Error(w, err)
		return
	}

	w.Write(robotsCacheContent)
}

func (rh *RobotsHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (rh *RobotsHandler) Route() string { return rh.route }

//HandlesGet retrieve whether this handler handles get requests
func (rh *RobotsHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (rh *RobotsHandler) HandlesPost() bool { return false }
