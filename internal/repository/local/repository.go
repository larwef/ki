package local

import (
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
	"log"
	"os"
)

// Repository representa a local storge object
type Repository struct {
	path string
}

// NewRepository returns a new Repository storage object
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

// StoreGroup stores a config in the local storage. Will not overwrite existing Group object.
func (r *Repository) StoreGroup(g adding.Group) error {
	fullPath := r.path + "/" + g.ID + ".json"
	if _, err := os.Stat(fullPath); err == nil {
		return adding.ErrGroupConflict
	}

	return r.storeGroup(g)
}

// This store function will overwrite the group object
func (r *Repository) storeGroup(g adding.Group) error {
	basePath := r.path + "/"
	fullPath := basePath + g.ID + ".json"

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	grp := Group{
		ID:      g.ID,
		Configs: g.Configs,
	}

	return storeJSON(file, grp)
}

// RetrieveGroup retrieves a group from the local storage specified by id
func (r *Repository) RetrieveGroup(id string) (*listing.Group, error) {
	file, err := os.OpenFile(r.path+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, listing.ErrGroupNotFound
	}

	var grp Group
	if err := retrieveJSON(file, &grp); err != nil {
		return &listing.Group{}, err
	}

	return &listing.Group{
		ID:      grp.ID,
		Configs: grp.Configs,
	}, nil

}

// StoreConfig stores a config in the local storage
func (r *Repository) StoreConfig(c adding.Config) error {
	grp, err := r.RetrieveGroup(c.Group)
	if err != nil {
		return err
	}

	basePath := r.path + "/" + c.Group + "/"

	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(basePath+c.ID+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// TODO: There is a chance that the config will get created and storing the new group with config added will fail. Fix fix.
	// TODO: Sort array?
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

	addGrp := adding.Group{
		ID:      grp.ID,
		Configs: grp.Configs,
	}

	if err := r.storeGroup(addGrp); err != nil {
		log.Println("Failed persisting Group when new config was added. Config not added.")
		return err
	}

	conf := Config{
		ID:           c.ID,
		Name:         c.Name,
		LastModified: c.LastModified,
		Version:      c.Version,
		Group:        c.Group,
		Properties:   c.Properties,
	}

	return storeJSON(file, conf)
}

// RetrieveConfig retrieves a config from the local storage spesified by groupID and id of the config
func (r *Repository) RetrieveConfig(groupID string, id string) (*listing.Config, error) {
	if _, err := r.RetrieveGroup(groupID); err != nil {
		return &listing.Config{}, err
	}

	file, err := os.OpenFile(r.path+"/"+groupID+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, listing.ErrConfigNotFound
	}

	var c Config
	err = retrieveJSON(file, &c)

	return &listing.Config{
		ID:           c.ID,
		Name:         c.Name,
		LastModified: c.LastModified,
		Version:      c.Version,
		Group:        c.Group,
		Properties:   c.Properties,
	}, err
}

func storeJSON(file *os.File, v interface{}) error {
	return json.NewEncoder(file).Encode(v)
}

func retrieveJSON(file *os.File, v interface{}) error {
	return json.NewDecoder(file).Decode(v)
}
