package repository

import (
	"errors"
)

// TODO: Make error type with statuscode etc

// ErrGroupNotFound is returned when group cant be found
var ErrGroupNotFound = errors.New("group not found")

// ErrConfigNotFound is returned when config cant be found
var ErrConfigNotFound = errors.New("config not found")

// Repository is an interface defining behaviour for persisting
type Repository interface {
	StoreGroup(g Group) error
	RetrieveGroup(id string) (*Group, error)

	StoreConfig(c Config) error
	RetrieveConfig(group string, id string) (*Config, error)
}
