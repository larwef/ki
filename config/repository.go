package config

import "errors"

// Repository defines an interface for object who handles storing and retrieving Config objects from some persisted medium
type Repository interface {
	Store(c Config) error
	Retrieve(group string, id string) (*Config, error)
}

// RepositoryMock is used for mocking an object implementing the persistence interface during testing
type RepositoryMock struct {
	storedConfig Config
}

// NewRepositoryMock returns a RepositoryMock object for testing
func NewRepositoryMock(config Config) *RepositoryMock {
	return &RepositoryMock{storedConfig: config}
}

// Store stores a config object in the mock
func (l *RepositoryMock) Store(c Config) error {
	l.storedConfig = c
	return nil
}

// Retrieve retrieves a config object from the mock
func (l *RepositoryMock) Retrieve(group string, id string) (*Config, error) {
	if group != l.storedConfig.Group || id != l.storedConfig.ID {
		return nil, errors.New("Config not found")
	}

	return &l.storedConfig, nil
}
