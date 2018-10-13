// +build integration

package integration

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/http/grpc"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/test"
	goGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"strings"
	"testing"
)

var grpcAddress = "tlstest.wefald.no:8081"
var grpcTestDataFolder = "../testdata/"

// Integration tests needs a fresh instance running locally to work. The easiest is to run the ./test-docker.sh script.

func getConnection(t *testing.T) *goGrpc.ClientConn {
	// Set up a connection to the server.

	// Use this if testing live with TLS
	rootCAs, _ := x509.SystemCertPool()
	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}
	creds := credentials.NewTLS(tlsConfig)
	conn, err := goGrpc.Dial(grpcAddress, goGrpc.WithTransportCredentials(creds))

	//conn, err := goGrpc.Dial(grpcAddress, goGrpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	return conn
}

func Test_gRPCAddAndRetrieveGroup(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	groupClient := grpc.NewGroupServiceClient(conn)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someGRPCGroup"}

	storeGrpRes, err := groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, storeGrpRes.Id, "someGRPCGroup")
	test.AssertEqual(t, len(storeGrpRes.ConfigIds), 0)

	retrieveGrpReq := &grpc.RetrieveGroupRequest{Id: "someGRPCGroup"}

	retrieveGroupRes, err := groupClient.RetrieveGroup(context.Background(), retrieveGrpReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, retrieveGroupRes.Id, "someGRPCGroup")
	test.AssertEqual(t, len(retrieveGroupRes.ConfigIds), 0)
}

func Test_gRPCAddGroup_Duplicate(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	groupClient := grpc.NewGroupServiceClient(conn)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someGRPCGroupConflict"}

	storeGrpRes, err := groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, storeGrpRes.Id, "someGRPCGroupConflict")
	test.AssertEqual(t, len(storeGrpRes.ConfigIds), 0)

	storeGrpRes, err = groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertEqual(t, strings.Contains(err.Error(), adding.ErrGroupConflict.Error()), true)
}

func Test_gRPCAddAndRetrieveConfig(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	groupClient := grpc.NewGroupServiceClient(conn)
	configClient := grpc.NewConfigServiceClient(conn)

	properties, err := ioutil.ReadFile(grpcTestDataFolder + "properties.json")
	test.AssertNotError(t, err)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someGRPCOtherGroup"}

	_, err = groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)

	storeConfigreq := &grpc.StoreConfigRequest{
		Id:         "someGRPCId",
		Name:       "someGRPCName",
		Group:      "someGRPCOtherGroup",
		Properties: properties,
	}

	storeConfigRes, err := configClient.StoreConfig(context.Background(), storeConfigreq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, storeConfigRes.Id, "someGRPCId")
	test.AssertEqual(t, storeConfigRes.Name, "someGRPCName")
	test.AssertEqual(t, storeConfigRes.Group, "someGRPCOtherGroup")

	var propMap map[string]interface{}
	err = json.Unmarshal(storeConfigRes.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)

	retrieveConfigReq := &grpc.RetrieveConfigRequest{
		Id:      "someGRPCId",
		GroupId: "someGRPCOtherGroup",
	}

	retrieveConfigRes, err := configClient.RetrieveConfig(context.Background(), retrieveConfigReq)
	test.AssertNotError(t, err)
	test.AssertEqual(t, retrieveConfigRes.Id, "someGRPCId")
	test.AssertEqual(t, retrieveConfigRes.Name, "someGRPCName")
	test.AssertEqual(t, retrieveConfigRes.Group, "someGRPCOtherGroup")

	err = json.Unmarshal(retrieveConfigRes.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)
}

func Test_gRPCAddConfig_GroupNotFound(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	configClient := grpc.NewConfigServiceClient(conn)

	storeConfigReq := &grpc.StoreConfigRequest{
		Id:         "someGRPCId",
		Name:       "someGRPCName",
		Group:      "someGRPCNonExistingGroup",
	}

	_, err := configClient.StoreConfig(context.Background(), storeConfigReq)
	test.AssertEqual(t, strings.Contains(err.Error(), listing.ErrGroupNotFound.Error()), true)
}

func Test_gRPCRetrieveConfig_GroupNotFound(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	configClient := grpc.NewConfigServiceClient(conn)

	retrieveConfigReq := &grpc.RetrieveConfigRequest{
		Id:      "someGRPCId",
		GroupId: "someGRPCNonExistingGroup",
	}

	_, err := configClient.RetrieveConfig(context.Background(), retrieveConfigReq)
	test.AssertEqual(t, strings.Contains(err.Error(), listing.ErrGroupNotFound.Error()), true)
}

func Test_gRPCRetrieveConfig_ConfigNotFound(t *testing.T) {
	conn := getConnection(t)
	defer conn.Close()

	configClient := grpc.NewConfigServiceClient(conn)
	groupClient := grpc.NewGroupServiceClient(conn)

	storeGrpReq := &grpc.StoreGroupRequest{Id: "someGRPCNewGroup"}

	_, err := groupClient.StoreGroup(context.Background(), storeGrpReq)
	test.AssertNotError(t, err)


	retrieveConfigReq := &grpc.RetrieveConfigRequest{
		Id:      "someGRPCNonExistingId",
		GroupId: "someGRPCNewGroup",
	}

	_, err = configClient.RetrieveConfig(context.Background(), retrieveConfigReq)
	test.AssertEqual(t, strings.Contains(err.Error(), listing.ErrConfigNotFound.Error()), true)
}
