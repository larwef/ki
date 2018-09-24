package commontests

import (
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/listing"
	"github.com/larwef/ki/testutil"
	"testing"
)

// TODO: Check if this is according to best practise. Its certainly useful to have common test for different repositories, but they should perhaps be structured differently.

// Repo has to satisfy adding and listing repository. Used to make tests that can cover all (or at least several) implementations
// of repository
type Repo interface {
	adding.Repository
	listing.Repository
}

func StoreAndRetrieveGroup(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()
	grp := adding.Group{
		ID: "someGroup",
	}

	err := repo.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	retrieveGrp, err := repo.RetrieveGroup("someGroup")
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, retrieveGrp.ID, grp.ID)
}

func StoreGroup_NotOverWrite(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID:      "someGroup",
		Configs: []string{"config1", "config2", "config3", "config4"},
	}

	err := repo.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	_, err = repo.RetrieveGroup("someGroup")
	testutil.AssertNotError(t, err)

	newGrp := adding.Group{
		ID:      "someGroup",
		Configs: []string{},
	}

	err = repo.StoreGroup(newGrp)
	testutil.AssertIsError(t, err)
	testutil.AssertEqual(t, err, adding.ErrGroupConflict)
}

func RetrieveGroup_GroupNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	err := repo.StoreGroup(grp)
	testutil.AssertNotError(t, err)

	_, err = repo.RetrieveGroup("someOtherGroup")
	testutil.AssertEqual(t, err, listing.ErrGroupNotFound)
}

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
	testutil.AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	retrieveConf, err := repo.RetrieveConfig("someGroup", "someConfig")
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, retrieveConf.ID, conf.ID)
	testutil.AssertEqual(t, retrieveConf.Group, "someGroup")
}

func StoreConfig_NoDuplicateInGroup(t *testing.T, repo Repo, cleanup func()) {
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
	testutil.AssertNotError(t, err)

	err = repo.StoreConfig(conf1)
	testutil.AssertNotError(t, err)
	err = repo.StoreConfig(conf2)
	testutil.AssertNotError(t, err)
	err = repo.StoreConfig(conf3)
	testutil.AssertNotError(t, err)

	newGrp, err := repo.RetrieveGroup("someGroup")
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, newGrp.ID, grp.ID)
	testutil.AssertEqual(t, len(newGrp.Configs), 2)
}

func StoreConfig_GroupNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf := adding.Config{
		ID:    "someConfig",
		Group: "someOtherGroup",
	}

	err := repo.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	testutil.AssertEqual(t, err, listing.ErrGroupNotFound)
}

func RetrieveConfig_GroupNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf := adding.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := repo.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	_, err = repo.RetrieveConfig("someOtherGroup", "someConfig")
	testutil.AssertEqual(t, err, listing.ErrGroupNotFound)
}

func RetrieveConfig_ConfigNotExist(t *testing.T, repo Repo, cleanup func()) {
	defer cleanup()

	grp := adding.Group{
		ID: "someGroup",
	}

	conf := adding.Config{
		ID:    "someConfig",
		Group: "someGroup",
	}

	err := repo.StoreGroup(grp)
	testutil.AssertNotError(t, err)
	err = repo.StoreConfig(conf)
	testutil.AssertNotError(t, err)

	_, err = repo.RetrieveConfig("someGroup", "someOtherConfig")
	testutil.AssertEqual(t, err, listing.ErrConfigNotFound)
}
