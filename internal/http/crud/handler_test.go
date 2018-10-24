package crud

import (
	"bytes"
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/http/auth"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/internal/repository/memory"
	"github.com/larwef/ki/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testDataFolder = "../../../test/testdata/"

func setup(t *testing.T) (*Handler, *memory.Repository) {
	adminPw, err := auth.HashPassword("adminPassword123")
	test.AssertNotError(t, err)

	admin := auth.User{
		Username:     "admin",
		PasswordHash: adminPw,
		Role:         auth.ADMIN,
	}

	clientPw, err := auth.HashPassword("clientPassword321")
	test.AssertNotError(t, err)
	client := auth.User{
		Username:     "client",
		PasswordHash: clientPw,
		Role:         auth.CLIENT,
	}

	pool := auth.NewUserPool()
	err = pool.RegisterUser(admin)
	test.AssertNotError(t, err)
	err = pool.RegisterUser(client)
	test.AssertNotError(t, err)

	repository := memory.NewRepository()
	return &Handler{
		userpool: pool,
		adding:   adding.NewService(repository),
		listing:  listing.NewService(repository),
	}, repository
}

func TestHandler_HealthHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, healthPath, nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), contentType)
	test.AssertJSONEqual(t, res.Body.String(), test.GetTestFileAsString(t, testDataFolder+"healthResponse.json"))
}

func TestHandler_InvalidPath(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/invalidpath", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Not Found\n")
}

func TestHandler_InvalidConfigPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/config/someGroup/someId/somethingElse", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusBadRequest)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Invalid Path\n")
}

func TestHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("INVALID", "/config/someGroup/someId", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusMethodNotAllowed)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Method Not Allowed\n")
}

func TestHandler_PutGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/config/someGroup/", bytes.NewBufferString("{}"))
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)
	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), contentType)

	var grpResponse listing.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	test.AssertNotError(t, err)

	test.AssertEqual(t, grpResponse.ID, "someGroup")
	test.AssertEqual(t, len(grpResponse.Configs), 0)
}

func TestHandler_PutGroup_Duplicate(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/config/someGroup", bytes.NewBufferString("{}"))
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	handler.ServeHTTP(res, req)
	test.AssertEqual(t, res.Code, http.StatusConflict)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), adding.ErrGroupConflict.Error()+"\n")
}

func TestHandler_GetGroup(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	handler.ServeHTTP(res, req)
	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), contentType)

	var grpResponse listing.Group
	err = json.NewDecoder(res.Body).Decode(&grpResponse)
	test.AssertNotError(t, err)

	test.AssertEqual(t, grpResponse.ID, "someGroup")
	test.AssertEqual(t, len(grpResponse.Configs), 3)
}

func TestHandler_GetGroup_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	handler.ServeHTTP(res, req)
	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestHandler_PutConfig(t *testing.T) {
	file, err := os.OpenFile(testDataFolder+"configExample.json", os.O_RDONLY, 0644)
	test.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

	repository.StoreGroup(adding.Group{
		ID: "someOtherGroup",
	})

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), contentType)

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

func TestHandler_PutConfig_GroupNotFound(t *testing.T) {
	file, err := os.OpenFile(testDataFolder+"configExample.json", os.O_RDONLY, 0644)
	test.AssertNotError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/config/someOtherGroup/someOtherId", file)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusNotFound)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), listing.ErrGroupNotFound.Error()+"\n")
}

func TestHandler_GetConfig(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someId", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	var c adding.Config
	test.UnmarshalJSONFromFile(t, testDataFolder+"configExample.json", &c)
	repository.StoreConfig(c)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusOK)
	test.AssertEqual(t, res.Header().Get("Content-Type"), contentType)
	test.AssertJSONEqual(t, res.Body.String(), test.GetTestFileAsString(t, testDataFolder+"configExample.json"))
}

func TestHandler_GetConfig_GroupNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someOtherGroup/someId", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

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

func TestHandler_GetConfig_ConfigNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/someOtherId", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "adminPassword123")

	res := httptest.NewRecorder()
	handler, repository := setup(t)

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

func TestHandler_AuthenticateNoAuthHeader(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/SomeId", nil)
	test.AssertNotError(t, err)

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusUnauthorized)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Unauthorized\n")
}

func TestHandler_AuthenticateNonExistingUser(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/SomeId", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("nonExistingUser", "adminPassword123")

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusUnauthorized)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Unauthorized\n")
}

func TestHandler_AuthenticateWrongPassword(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/config/someGroup/SomeId", nil)
	test.AssertNotError(t, err)
	req.SetBasicAuth("admin", "wrongPassword")

	res := httptest.NewRecorder()
	handler, _ := setup(t)

	handler.ServeHTTP(res, req)

	test.AssertEqual(t, res.Code, http.StatusUnauthorized)
	test.AssertEqual(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, res.Body.String(), "Unauthorized\n")
}

//func TestHandler_AuthenticateInsufficientRole(t *testing.T) {
//	t.Fatal("Test not implemented")
//}
