package memory

import (
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/listing"
)

// Repository representa a in memory storge object
type Repository struct {
	groups  map[string]Group
	configs map[string]Config
}

// NewRepository returns a new Repository storage object
func NewRepository() *Repository {
	return &Repository{
		groups:  make(map[string]Group),
		configs: make(map[string]Config),
	}
}

// StoreGroup stores a config in the memory storage. Will not overwrite existing Group object.
func (r *Repository) StoreGroup(g adding.Group) error {
	if _, exists := r.groups[g.ID]; exists {
		return adding.ErrGroupConflict
	}

	r.groups[g.ID] = Group{
		ID:      g.ID,
		Configs: g.Configs,
	}

	return nil
}

// RetrieveGroup retrieves a group from the memory storage specified by id
func (r *Repository) RetrieveGroup(id string) (*listing.Group, error) {
	if val, exists := r.groups[id]; exists {
		return &listing.Group{
			ID:      val.ID,
			Configs: val.Configs,
		}, nil
	}

	return &listing.Group{}, listing.ErrGroupNotFound
}

// StoreConfig stores a config in the memory storage
func (r *Repository) StoreConfig(c adding.Config) error {
	grp, exists := r.groups[c.Group]
	if !exists {
		return listing.ErrGroupNotFound
	}

	if len(grp.Configs) == 0 {
		grp.Configs = append(grp.Configs, c.ID)
	} else {
		for i := 0; i <= len(grp.Configs); i++ {
			if grp.Configs[i] == c.ID {
				break
			}

			if i >= len(grp.Configs)-1 {
				grp.Configs = append(grp.Configs, c.ID)
				break
			}
		}
	}

	r.groups[c.Group] = grp
	r.configs[c.ID] = Config{
		ID:           c.ID,
		Name:         c.Name,
		LastModified: c.LastModified,
		Version:      c.Version,
		Group:        c.Group,
		Properties:   c.Properties,
	}

	return nil
}

// RetrieveConfig retrieves a config from the memory storage spesified by groupID and id of the config
// TODO: Doesnt behave same way as local. Cant handle configs with same id and different group
func (r *Repository) RetrieveConfig(groupID string, id string) (*listing.Config, error) {

	if _, exists := r.groups[groupID]; !exists {
		return &listing.Config{}, listing.ErrGroupNotFound
	}

	c, exists := r.configs[id]
	if !exists {
		return &listing.Config{}, listing.ErrConfigNotFound
	}

	return &listing.Config{
		ID:           c.ID,
		Name:         c.Name,
		LastModified: c.LastModified,
		Version:      c.Version,
		Group:        c.Group,
		Properties:   c.Properties,
	}, nil
}
