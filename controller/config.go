package controller

import (
	"encoding/json"
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
		c.handlePut().ServeHTTP(res, req)
	case http.MethodGet:
		c.handleGet().ServeHTTP(res, req)
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (c *configHandler) handlePut() http.Handler {
	log.Println("PUT invoked")
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		configPathValidator(setCommonHeaders(c.putConfig(c.getConfig(emptyHandler())))).ServeHTTP(res, req)
	})
}

func (c *configHandler) handleGet() http.Handler {
	log.Println("GET invoked")
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		configPathValidator(setCommonHeaders(c.getConfig(emptyHandler()))).ServeHTTP(res, req)
	})
}

func (c *configHandler) putConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var conf config.Config

		if err := json.NewDecoder(req.Body).Decode(&conf); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Unable to unmarshal request object", http.StatusInternalServerError)
			return
		}

		defer req.Body.Close()

		conf.Group, conf.Id, _ = getPathVariables(req.URL.Path)

		if err := c.persistence.Store(conf); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Error persisting request object", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (c *configHandler) getConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		group, id, _ := getPathVariables(req.URL.Path)

		var conf *config.Config
		var err error
		if conf, err = c.persistence.Retrieve(group, id); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		}

		if err = json.NewEncoder(res).Encode(conf); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Error marshalling response", http.StatusInternalServerError)
			return
		}
	})
}

func setCommonHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")

		h.ServeHTTP(res, req)
	})
}

func configPathValidator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		group, id, remainder := getPathVariables(req.URL.Path)

		if group == "" || id == "" || remainder != "/" {
			http.Error(res, "Invalid path", http.StatusBadRequest)
			return
		}

		h.ServeHTTP(res, req)
	})
}

// TODO: Generalize?
func getPathVariables(url string) (string, string, string) {
	var group, id string
	group, url = shiftPath(url)
	id, url = shiftPath(url)

	return group, id, url
}
