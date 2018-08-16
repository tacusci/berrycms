package web

import "net/http"

//LoginHandler contains response functions for admin login
type LogoutHandler struct {
	Router *MutableRouter
	route  string
}

func (lh *LogoutHandler) Get(w http.ResponseWriter, r *http.Request) {}

func (lh *LogoutHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (lh *LogoutHandler) Route() string { return lh.route }

func (lh *LogoutHandler) HandlesGet() bool  { return false }
func (lh *LogoutHandler) HandlesPost() bool { return true }
