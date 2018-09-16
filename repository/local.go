package repository

import (
	"encoding/json"
	"log"
	"os"
)

// Local representa a local storge object
type Local struct {
	path string
}

// NewLocal returns a new Local storage object
func NewLocal(path string) *Local {
	return &Local{path: path}
}

// StoreGroup stores a config in the local storage. Will not overwrite existing Group object.
func (l *Local) StoreGroup(g Group) error {
	fullPath := l.path + "/" + g.ID + ".json"
	if _, err := os.Stat(fullPath); err == nil {
		return ErrConflict
	}

	return l.storeGroup(g)
}

// This store function will overwrite the group object
func (l *Local) storeGroup(g Group) error {
	basePath := l.path + "/"
	fullPath := basePath + g.ID + ".json"

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return ErrInternal
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return ErrInternal
	}

	return storeJSON(file, g)
}

// RetrieveGroup retrieves a group from the local storage specified by id
func (l *Local) RetrieveGroup(id string) (*Group, error) {

	file, err := os.OpenFile(l.path+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	var g Group
	err = retrieveJSON(file, &g)
	return &g, err

}

// StoreConfig stores a config in the local storage
func (l *Local) StoreConfig(c Config) error {
	grp, err := l.RetrieveGroup(c.Group)
	if err != nil {
		return err
	}

	basePath := l.path + "/" + c.Group + "/"

	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return ErrInternal
	}

	file, err := os.OpenFile(basePath+c.ID+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return ErrInternal
	}

	// TODO: There is a chance that the config will get created and storing the new group with config added will fail. Fix fix.
	// TODO: Should not append if configID already exists
	// TODO: Sort array?
	grp.Configs = append(grp.Configs, c.ID)
	if err := l.storeGroup(*grp); err != nil {
		log.Println("Failed persisting Group when new config was added. The config will most likely exist but not be added to Group config array.")
		return err
	}

	return storeJSON(file, c)
}

// RetrieveConfig retrieves a config from the local storage spesified by groupID and id of the config
func (l *Local) RetrieveConfig(groupID string, id string) (*Config, error) {

	if _, err := l.RetrieveGroup(groupID); err != nil {
		return &Config{}, err
	}

	file, err := os.OpenFile(l.path+"/"+groupID+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, ErrConfigNotFound
	}

	var c Config
	err = retrieveJSON(file, &c)

	return &c, err
}

func storeJSON(file *os.File, v interface{}) error {
	return json.NewEncoder(file).Encode(v)
}

func retrieveJSON(file *os.File, v interface{}) error {
	return json.NewDecoder(file).Decode(v)
}
