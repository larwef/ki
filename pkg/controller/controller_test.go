package controller

import (
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/listing"
	"github.com/larwef/ki/pkg/repository/memory"
	"github.com/larwef/ki/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseHttpHandler_InvalidPath(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/invalidpath", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), "Not Found\n")
}
