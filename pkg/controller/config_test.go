package controller

import (
	"bytes"
	"encoding/json"
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/listing"
	"github.com/larwef/ki/pkg/repository/memory"
	"github.com/larwef/ki/testutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testDataFolder = "testdata/"

func TestConfigHandler_InvalidConfigPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/config/someGroup/someId/somethingElse", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusBadRequest)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), "Invalid Path\n")
}

func TestConfigHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("INVALID", "/config/someGroup/someId", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusMethodNotAllowed)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), "Method Not Allowed\n")
}

func TestConfigHandler_PutGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/config/someGroup/", bytes.NewBufferString("{}"))
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)
	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var grpResponse listing.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	testutil.AssertNotError(t, err)

	testutil.AssertEqual(t, grpResponse.ID, "someGroup")
	testutil.AssertEqual(t, len(grpResponse.Configs), 0)
}

func TestConfigHandler_GetGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	handler.ServeHTTP(res, req)
	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var grpResponse listing.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	testutil.AssertNotError(t, err)

	testutil.AssertEqual(t, grpResponse.ID, "someGroup")
	testutil.AssertEqual(t, len(grpResponse.Configs), 3)
}

func TestConfigHandler_GetGroup_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	handler.ServeHTTP(res, req)
	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_PutConfig(t *testing.T) {
	file, err := os.OpenFile(testDataFolder+"configExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someOtherGroup",
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var configResponse listing.Config
	err = json.NewDecoder(res.Body).Decode(&configResponse)
	testutil.AssertNotError(t, err)

	file, err = os.OpenFile(testDataFolder+"storedConfigExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)
	var configExpected listing.Config
	err = json.NewDecoder(file).Decode(&configExpected)
	testutil.AssertNotError(t, err)

	testutil.AssertEqual(t, configResponse.ID, configExpected.ID)
	testutil.AssertEqual(t, configResponse.Name, configExpected.Name)
	testutil.AssertEqual(t, configResponse.Group, configExpected.Group)
}

func TestConfigHandler_PutConfig_GroupNotFound(t *testing.T) {
	file, err := os.OpenFile(testDataFolder+"configExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_GetConfig(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someId", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	var c adding.Config
	testutil.UnmarshalJSONFromFile(t, testDataFolder+"configExample.json", &c)
	repository.StoreConfig(c)

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testutil.AssertJSONEqual(t, res.Body.String(), testutil.GetTestFileAsString(t, testDataFolder+"configExample.json"))
}

func TestConfigHandler_GetConfig_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/someId", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	repository.StoreConfig(adding.Config{
		ID: "someId",
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_GetConfig_ConfigNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someOtherId", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	repository.StoreConfig(adding.Config{
		ID: "someId",
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), listing.ErrConfigNotFound.Error()+"\n")
}

func TestConfigHandler_pathValidation(t *testing.T) {
	testutil.AssertEqual(t, invalidConfigPath("/grp"), false)
	testutil.AssertEqual(t, invalidConfigPath("/grp/"), false)
	testutil.AssertEqual(t, invalidConfigPath("/grp/config"), false)
	testutil.AssertEqual(t, invalidConfigPath("/grp/config/"), false)

	testutil.AssertEqual(t, invalidConfigPath(""), true)
	testutil.AssertEqual(t, invalidConfigPath("/"), true)
	testutil.AssertEqual(t, invalidConfigPath("/grp/config/stuff"), true)
	testutil.AssertEqual(t, invalidConfigPath("/grp/config/stuff/"), true)
	testutil.AssertEqual(t, invalidConfigPath("/grp/config/stuff/morestuff"), true)
}
