package controller

import (
	"encoding/json"
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/group"
	"github.com/larwef/ki/repository"
	"log"
	"net/http"
	"time"
)

type configHandler struct {
	repository repository.Repository
}

// TODO: Should check path -> method, not method -> path
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

// TODO: The routing could probably be more elegant
func (c *configHandler) handlePut(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		chain := newHandlerChain(h)

		_, id, _ := getPathVariables(req.URL.Path)

		if id == "" {
			chain.
				add(c.storeGroup).
				add(c.retrieveGroup)
		} else {
			chain.
				add(c.storeConfig).
				add(c.retrieveConfig)
		}

		chain.ServeHTTP(res, req)
	})
}

// TODO: The routing could probably be more elegant
func (c *configHandler) handleGet(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		chain := newHandlerChain(h)

		_, id, _ := getPathVariables(req.URL.Path)

		if id == "" {
			chain.
				add(c.retrieveGroup)
		} else {
			chain.
				add(c.retrieveConfig)
		}

		chain.ServeHTTP(res, req)
	})
}

func (c *configHandler) storeGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var grp group.Group

		if err := json.NewDecoder(req.Body).Decode(&grp); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Unable to unmarshal request object", http.StatusInternalServerError)
			return
		}

		defer req.Body.Close()

		grp.ID, _, _ = getPathVariables(req.URL.Path)

		if err := c.repository.StoreGroup(grp); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Error persisting request object", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (c *configHandler) retrieveGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		grp, _, _ := getPathVariables(req.URL.Path)

		var conf *group.Group
		var err error
		if conf, err = c.repository.RetrieveGroup(grp); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, err.Error(), http.StatusNotFound)
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

		if err := c.repository.StoreConfig(conf); err != nil {
			log.Printf("Error: %v", err)
			// TODO: Perfect example for why a custom error type is needed. Here there are several errors which may have different http codes
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (c *configHandler) retrieveConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		grp, id, _ := getPathVariables(req.URL.Path)

		var conf *config.Config
		var err error
		if conf, err = c.repository.RetrieveConfig(grp, id); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, err.Error(), http.StatusNotFound)
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
		if invalidConfigPath(req.URL.Path) {
			http.Error(res, "Invalid Path", http.StatusBadRequest)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func invalidConfigPath(url string) bool {
	grp, _, remainder := getPathVariables(url)
	if grp == "" || remainder != "/" {
		return true
	}

	return false
}

// TODO: Generalize?
func getPathVariables(url string) (string, string, string) {
	var grp, id string
	grp, url = shiftPath(url)
	id, url = shiftPath(url)

	return grp, id, url
}
