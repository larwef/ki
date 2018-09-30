package grpc

import (
	"context"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
	"time"
)

// Handler handles processing of gRPC calls
type Handler struct {
	adding  adding.Service
	listing listing.Service
}

// NewHandler returns a new Handler
func NewHandler(adding adding.Service, listing listing.Service) *Handler {
	return &Handler{
		adding:  adding,
		listing: listing,
	}
}

// StoreGroup maps a request to a group object and stores it in the repository. Subsequently fetches the object and returns it
// to the caller.
func (s *Handler) StoreGroup(ctx context.Context, req *StoreGroupRequest) (*Group, error) {
	addGrp := adding.Group{ID: req.Id}

	if err := s.adding.AddGroup(addGrp); err != nil {
		return &Group{}, err
	}

	return s.retrieveGroup(req.Id)
}

// RetrieveGroup fetches a group object from repository and maps it to a gRPC response
func (s *Handler) RetrieveGroup(ctx context.Context, req *RetrieveGroupRequest) (*Group, error) {
	return s.retrieveGroup(req.Id)
}

func (s *Handler) retrieveGroup(groupID string) (*Group, error) {
	grp, err := s.listing.GetGroup(groupID)

	return &Group{
		Id:        grp.ID,
		ConfigIds: grp.Configs,
	}, err
}

// StoreConfig maps a request to a config object and stores it in the repository. Subsequently fetches the object and returns it
// to the caller.
func (s *Handler) StoreConfig(ctx context.Context, req *StoreConfigRequest) (*Config, error) {
	addConf := adding.Config{
		ID:           req.Id,
		Name:         req.Name,
		LastModified: time.Now(),
		Group:        req.Group,
		Properties:   req.Properties,
	}

	if err := s.adding.AddConfig(addConf); err != nil {
		return &Config{}, err
	}

	return s.retrieveConfig(req.Group, req.Id)
}

// RetrieveConfig fetches a config object from repository and maps it to a gRPC response
func (s *Handler) RetrieveConfig(ctx context.Context, req *RetrieveConfigRequest) (*Config, error) {
	return s.retrieveConfig(req.GroupId, req.Id)
}

func (s *Handler) retrieveConfig(groupID string, configID string) (*Config, error) {
	conf, err := s.listing.GetConfig(groupID, configID)

	return &Config{
		Id:           conf.ID,
		Name:         conf.Name,
		LastModified: conf.LastModified.Unix(),
		Version:      int32(conf.Version),
		Group:        conf.Group,
		Properties:   conf.Properties,
	}, err
}
