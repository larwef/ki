package persistence

import "github.com/larwef/ki/config"

type Persistence interface {
	Store(c config.Config) error
	Retrieve(group string, id string) (*config.Config, error)
}
