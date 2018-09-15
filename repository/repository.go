package repository

import (
	"errors"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/group"
)

var errGroupNotFound = errors.New("group not found")
var errConfigNotFound = errors.New("config not found")

type Repository interface {
	StoreGroup(g group.Group) error
	RetrieveGroup(id string) (*group.Group, error)

	StoreConfig(c config.Config) error
	RetrieveConfig(group string, id string) (*config.Config, error)
}

// Mock is mocking the Repository interfce for testing
type Mock struct {
	StoredGroup  group.Group
	StoredConfig config.Config
}

func (r *Mock) StoreGroup(g group.Group) error {
	r.StoredGroup = g
	return nil
}

func (r *Mock) RetrieveGroup(id string) (*group.Group, error) {
	if id != r.StoredGroup.ID {
		return nil, errGroupNotFound
	}

	return &r.StoredGroup, nil
}

func (r *Mock) StoreConfig(c config.Config) error {
	r.StoredConfig = c
	return nil
}

func (r *Mock) RetrieveConfig(groupId string, id string) (*config.Config, error) {
	if groupId != r.StoredConfig.Group {
		return nil, errGroupNotFound
	}

	if id != r.StoredConfig.ID {
		return nil, errConfigNotFound
	}

	return &r.StoredConfig, nil
}
