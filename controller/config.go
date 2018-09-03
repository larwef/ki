package controller

import (
	"encoding/json"
	"errors"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/config/persistence"
	"log"
	"net/http"
)

type configHandler struct {
	persistence persistence.Persistence
}

func (c *configHandler) handleConfig(res http.ResponseWriter, req *http.Request) {
	log.Printf("Config invoked")

	switch req.Method {
	case http.MethodPut:
		c.handlePut(res, req)
	case http.MethodGet:
		c.handleGet(res, req)
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (c *configHandler) handlePut(res http.ResponseWriter, req *http.Request) {
	log.Println("PUT invoked")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	group, id, err := getPathVariables(req.URL.Path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var conf config.Config

	err = json.NewDecoder(req.Body).Decode(&conf)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Unable to unmarshal request object", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	conf.Id = id
	conf.Group = group

	err = c.persistence.Store(conf)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Error persisting  request object", http.StatusInternalServerError)
		return
	}

	newConf, err := c.persistence.Retrieve(group, id)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(res).Encode(newConf)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Error marshalling response", http.StatusInternalServerError)
		return
	}
}

func (c *configHandler) handleGet(res http.ResponseWriter, req *http.Request) {
	log.Println("GET invoked")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	group, id, err := getPathVariables(req.URL.Path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	conf, err := c.persistence.Retrieve(group, id)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(res).Encode(conf)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Error marshalling response", http.StatusInternalServerError)
		return
	}
}

func getPathVariables(url string) (string, string, error) {
	var id, group string
	group, url = shiftPath(url)
	id, url = shiftPath(url)

	if group == "" || id == "" || url != "/" {
		return "", "", errors.New("invalid path")
	}

	return group, id, nil
}
