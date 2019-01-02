package web

import (
	"net/http"
)

type RobotsHandler struct {
	Router *MutableRouter
	route  string
}

func (rh *RobotsHandler) Get(w http.ResponseWriter, r *http.Request) {

}
