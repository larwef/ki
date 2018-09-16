package repository

import (
	"encoding/json"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/group"
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

// StoreGroup stores a config in the local storage
func (l *Local) StoreGroup(g group.Group) error {
	basePath := l.path + "/"

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(basePath+g.ID+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return storeJSON(file, g)
}

// RetrieveGroup retrieves a group from the local storage specified by id
func (l *Local) RetrieveGroup(id string) (*group.Group, error) {

	file, err := os.OpenFile(l.path+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var g group.Group
	err = retrieveJSON(file, &g)
	return &g, err

}

// StoreConfig stores a config in the local storage
func (l *Local) StoreConfig(c config.Config) error {
	grp, err := l.RetrieveGroup(c.Group)
	if err != nil {
		return err
	}

	basePath := l.path + "/" + c.Group + "/"

	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(basePath+c.ID+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// TODO: There is a chance that the config will get created and storing the new group with config added will fail. Fix fix.
	grp.Configs = append(grp.Configs, c.ID)
	if err := l.StoreGroup(*grp); err != nil {
		log.Println("Failed persisting Group when new config was added. The config will most likely exist but not be added to Group config array.")
		return err
	}

	return storeJSON(file, c)
}

// RetrieveConfig retrieves a config from the local storage spesified by groupID and id of the config
func (l *Local) RetrieveConfig(groupID string, id string) (*config.Config, error) {

	if _, err := l.RetrieveGroup(groupID); err != nil {
		return &config.Config{}, err
	}

	file, err := os.OpenFile(l.path+"/"+groupID+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var c config.Config
	err = retrieveJSON(file, &c)

	return &c, err
}

func storeJSON(file *os.File, v interface{}) error {
	return json.NewEncoder(file).Encode(v)
}

func retrieveJSON(file *os.File, v interface{}) error {
	return json.NewDecoder(file).Decode(v)
}
