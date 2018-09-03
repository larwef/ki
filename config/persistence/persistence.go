package persistence

import "github.com/larwef/ki/config"

type Persistence interface {
	Store(c config.Config) error
	Retrieve(id string, group string) (*config.Config, error)
}
