package persistence

import (
	"github.com/larwef/ki/config"
	"github.com/pkg/errors"
)

// Mock is used for mocking an object implementing the persistence interface during testing
type Mock struct {
	storedConfig config.Config
}

// NewMock returns a Mock object for testing
func NewMock(config config.Config) *Mock {
	return &Mock{storedConfig: config}
}

// Store stores a config object in the mock
func (l *Mock) Store(c config.Config) error {
	l.storedConfig = c
	return nil
}

// Retrieve retrieves a config object from the mock
func (l *Mock) Retrieve(group string, id string) (*config.Config, error) {
	if group != l.storedConfig.Group || id != l.storedConfig.ID {
		return nil, errors.New("Config not found")
	}

	return &l.storedConfig, nil
}
