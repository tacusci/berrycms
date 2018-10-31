package web

import (
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	audh := AdminUsersDeleteHandler{}
	req := httptest.NewRequest("POST", "/admin/users/delete", nil)
	responseRecorder := httptest.NewRecorder()

	//TODO: create and set PostForm to request here

	audh.Post(responseRecorder, req)
}
