package web

import (
	"net/http"
)

//AdminUserGroupsEditAddHandler contains response functions for pages admin page
type AdminUserGroupsEditAddHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Get(w http.ResponseWriter, r *http.Request) {}

//Post handles post requests to URI
func (augeah *AdminUserGroupsEditAddHandler) Post(w http.ResponseWriter, r *http.Request) {

}

//Route get URI route for handler
func (augeah *AdminUserGroupsEditAddHandler) Route() string { return augeah.route }

//HandlesGet retrieve whether this handler handles get requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesGet() bool { return false }

//HandlesPost retrieve whether this handler handles post requests
func (augeah *AdminUserGroupsEditAddHandler) HandlesPost() bool { return true }
