package local

import (
	"github.com/larwef/ki/pkg/repository/commontests"
	"os"
	"testing"
)

var testDir = "testdir"

func TestRepository_StoreAndRetrieveGroup(t *testing.T) {
	commontests.StoreAndRetrieveGroup(t, NewRepository(testDir), clean)
}

func TestRepository_StoreGroup_NotOverWrite(t *testing.T) {
	commontests.StoreGroup_NotOverWrite(t, NewRepository(testDir), clean)
}

func TestRepository_RetrieveGroup_GroupNotExist(t *testing.T) {
	commontests.RetrieveGroup_GroupNotExist(t, NewRepository(testDir), clean)
}

func TestRepository_StoreAndRetrieveConfig(t *testing.T) {
	commontests.StoreAndRetrieveConfig(t, NewRepository(testDir), clean)
}

func TestRepository_StoreConfig_NoDuplicateInGroup(t *testing.T) {
	commontests.StoreConfig_NoDuplicateInGroup(t, NewRepository(testDir), clean)
}

func TestRepository_StoreConfig_GroupNotExist(t *testing.T) {
	commontests.StoreConfig_GroupNotExist(t, NewRepository(testDir), clean)
}

func TestRepository_RetrieveConfig_GroupNotExist(t *testing.T) {
	commontests.RetrieveConfig_GroupNotExist(t, NewRepository(testDir), clean)
}

func TestRepository_RetrieveConfig_ConfigNotExist(t *testing.T) {
	commontests.RetrieveConfig_ConfigNotExist(t, NewRepository(testDir), clean)
}

func clean() {
	os.RemoveAll(testDir)
}
