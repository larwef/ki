package repository

import (
	"github.com/larwef/ki/testutil"
	"os"
	"testing"
)

var testDir = "testdir"

func TestLocal_StoreAndRetrieveGroup(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	retrieveGrp, err := lcl.RetrieveGroup("someGroup")
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, retrieveGrp.ID, grp.ID)
}

func TestLocal_StoreGroup_NotOverWrite(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3", "config4"},
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	_, err = lcl.RetrieveGroup("someGroup")
	testutil.AssertNotError(t, err)

	newGrp := Group{
		ID:      "someGroup",
		Configs: []string{},
	}

	err = lcl.StoreGroup(newGrp)
	testutil.AssertIsError(t, err)
	testutil.AssertEqual(t, err, ErrConflict)
}

func TestLocal_RetrieveGroup_GroupNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	_, err = lcl.RetrieveGroup("someOtherGroup")
	testutil.AssertEqual(t, err, ErrGroupNotFound)
}

func TestLocal_StoreAndRetrieveConfig(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID: "someGroup",
	}

	conf := Config{
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

func TestLocal_StoreConfig_NoDuplicateInGroup(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID: "someGroup",
	}

	conf1 := Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	conf2 := Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	// This is essentialy an update of conf1
	conf3 := Config{
		ID:    "someOtherConfig",
		Group: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	err = lcl.StoreConfig(conf1)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf2)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf3)
	testutil.AssertNotError(t, err)

	newGrp, err := lcl.RetrieveGroup("someGroup")
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, newGrp.ID, grp.ID)
	testutil.AssertEqual(t, len(newGrp.Configs), 2)
}

func TestLocal_StoreConfig_GroupNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID: "someGroup",
	}

	conf := Config{
		ID:    "someConfig",
		Group: "someOtherGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf)
	testutil.AssertEqual(t, err, ErrGroupNotFound)
}

func TestLocal_RetrieveConfig_GroupNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID: "someGroup",
	}

	conf := Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	_, err = lcl.RetrieveConfig("someOtherGroup", "someConfig")
	testutil.AssertEqual(t, err, ErrGroupNotFound)
}

func TestLocal_RetrieveConfig_ConfigNotExist(t *testing.T) {
	defer clean()
	lcl := NewLocal(testDir)

	grp := Group{
		ID: "someGroup",
	}

	conf := Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := lcl.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = lcl.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	_, err = lcl.RetrieveConfig("someGroup", "someOtherConfig")
	testutil.AssertEqual(t, err, ErrConfigNotFound)
}

func clean() {
	os.RemoveAll(testDir)
}
