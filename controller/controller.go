package controller

import (
	"encoding/json"
	"github.com/larwef/ki/config/persistence"
	"github.com/larwef/ki/config"
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	statusPath = "status"
	configPath = "config"
)

type (
	BaseHttpHandler struct {
		configHandler *configHandler
	}

	configHandler struct {
		persistence persistence.Persistence
	}
)

func NewBaseHttpHandler(persistence persistence.Persistence) *BaseHttpHandler {
	return &BaseHttpHandler{
		configHandler: &configHandler{persistence: persistence},
	}
}

func (b *BaseHttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = shiftPath(req.URL.Path)

	switch head {
	case configPath:
		http.HandlerFunc(b.configHandler.handleConfig).ServeHTTP(res, req)
	default:
		log.Printf("Invalid path %s called", head)
		http.Error(res, "Not Found", http.StatusNotFound)
	}
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

//TODO: Refractor duplicate code
func (c *configHandler) handlePut(res http.ResponseWriter, req *http.Request) {
	log.Println("PUT invoked")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	var id, group string

	group, req.URL.Path = shiftPath(req.URL.Path)
	id, req.URL.Path = shiftPath(req.URL.Path)

	if group == "" || id == "" || req.URL.Path != "/" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	var conf config.Config

	err := json.NewDecoder(req.Body).Decode(&conf)
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

	storedConf, err := c.persistence.Retrieve(id, group)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(res).Encode(storedConf)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(res, "Error marshalling response", http.StatusInternalServerError)
		return
	}
}

func (c *configHandler) handleGet(res http.ResponseWriter, req *http.Request) {
	log.Println("GET invoked")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	var id, group string

	group, req.URL.Path = shiftPath(req.URL.Path)
	id, req.URL.Path = shiftPath(req.URL.Path)

	if group == "" || id == "" || req.URL.Path != "/" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	conf, err := c.persistence.Retrieve(id, group)
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

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
