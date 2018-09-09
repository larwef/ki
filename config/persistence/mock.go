package persistence

import (
	"github.com/larwef/ki/config"
	"github.com/pkg/errors"
)

type mock struct {
	storedConfig config.Config
}

func NewMock(config config.Config) *mock {
	return &mock{storedConfig: config}
}

func (l *mock) Store(c config.Config) error {
	l.storedConfig = c
	return nil
}

func (l *mock) Retrieve(group string, id string) (*config.Config, error) {
	if group != l.storedConfig.Group || id != l.storedConfig.Id {
		return nil, errors.New("Config not found")
	}

	return &l.storedConfig, nil
}
