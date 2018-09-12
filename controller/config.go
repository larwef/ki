package controller

import (
	"encoding/json"
	"github.com/larwef/ki/config"
	"log"
	"net/http"
	"time"
)

type configHandler struct {
	configRepo config.Repository
}

func (c *configHandler) handleConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPut:
			newHandlerChain(h).
				add(configPathValidator).
				add(setCommonHeaders).
				add(c.handlePut).
				ServeHTTP(res, req)
		case http.MethodGet:
			newHandlerChain(h).
				add(configPathValidator).
				add(setCommonHeaders).
				add(c.handleGet).
				ServeHTTP(res, req)
		default:
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (c *configHandler) handlePut(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newHandlerChain(h).
			add(c.storeConfig).
			add(c.retrieveConfig).
			ServeHTTP(res, req)
	})
}

func (c *configHandler) handleGet(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newHandlerChain(h).
			add(c.retrieveConfig).
			ServeHTTP(res, req)
	})
}

func (c *configHandler) storeConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var conf config.Config

		if err := json.NewDecoder(req.Body).Decode(&conf); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Unable to unmarshal request object", http.StatusInternalServerError)
			return
		}

		defer req.Body.Close()

		conf.Group, conf.ID, _ = getPathVariables(req.URL.Path)
		conf.LastModified = time.Now()

		if err := c.configRepo.Store(conf); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Error persisting request object", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (c *configHandler) retrieveConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		group, id, _ := getPathVariables(req.URL.Path)

		var conf *config.Config
		var err error
		if conf, err = c.configRepo.Retrieve(group, id); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		}

		if err = json.NewEncoder(res).Encode(conf); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Error marshalling response", http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(res, req)
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
