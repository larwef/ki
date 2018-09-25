package memory

import (
	"github.com/larwef/ki/test"
	"testing"
)

func TestRepository_StoreAndRetrieveGroup(t *testing.T) {
	test.StoreAndRetrieveGroup(t, NewRepository(), clean)
}

func TestRepository_StoreGroup_NotOverWrite(t *testing.T) {
	test.StoreGroupAndNotOverWrite(t, NewRepository(), clean)
}

func TestRepository_RetrieveGroup_GroupNotExist(t *testing.T) {
	test.RetrieveGroupWhenGroupNotExist(t, NewRepository(), clean)
}

func TestRepository_StoreAndRetrieveConfig(t *testing.T) {
	test.StoreAndRetrieveConfig(t, NewRepository(), clean)
}

func TestRepository_StoreConfig_NoDuplicateInGroup(t *testing.T) {
	test.StoreConfigNoDuplicateInGroup(t, NewRepository(), clean)
}

func TestRepository_StoreConfig_GroupNotExist(t *testing.T) {
	test.StoreConfigWhenGroupNotExist(t, NewRepository(), clean)
}

func TestRepository_RetrieveConfig_GroupNotExist(t *testing.T) {
	test.RetrieveConfigWhenGroupNotExist(t, NewRepository(), clean)
}

func TestRepository_RetrieveConfig_ConfigNotExist(t *testing.T) {
	test.RetrieveConfigWhenConfigNotExist(t, NewRepository(), clean)
}

func clean() {}
