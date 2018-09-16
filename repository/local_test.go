package repository

import (
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/group"
	"github.com/larwef/ki/testutil"
	"os"
	"testing"
)

var testDir = "testdir"

func TestLocal_StoreAndRetrieveGroup(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := group.Group{
		ID: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	retrieveGrp, err := lcl.RetrieveGroup("someGroup")
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, retrieveGrp.ID, grp.ID)
}

func TestLocal_RetrieveGroup_GroupNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := group.Group{
		ID: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	_, err = lcl.RetrieveGroup("someOtherGroup")
	testutil.AssertEqual(t, err.Error(), ErrGroupNotFound.Error())
}

func TestLocal_StoreAndRetrieveConfig(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := group.Group{
		ID: "someGroup",
	}

	conf := config.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	retrieveConf, err := lcl.RetrieveConfig("someGroup", "someConfig")
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, retrieveConf.ID, conf.ID)
	testutil.AssertEqual(t, retrieveConf.Group, "someGroup")
}

func TestLocal_StoreConfig_GroupNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := group.Group{
		ID: "someGroup",
	}

	conf := config.Config{
		ID:    "someConfig",
		Group: "someOtherGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf)
	testutil.AssertEqual(t, err.Error(), ErrGroupNotFound.Error())
}

func TestLocal_RetrieveConfig_GroupNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := group.Group{
		ID: "someGroup",
	}

	conf := config.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	_, err = lcl.RetrieveConfig("someOtherGroup", "someConfig")
	testutil.AssertEqual(t, err.Error(), ErrGroupNotFound.Error())
}

func TestLocal_RetrieveConfig_ConfigNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := group.Group{
		ID: "someGroup",
	}

	conf := config.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	_, err = lcl.RetrieveConfig("someGroup", "someOtherConfig")
	testutil.AssertEqual(t, err.Error(), ErrConfigNotFound.Error())
}

func clean() {
	os.RemoveAll(testDir)
}
