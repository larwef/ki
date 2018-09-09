package persistence

import "github.com/larwef/ki/config"

// Persistence defines an interface for object who handles storing and retrieving object from some persisted medium
type Persistence interface {
	Store(c config.Config) error
	Retrieve(group string, id string) (*config.Config, error)
}
