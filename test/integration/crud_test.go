// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/test"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var (
	crudAddress        = "http://localhost:8080"
	crudTestDataFolder = "../testdata/"
)

// Integration tests needs a fresh instance running locally to work. The easiest is to run the ./test-docker.sh script.

type TestClient struct {
	client  *http.Client
	baseURL string
}

var testClient TestClient

func TestMain(m *testing.M) {
	testClient = TestClient{
		client:  http.DefaultClient,
		baseURL: crudAddress + "/",
	}

	os.Exit(m.Run())
}

func Test_CRUDAddAndRetrieveGroup(t *testing.T) {
	grpPutReq := adding.Group{ID: "someCrudGroup"}

	var grpPutRes listing.Group
	_, err := testClient.put("config/someCrudGroup", &grpPutReq, &grpPutRes)
	test.AssertNotError(t, err)
	test.AssertEqual(t, grpPutRes.ID, "someCrudGroup")
	test.AssertEqual(t, len(grpPutRes.Configs), 0)

	var grpGetRes listing.Group
	_, err = testClient.get("config/someCrudGroup", &grpGetRes)
	test.AssertNotError(t, err)
	test.AssertEqual(t, grpGetRes.ID, "someCrudGroup")
	test.AssertEqual(t, len(grpGetRes.Configs), 0)
}

func Test_CRUDAddGroup_Duplicate(t *testing.T) {
	grpPutReq := adding.Group{ID: "someCrudGroup"}

	var grpPutRes listing.Group
	_, err := testClient.put("config/someCrudGroupConflict", &grpPutReq, &grpPutRes)
	test.AssertNotError(t, err)

	res, _ := testClient.put("config/someCrudGroupConflict", &grpPutReq, &grpPutRes)
	payload, err := ioutil.ReadAll(res.Body)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res.StatusCode, http.StatusConflict)
	test.AssertEqual(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, string(payload), adding.ErrGroupConflict.Error()+"\n")
}

func Test_CRUDAddAndRetrieveConfig(t *testing.T) {
	grpPutReq := adding.Group{ID: "someCrudGroup"}

	var grpPutRes listing.Group
	_, err := testClient.put("config/someCrudOtherGroup", &grpPutReq, &grpPutRes)

	properties, err := ioutil.ReadFile(grpcTestDataFolder + "properties.json")
	test.AssertNotError(t, err)
	confPutReq := &adding.Config{
		ID:         "someCrudId",
		Name:       "someCrudName",
		Group:      "someCrudOtherGroup",
		Properties: properties,
	}

	var confPutRes listing.Config
	_, err = testClient.put("config/someCrudOtherGroup/someCrudId", &confPutReq, &confPutRes)
	test.AssertNotError(t, err)
	test.AssertEqual(t, confPutRes.ID, "someCrudId")
	test.AssertEqual(t, confPutRes.Name, "someCrudName")
	test.AssertEqual(t, confPutRes.Group, "someCrudOtherGroup")

	var propMap map[string]interface{}
	err = json.Unmarshal(confPutRes.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)

	var confGetRes listing.Config
	_, err = testClient.get("config/someCrudOtherGroup/someCrudId", &confGetRes)
	test.AssertNotError(t, err)
	test.AssertEqual(t, confGetRes.ID, "someCrudId")
	test.AssertEqual(t, confGetRes.Name, "someCrudName")
	test.AssertEqual(t, confGetRes.Group, "someCrudOtherGroup")

	err = json.Unmarshal(confGetRes.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)
}

func Test_CRUDAddConfig_GroupNotFound(t *testing.T) {
	confPutReq := &adding.Config{
		ID:         "someCrudId",
		Name:       "someCrudName",
		Group:      "someCrudNonExistingGroup",
	}

	var confPutRes listing.Config
	res, _ := testClient.put("config/someCrudNonExistingGroup/someCrudId", &confPutReq, &confPutRes)
	payload, err := ioutil.ReadAll(res.Body)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res.StatusCode, http.StatusNotFound)
	test.AssertEqual(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, string(payload), listing.ErrGroupNotFound.Error()+"\n")
}

func Test_CRUDRetrieveConfig_GroupNotFound(t *testing.T) {
	var confGetRes listing.Config
	res, _ := testClient.get("config/someCrudNonExistingGroup/someCrudId", &confGetRes)
	payload, err := ioutil.ReadAll(res.Body)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res.StatusCode, http.StatusNotFound)
	test.AssertEqual(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, string(payload), listing.ErrGroupNotFound.Error()+"\n")
}

func Test_CRUDRetrieveConfig_ConfigNotFound(t *testing.T) {
	grpPutReq := adding.Group{ID: "someCrudGroup"}

	var grpPutRes listing.Group
	_, err := testClient.put("config/someCrudOtherGroup", &grpPutReq, &grpPutRes)

	var confGetRes listing.Config
	res, _ := testClient.get("config/someCrudGroup/someCrudNonExistingId", &confGetRes)
	payload, err := ioutil.ReadAll(res.Body)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res.StatusCode, http.StatusNotFound)
	test.AssertEqual(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	test.AssertEqual(t, string(payload), listing.ErrConfigNotFound.Error()+"\n")
}

func (c *TestClient) get(path string, responseObj interface{}) (*http.Response, error) {
	request, err := c.getRequest(path, http.MethodGet, nil, responseObj)
	if err != nil {
		return nil, err
	}

	return c.do(request, responseObj)
}

func (c *TestClient) put(path string, requestObj interface{}, responseObj interface{}) (*http.Response, error) {
	request, err := c.getRequest(path, http.MethodPut, requestObj, responseObj)
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	return c.do(request, responseObj)
}

func (c *TestClient) post(path string, requestObj interface{}, responseObj interface{}) (*http.Response, error) {
	request, err := c.getRequest(path, http.MethodPost, requestObj, responseObj)
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	return c.do(request, responseObj)
}

func (c *TestClient) getRequest(path string, method string, requestObj interface{}, responseObj interface{}) (*http.Request, error) {
	payload, err := json.Marshal(requestObj)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, c.baseURL+path, bytes.NewBuffer(payload))
	if err != nil {
		return request, err
	}

	return request, err
}

func (c *TestClient) do(req *http.Request, responseObj interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, err
	}

	defer resp.Body.Close()

	if responseObj != nil {
		if w, ok := responseObj.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(responseObj)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, err
}
