package persistence

import (
	"encoding/json"
	"github.com/larwef/ki/config"
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

// Store stores a config in the local storage
func (l *Local) Store(c config.Config) error {
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

func storeJSON(file *os.File, c config.Config) error {
	return json.NewEncoder(file).Encode(c)
}

// Retrieve retrieves a config from the local storage
func (l *Local) Retrieve(group string, id string) (*config.Config, error) {

	file, err := os.OpenFile(l.path+"/"+group+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var c config.Config
	err = retrieveJSON(file, &c)
	return &c, err

}

func retrieveJSON(file *os.File, c *config.Config) error {
	return json.NewDecoder(file).Decode(&c)
}
