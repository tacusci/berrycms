package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSavedPageGet(t *testing.T) {
	sph := SavedPageHandler{}
	req := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()

	sph.Get(responseRecorder, req)

	resp := responseRecorder.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Test get retrieved a response which is not 404...")
	}
}
