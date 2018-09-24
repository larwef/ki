package listing

import "errors"

// ErrGroupNotFound is used when a group resource could not be found.
var ErrGroupNotFound = errors.New("group not found")

// ErrConfigNotFound is used when a config resource could not be found.
var ErrConfigNotFound = errors.New("config not found")

// Service provides adding operations
type Service interface {
	GetGroup(id string) (*Group, error)
	GetConfig(groupID string, id string) (*Config, error)
}

// Repository provides access to repository
type Repository interface {
	RetrieveGroup(id string) (*Group, error)
	RetrieveConfig(groupID string, id string) (*Config, error)
}

type service struct {
	repo Repository
}

// NewService created a new adding service
func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetGroup(id string) (*Group, error) {
	return s.repo.RetrieveGroup(id)
}

func (s *service) GetConfig(groupID string, id string) (*Config, error) {
	return s.repo.RetrieveConfig(groupID, id)
}
