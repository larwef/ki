package test

import (
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/listing"
	"testing"
)

// Repo has to satisfy adding and listing repository. Used to make tests that can cover all (or at least several) implementations
// of repository
type Repo interface {
	adding.Repository
	listing.Repository
}

// StoreAndRetrieveGroup tests that a group object can be stored and subsequently retrieved
func StoreAndRetrieveGroup(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()
	grp := adding.Group{
		ID: "someGroup",
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)

	retrieveGrp, err := repo.RetrieveGroup("someGroup")
	AssertNotError(t, err)
	AssertEqual(t, retrieveGrp.ID, grp.ID)
}

// StoreGroupAndNotOverWrite tests that when store group is called on an existing group it will not overwrite the existing group
func StoreGroupAndNotOverWrite(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3", "config4"},
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)

	_, err = repo.RetrieveGroup("someGroup")
	AssertNotError(t, err)

	newGrp := adding.Group{
		ID:      "someGroup",
		Configs: []string{},
	}

	err = repo.StoreGroup(newGrp)
	AssertIsError(t, err)
	AssertEqual(t, err, adding.ErrGroupConflict)
}

// RetrieveGroupWhenGroupNotExist tests retrieving a group that doesnt exist
func RetrieveGroupWhenGroupNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)

	_, err = repo.RetrieveGroup("someOtherGroup")
	AssertEqual(t, err, listing.ErrGroupNotFound)
}

// StoreAndRetrieveConfig tests storing and subsequently retrieving a Config
func StoreAndRetrieveConfig(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf := adding.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	AssertNotError(t, err)

	retrieveConf, err := repo.RetrieveConfig("someGroup", "someConfig")
	AssertNotError(t, err)
	AssertEqual(t, retrieveConf.ID, conf.ID)
	AssertEqual(t, retrieveConf.Group, "someGroup")
}

// StoreConfigNoDuplicateInGroup tests that when adding a group with the same id, the Group object wont have duplicates in its Config array
func StoreConfigNoDuplicateInGroup(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf1 := adding.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	conf2 := adding.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	// This is essentialy an update of conf1
	conf3 := adding.Config{
		ID:    "someOtherConfig",
		Group: "someGroup",
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)

	err = repo.StoreConfig(conf1)
	AssertNotError(t, err)
	err = repo.StoreConfig(conf2)
	AssertNotError(t, err)
	err = repo.StoreConfig(conf3)
	AssertNotError(t, err)

	newGrp, err := repo.RetrieveGroup("someGroup")
	AssertNotError(t, err)
	AssertEqual(t, newGrp.ID, grp.ID)
	AssertEqual(t, len(newGrp.Configs), 2)
}

// StoreConfigWhenGroupNotExist tests storing a Config when the Group doesnt exist
func StoreConfigWhenGroupNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf := adding.Config{
		ID:    "someConfig",
		Group: "someOtherGroup",
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	AssertEqual(t, err, listing.ErrGroupNotFound)
}

// RetrieveConfigWhenGroupNotExist tests retireving a Config with a non existing Group
func RetrieveConfigWhenGroupNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf := adding.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	AssertNotError(t, err)

	_, err = repo.RetrieveConfig("someOtherGroup", "someConfig")
	AssertEqual(t, err, listing.ErrGroupNotFound)
}

// RetrieveConfigWhenConfigNotExist tests retireving a non existing Config
func RetrieveConfigWhenConfigNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf := adding.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := repo.StoreGroup(grp)
	AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	AssertNotError(t, err)

	_, err = repo.RetrieveConfig("someGroup", "someOtherConfig")
	AssertEqual(t, err, listing.ErrConfigNotFound)
}
