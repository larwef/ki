package repository

import (
	"fmt"
	"net/http"
)

// Error defines an error type for repository errors
type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

// ErrGroupNotFound is returned when group cant be found
var ErrGroupNotFound = Error{
	Code:    http.StatusNotFound,
	Message: "Group not found",
}

// ErrConfigNotFound is returned when config cant be found
var ErrConfigNotFound = Error{
	Code:    http.StatusNotFound,
	Message: "Config not found",
}

// ErrInternal is returned when there was an error from the underlying storage. Eg. error persisting or retrieving
var ErrInternal = Error{
	Code:    http.StatusInternalServerError,
	Message: "Internal server error",
}

// ErrConflict is returned when an object already exists thats not overwritable
var ErrConflict = Error{
	Code:    http.StatusConflict,
	Message: "Resource already exist and is not overwritable",
}

// Repository is an interface defining behaviour for persisting
type Repository interface {
	StoreGroup(g Group) error
	RetrieveGroup(id string) (*Group, error)

	StoreConfig(c Config) error
	RetrieveConfig(group string, id string) (*Config, error)
}
