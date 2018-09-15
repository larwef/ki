package controller

import (
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/repository"
	"github.com/larwef/ki/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseHttpHandler_InvalidPath(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/invalidpath", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{StoredConfig: config.Config{}})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), "Not Found\n")
}
