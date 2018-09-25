package local

import (
	"github.com/larwef/ki/test"
	"os"
	"testing"
)

var testDir = "testdir"

func TestRepository_StoreAndRetrieveGroup(t *testing.T) {
	test.StoreAndRetrieveGroup(t, NewRepository(testDir), clean)
}

func TestRepository_StoreGroup_NotOverWrite(t *testing.T) {
	test.StoreGroupAndNotOverWrite(t, NewRepository(testDir), clean)
}

func TestRepository_RetrieveGroup_GroupNotExist(t *testing.T) {
	test.RetrieveGroupWhenGroupNotExist(t, NewRepository(testDir), clean)
}

func TestRepository_StoreAndRetrieveConfig(t *testing.T) {
	test.StoreAndRetrieveConfig(t, NewRepository(testDir), clean)
}

func TestRepository_StoreConfig_NoDuplicateInGroup(t *testing.T) {
	test.StoreConfigNoDuplicateInGroup(t, NewRepository(testDir), clean)
}

func TestRepository_StoreConfig_GroupNotExist(t *testing.T) {
	test.StoreConfigWhenGroupNotExist(t, NewRepository(testDir), clean)
}

func TestRepository_RetrieveConfig_GroupNotExist(t *testing.T) {
	test.RetrieveConfigWhenGroupNotExist(t, NewRepository(testDir), clean)
}

func TestRepository_RetrieveConfig_ConfigNotExist(t *testing.T) {
	test.RetrieveConfigWhenConfigNotExist(t, NewRepository(testDir), clean)
}

func clean() {
	os.RemoveAll(testDir)
}
