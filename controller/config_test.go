package controller

import (
	"bytes"
	"encoding/json"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/group"
	"github.com/larwef/ki/repository"
	"github.com/larwef/ki/testutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestConfigHandler_InvalidConfigPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/config/someGroup/someId/somethingElse", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusBadRequest)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), "Invalid Path\n")
}

func TestConfigHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("INVALID", "/config/someGroup/someId", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{StoredConfig: config.Config{}})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusMethodNotAllowed)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), "Method Not Allowed\n")
}

func TestConfigHanler_PutGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/config/someGroup/", bytes.NewBufferString("{}"))
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{})

	handler.ServeHTTP(res, req)
	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var grpResponse group.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	testutil.AssertNotError(t, err)

	testutil.AssertEqual(t, grpResponse.ID, "someGroup")
	testutil.AssertEqual(t, len(grpResponse.Configs), 0)
}

func TestConfigHanler_GetGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{
		StoredGroup: group.Group{
			ID:      "someGroup",
			Configs: []string{"config1", "config2", "config3"},
		},
	})

	handler.ServeHTTP(res, req)
	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var grpResponse group.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	testutil.AssertNotError(t, err)

	testutil.AssertEqual(t, grpResponse.ID, "someGroup")
	testutil.AssertEqual(t, len(grpResponse.Configs), 3)
}

func TestConfigHanler_GetGroup_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{
		StoredGroup: group.Group{
			ID:      "someGroup",
			Configs: []string{"config1", "config2", "config3"},
		},
	})

	handler.ServeHTTP(res, req)
	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), repository.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_PutConfig(t *testing.T) {
	file, err := os.OpenFile("../testdata/configExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{
		StoredGroup: group.Group{ID: "someOtherGroup"},
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var configResponse config.Config
	err = json.NewDecoder(res.Body).Decode(&configResponse)
	testutil.AssertNotError(t, err)

	file, err = os.OpenFile("../testdata/storedConfigExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)
	var configExpected config.Config
	err = json.NewDecoder(file).Decode(&configExpected)
	testutil.AssertNotError(t, err)

	testutil.AssertEqual(t, configResponse.ID, configExpected.ID)
	testutil.AssertEqual(t, configResponse.Name, configExpected.Name)
	testutil.AssertEqual(t, configResponse.Group, configExpected.Group)
}

func TestConfigHandler_PutConfig_GroupNotFound(t *testing.T) {
	file, err := os.OpenFile("../testdata/configExample.json", os.O_RDONLY, 0644)
	testutil.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{
		StoredGroup: group.Group{ID: "someGroup"},
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), repository.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_GetConfig(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someId", nil)
	testutil.AssertNotError(t, err)

	var c config.Config
	testutil.UnmarshalJSONFromFile(t, "../testdata/configExample.json", &c)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{
		StoredGroup:  group.Group{ID: "someGroup"},
		StoredConfig: c,
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusOK)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testutil.AssertJSONEqual(t, res.Body.String(), testutil.GetTestFileAsString(t, "../testdata/configExample.json"))
}

func TestConfigHandler_GetConfig_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/someId", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{
		StoredGroup: group.Group{
			ID: "someGroup",
		},
		StoredConfig: config.Config{
			ID: "someId",
		},
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), repository.ErrGroupNotFound.Error()+"\n")
}

func TestConfigHandler_GetConfig_ConfigNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someOtherId", nil)
	testutil.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler := NewBaseHTTPHandler(&repository.Mock{
		StoredGroup: group.Group{
			ID: "someGroup",
		},
		StoredConfig: config.Config{
			ID: "someId",
		},
	})

	handler.ServeHTTP(res, req)

	testutil.AssertEqual(t, res.Code, http.StatusNotFound)
	testutil.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testutil.AssertEqual(t, res.Body.String(), repository.ErrConfigNotFound.Error()+"\n")
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
