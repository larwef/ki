package controller

import (
	"bytes"
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/internal/repository/memory"
	"github.com/larwef/ki/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testDataFolder = "../../test/testdata/"

func TestConfigHandler_InvalidConfigPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/config/someGroup/someId/somethingElse", nil)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusBadRequest)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Invalid Path\n")
}

func TestConfigHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("INVALID", "/config/someGroup/someId", nil)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusMethodNotAllowed)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Method Not Allowed\n")
}

func TestConfigHandler_PutGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/config/someGroup/", bytes.NewBufferString("{}"))
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	handler.ServeHTTP(res, req)
	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var grpResponse listing.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	test.AssertNotError(t, err)

	test.AssertEqual(t, grpResponse.ID, "someGroup")
	test.AssertEqual(t, len(grpResponse.Configs), 0)
}

func TestConfigHandler_GetGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/", nil)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	handler.ServeHTTP(res, req)
	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var grpResponse listing.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	test.AssertNotError(t, err)

	test.AssertEqual(t, grpResponse.ID, "someGroup")
	test.AssertEqual(t, len(grpResponse.Configs), 3)
}

func TestConfigHandler_GetGroup_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/", nil)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	handler.ServeHTTP(res, req)
	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_PutConfig(t *testing.T) {
	file, err := os.OpenFile(testDataFolder+"configExample.json", os.O_RDONLY, 0644)
	test.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someOtherGroup",
	})

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var configResponse listing.Config
	err = json.NewDecoder(res.Body).Decode(&configResponse)
	test.AssertNotError(t, err)

	file, err = os.OpenFile(testDataFolder+"storedConfigExample.json", os.O_RDONLY, 0644)
	test.AssertNotError(t, err)
	var configExpected listing.Config
	err = json.NewDecoder(file).Decode(&configExpected)
	test.AssertNotError(t, err)

	test.AssertEqual(t, configResponse.ID, configExpected.ID)
	test.AssertEqual(t, configResponse.Name, configExpected.Name)
	test.AssertEqual(t, configResponse.Group, configExpected.Group)
}

func TestConfigHandler_PutConfig_GroupNotFound(t *testing.T) {
	file, err := os.OpenFile(testDataFolder+"configExample.json", os.O_RDONLY, 0644)
	test.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_GetConfig(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someId", nil)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	repository := memory.NewRepository()
	handler := NewBaseHTTPHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	var c adding.Config
	test.UnmarshalJSONFromFile(t, testDataFolder+"configExample.json", &c)
	repository.StoreConfig(c)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")
	test.AssertJSONEqual(t, res.Body.String(), test.GetTestFileAsString(t, testDataFolder+"configExample.json"))
}

func TestConfigHandler_GetConfig_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/someId", nil)
	test.AssertNotError(t, err)

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

	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_GetConfig_ConfigNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someOtherId", nil)
	test.AssertNotError(t, err)

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

	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), listing.ErrConfigNotFound.Error()+"\n")
}

func TestConfigHandler_pathValidation(t *testing.T) {
	test.AssertEqual(t, invalidConfigPath("/grp"), false)
	test.AssertEqual(t, invalidConfigPath("/grp/"), false)
	test.AssertEqual(t, invalidConfigPath("/grp/config"), false)
	test.AssertEqual(t, invalidConfigPath("/grp/config/"), false)

	test.AssertEqual(t, invalidConfigPath(""), true)
	test.AssertEqual(t, invalidConfigPath("/"), true)
	test.AssertEqual(t, invalidConfigPath("/grp/config/stuff"), true)
	test.AssertEqual(t, invalidConfigPath("/grp/config/stuff/"), true)
	test.AssertEqual(t, invalidConfigPath("/grp/config/stuff/morestuff"), true)
}
