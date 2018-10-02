// +build integration

package integration

import (
	"context"
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/http/grpc"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/test"
	goGrpc "google.golang.org/grpc"
	"io/ioutil"
	"strings"
	"testing"
)

var address = "localhost:8081"
var testDataFolder = "../testdata/"

// Integration tests needs a fresh instance running locally to work. The easiest is to run the ./test-docker.sh script.

func getConnection(t *testing.T) *goGrpc.ClientConn {
	// Set up a connection to the server.
	conn, err := goGrpc.Dial(address, goGrpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	return conn
}

func Test_AddAndRetrieveGroup(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	groupClient := grpc.NewGroupServiceClient(conn)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someGroup"}

	storeGrpRes, err := groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, storeGrpRes.Id, "someGroup")
	test.AssertEqual(t, len(storeGrpRes.ConfigIds), 0)

	retrieveGrpReq := &grpc.RetrieveGroupRequest{Id: "someGroup"}

	retrieveGroupRes, err := groupClient.RetrieveGroup(context.Background(), retrieveGrpReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, retrieveGroupRes.Id, "someGroup")
	test.AssertEqual(t, len(retrieveGroupRes.ConfigIds), 0)
}

func Test_AddGroup_Duplicate(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	groupClient := grpc.NewGroupServiceClient(conn)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someGroupConflict"}

	storeGrpRes, err := groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, storeGrpRes.Id, "someGroupConflict")
	test.AssertEqual(t, len(storeGrpRes.ConfigIds), 0)

	storeGrpRes, err = groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertEqual(t, strings.Contains(err.Error(), adding.ErrGroupConflict.Error()), true)
}

func Test_AddAndRetrieveConfig(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	groupClient := grpc.NewGroupServiceClient(conn)
	configClient := grpc.NewConfigServiceClient(conn)

	properties, err := ioutil.ReadFile(testDataFolder + "properties.json")
	test.AssertNotError(t, err)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someOtherGroup"}

	_, err = groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)

	storeConfigreq := &grpc.StoreConfigRequest{
		Id:         "someId",
		Name:       "someName",
		Group:      "someOtherGroup",
		Properties: properties,
	}

	storeConfigRes, err := configClient.StoreConfig(context.Background(), storeConfigreq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, storeConfigRes.Id, "someId")
	test.AssertEqual(t, storeConfigRes.Name, "someName")
	test.AssertEqual(t, storeConfigRes.Group, "someOtherGroup")

	var propMap map[string]interface{}
	err = json.Unmarshal(storeConfigRes.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)

	retrieveConfigReq := &grpc.RetrieveConfigRequest{
		Id:      "someId",
		GroupId: "someOtherGroup",
	}

	retrieveConfigRes, err := configClient.RetrieveConfig(context.Background(), retrieveConfigReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, retrieveConfigRes.Id, "someId")
	test.AssertEqual(t, retrieveConfigRes.Name, "someName")
	test.AssertEqual(t, retrieveConfigRes.Group, "someOtherGroup")

	err = json.Unmarshal(retrieveConfigRes.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)
}

func Test_AddConfig_GroupNotFound(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	configClient := grpc.NewConfigServiceClient(conn)

	storeConfigReq := &grpc.StoreConfigRequest{
		Id:         "someId",
		Name:       "someName",
		Group:      "someNonExistingGroup",
	}

	_, err := configClient.StoreConfig(context.Background(), storeConfigReq)
	test.AssertEqual(t, strings.Contains(err.Error(), listing.ErrGroupNotFound.Error()), true)
}

func Test_RetrieveConfig_GroupNotFound(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	configClient := grpc.NewConfigServiceClient(conn)

	retrieveConfigReq := &grpc.RetrieveConfigRequest{
		Id:      "someId",
		GroupId: "someNonExistingGroup",
	}

	_, err := configClient.RetrieveConfig(context.Background(), retrieveConfigReq)
	test.AssertEqual(t, strings.Contains(err.Error(), listing.ErrGroupNotFound.Error()), true)
}

func Test_RetrieveConfig_ConfigNotFound(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	configClient := grpc.NewConfigServiceClient(conn)
	groupClient := grpc.NewGroupServiceClient(conn)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someNewGroup"}

	_, err := groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)


	retrieveConfigReq := &grpc.RetrieveConfigRequest{
		Id:      "someNonExistingId",
		GroupId: "someNewGroup",
	}

	_, err = configClient.RetrieveConfig(context.Background(), retrieveConfigReq)
	test.AssertEqual(t, strings.Contains(err.Error(), listing.ErrConfigNotFound.Error()), true)
}
