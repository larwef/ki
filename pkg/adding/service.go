package adding

import "errors"

// ErrGroupConflict is used when a group already exists.
var ErrGroupConflict = errors.New("group already exists and is not overwritable")

// Service provides adding operations
type Service interface {
	AddGroup(g Group) error
	AddConfig(c Config) error
}

// Repository provides access to repository
type Repository interface {
	StoreGroup(g Group) error
	StoreConfig(c Config) error
}

type service struct {
	repo Repository
}

// NewService created a new adding service
func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) AddGroup(g Group) error {
	return s.repo.StoreGroup(Group{ID: g.ID})
}

func (s *service) AddConfig(c Config) error {
	return s.repo.StoreConfig(c)
}
