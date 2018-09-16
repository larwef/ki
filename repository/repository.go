package repository

import (
	"errors"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/group"
)

// TODO: Make error type with statuscode etc

// ErrGroupNotFound is returned when group cant be found
var ErrGroupNotFound = errors.New("group not found")

// ErrConfigNotFound is returned when config cant be found
var ErrConfigNotFound = errors.New("config not found")

// Repository is an interface defining behaviour for persisting
type Repository interface {
	StoreGroup(g group.Group) error
	RetrieveGroup(id string) (*group.Group, error)

	StoreConfig(c config.Config) error
	RetrieveConfig(group string, id string) (*config.Config, error)
}

// Mock is mocking the Repository interface for testing
type Mock struct {
	StoredGroup  group.Group
	StoredConfig config.Config
}

// StoreGroup stores Group object to the mock
func (r *Mock) StoreGroup(g group.Group) error {
	r.StoredGroup = g

	return nil
}

// RetrieveGroup retrieves Group object from the mock
func (r *Mock) RetrieveGroup(id string) (*group.Group, error) {
	if id != r.StoredGroup.ID {
		return nil, ErrGroupNotFound
	}

	return &r.StoredGroup, nil
}

// StoreConfig stores Config object to the mock
func (r *Mock) StoreConfig(c config.Config) error {
	grp, err := r.RetrieveGroup(c.Group)
	if err != nil {
		return err
	}

	r.StoredConfig = c

	grp.Configs = append(grp.Configs, c.ID)
	err = r.StoreGroup(*grp)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveConfig retrieves Config object from the mock
func (r *Mock) RetrieveConfig(groupID string, id string) (*config.Config, error) {
	if groupID != r.StoredGroup.ID {
		return nil, ErrGroupNotFound
	}

	if id != r.StoredConfig.ID {
		return nil, ErrConfigNotFound
	}

	return &r.StoredConfig, nil
}
