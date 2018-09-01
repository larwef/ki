package persistence

import "github.com/larwef/ki/config"

type persistence interface {
	Store(c config.Config) error
	Retrieve(id string) (*config.Config, error)
}
