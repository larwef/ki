package persistence

import (
	"encoding/json"
	"github.com/larwef/ki/config"
	"os"
)

type local struct {
	path string
}

func NewLocal(path string) *local {
	return &local{path: path}
}

func (l *local) Store(c config.Config) error {
	basePath := l.path + "/" + c.Group + "/"

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(basePath+c.Id+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return storeJson(file, c)
}

func storeJson(file *os.File, c config.Config) error {
	return json.NewEncoder(file).Encode(c)
}

func (l *local) Retrieve(group string, id string) (*config.Config, error) {

	file, err := os.OpenFile(l.path+"/"+group+"/"+id+".json", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var c config.Config
	err = retrieveJson(file, &c)
	return &c, err

}

func retrieveJson(file *os.File, c *config.Config) error {
	return json.NewDecoder(file).Decode(&c)
}
