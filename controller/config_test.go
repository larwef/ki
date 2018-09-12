package controller

import (
	"encoding/json"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/testutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestConfigHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("INVALID", "/config/someGroup/someId", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(config.NewRepositoryMock(config.Config{}))

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusMethodNotAllowed)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), "Method Not Allowed\n")
}

func TestConfigHandler_HandleGet(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someId", nil)
	if err != nil {
		t.Fatal(err)
	}

	var c config.Config
	testutil.UnmarshalJSONFromFile(t, "../testdata/configExample.json", &c)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(config.NewRepositoryMock(c))

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testutil.AssertJSONEqual(t, res.Body.String(), testutil.GetTestFileAsString(t, "../testdata/configExample.json"))
}

func TestConfigHandler_HandlePut(t *testing.T) {
	file, err := os.OpenFile("../testdata/configExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(config.NewRepositoryMock(config.Config{}))

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var configRespone config.Config
	err = json.NewDecoder(res.Body).Decode(&configRespone)
	testutil.AssertNotError(t, err)

	file, err = os.OpenFile("../testdata/storedConfigExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)
	var configExpected config.Config
	err = json.NewDecoder(file).Decode(&configExpected)
	testutil.AssertNotError(t, err)

	testutil.AssertEqual(t, configRespone.ID, configExpected.ID)
	testutil.AssertEqual(t, configRespone.Name, configExpected.Name)
	testutil.AssertEqual(t, configRespone.Group, configExpected.Group)
}
