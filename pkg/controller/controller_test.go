package controller

import (
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/listing"
	"github.com/larwef/ki/pkg/repository/memory"
	"github.com/larwef/ki/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseHttpHandler_InvalidPath(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/invalidpath", nil)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Not Found\n")
}
