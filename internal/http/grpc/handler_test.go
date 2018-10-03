package grpc

import (
	"context"
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/internal/repository/memory"
	"github.com/larwef/ki/test"
	"io/ioutil"
	"testing"
)

var testDataFolder = "../../../test/testdata/"

func TestHandler_StoreGroup(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	req := &StoreGroupRequest{Id: "someGroup"}
	ctx := context.Background()
	res, err := handler.StoreGroup(ctx, req)

	test.AssertNotError(t, err)
	test.AssertEqual(t, res.Id, "someGroup")
	test.AssertEqual(t, len(res.ConfigIds), 0)
}

func TestHandler_StoreGroup_Duplicate(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	req := &StoreGroupRequest{Id: "someGroup"}
	ctx := context.Background()
	handler.StoreGroup(ctx, req)
	_, err := handler.StoreGroup(ctx, req)

	test.AssertEqual(t, err, adding.ErrGroupConflict)
}

func TestHandler_RetrieveGroup(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	ctx := context.Background()
	req := &RetrieveGroupRequest{Id: "someGroup"}
	res, err := handler.RetrieveGroup(ctx, req)

	test.AssertNotError(t, err)
	test.AssertEqual(t, res.Id, "someGroup")
	test.AssertEqual(t, len(res.ConfigIds), 3)
}

func TestHandler_RetrieveGroup_GroupNotFound(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3"},
	})

	ctx := context.Background()
	req := &RetrieveGroupRequest{Id: "someOtherGroup"}
	_, err := handler.RetrieveGroup(ctx, req)

	test.AssertEqual(t, err, listing.ErrGroupNotFound)
}

func TestHandler_StoreConfig(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	ctx := context.Background()
	properties, err := ioutil.ReadFile(testDataFolder + "properties.json")
	test.AssertNotError(t, err)
	req := &StoreConfigRequest{
		Id:         "someId",
		Name:       "someName",
		Group:      "someGroup",
		Properties: properties,
	}
	res, err := handler.StoreConfig(ctx, req)

	test.AssertNotError(t, err)
	test.AssertEqual(t, res.Id, "someId")
	test.AssertEqual(t, res.Name, "someName")
	test.AssertEqual(t, res.Group, "someGroup")

	var propMap map[string]interface{}
	err = json.Unmarshal(res.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)
}

func TestHandler_StoreConfig_GroupNotFound(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	ctx := context.Background()
	properties, err := ioutil.ReadFile(testDataFolder + "properties.json")
	test.AssertNotError(t, err)
	req := &StoreConfigRequest{
		Id:         "someId",
		Name:       "someName",
		Group:      "someOtherGroup",
		Properties: properties,
	}
	_, err = handler.StoreConfig(ctx, req)

	test.AssertEqual(t, err, listing.ErrGroupNotFound)
}

func TestHandler_RetrieveConfig(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	var c adding.Config
	test.UnmarshalJSONFromFile(t, testDataFolder+"configExample.json", &c)
	repository.StoreConfig(c)

	ctx := context.Background()
	req := &RetrieveConfigRequest{
		Id:      "someId",
		GroupId: "someGroup",
	}
	res, err := handler.RetrieveConfig(ctx, req)

	test.AssertNotError(t, err)
	test.AssertEqual(t, res.Id, "someId")
	test.AssertEqual(t, res.Name, "someName")
	test.AssertEqual(t, res.Group, "someGroup")

	var propMap map[string]interface{}
	err = json.Unmarshal(res.Properties, &propMap)
	test.AssertNotError(t, err)
	test.AssertEqual(t, propMap["property1"], float64(12))
	test.AssertEqual(t, propMap["property2"], "12")
	test.AssertEqual(t, propMap["property3"], "someString")
	test.AssertEqual(t, propMap["property4"], "someOtherString")
	test.AssertEqual(t, propMap["property5"], 12.1)
}

func TestHandler_RetrieveConfig_GroupNotFound(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	ctx := context.Background()
	req := &RetrieveConfigRequest{
		Id:      "someId",
		GroupId: "someOtherGroup",
	}
	_, err := handler.RetrieveConfig(ctx, req)

	test.AssertEqual(t, err, listing.ErrGroupNotFound)
}

func TestHandler_RetrieveConfig_ConfigNotFound(t *testing.T) {
	repository := memory.NewRepository()
	handler := NewHandler(adding.NewService(repository), listing.NewService(repository))

	repository.StoreGroup(adding.Group{
		ID: "someGroup",
	})

	ctx := context.Background()
	req := &RetrieveConfigRequest{
		Id:      "someOtherId",
		GroupId: "someGroup",
	}
	_, err := handler.RetrieveConfig(ctx, req)

	test.AssertEqual(t, err, listing.ErrConfigNotFound)
}
