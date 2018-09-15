package repository

import (
	"encoding/json"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/group"
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

// StorConfig stores a config in the local storage
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

// RetrieveConfig retrieves a config from the local storage
func (l *Local) RetrieveGroup(id string) (*group.Group, error) {

	file, err := os.OpenFile(l.path+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var g group.Group
	err = retrieveJSON(file, &g)
	return &g, err

}

// StorConfig stores a config in the local storage
func (l *Local) StoreConfig(c config.Config) error {
	basePath := l.path + "/" + c.Group + "/"

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(basePath+c.ID+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return storeJSON(file, c)
}

// RetrieveConfig retrieves a config from the local storage
func (l *Local) RetrieveConfig(group string, id string) (*config.Config, error) {

	file, err := os.OpenFile(l.path+"/"+group+"/"+id+".json", os.O_RDONLY, 0644)
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
