package controller

import (
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
	"log"
	"net/http"
	"time"
)

// TODO: Should check path -> method, not method -> path
func (handler *Handler) handleConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPut:
			newHandlerChain(h).
				add(configPathValidator).
				add(setCommonHeaders).
				add(handler.handlePut).
				ServeHTTP(res, req)
		case http.MethodGet:
			newHandlerChain(h).
				add(configPathValidator).
				add(setCommonHeaders).
				add(handler.handleGet).
				ServeHTTP(res, req)
		default:
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

// TODO: The routing could probably be more elegant
func (handler *Handler) handlePut(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		chain := newHandlerChain(h)

		_, id, _ := getPathVariables(req.URL.Path)

		if id == "" {
			chain.
				add(handler.storeGroup).
				add(handler.retrieveGroup)
		} else {
			chain.
				add(handler.storeConfig).
				add(handler.retrieveConfig)
		}

		chain.ServeHTTP(res, req)
	})
}

// TODO: The routing could probably be more elegant
func (handler *Handler) handleGet(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		chain := newHandlerChain(h)

		_, id, _ := getPathVariables(req.URL.Path)

		if id == "" {
			chain.
				add(handler.retrieveGroup)
		} else {
			chain.
				add(handler.retrieveConfig)
		}

		chain.ServeHTTP(res, req)
	})
}

func (handler *Handler) storeGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var grp adding.Group

		if err := json.NewDecoder(req.Body).Decode(&grp); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Unable to unmarshal request object", http.StatusInternalServerError)
			return
		}

		defer req.Body.Close()

		grp.ID, _, _ = getPathVariables(req.URL.Path)

		if err := handler.adding.AddGroup(grp); err != nil {
			log.Printf("Error: %v", err)
			if err == adding.ErrGroupConflict {
				http.Error(res, err.Error(), http.StatusConflict)
				return
			}
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (handler *Handler) retrieveGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		grpID, _, _ := getPathVariables(req.URL.Path)

		var conf *listing.Group
		var err error
		if conf, err = handler.listing.GetGroup(grpID); err != nil {
			log.Printf("Error: %v", err)
			if err == listing.ErrGroupNotFound {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
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

func (handler *Handler) storeConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var conf adding.Config

		if err := json.NewDecoder(req.Body).Decode(&conf); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Unable to unmarshal request object", http.StatusInternalServerError)
			return
		}

		defer req.Body.Close()

		conf.Group, conf.ID, _ = getPathVariables(req.URL.Path)
		conf.LastModified = time.Now()

		if err := handler.adding.AddConfig(conf); err != nil {
			log.Printf("Error: %v", err)
			if err == listing.ErrGroupNotFound {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (handler *Handler) retrieveConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		grp, id, _ := getPathVariables(req.URL.Path)

		var conf *listing.Config
		var err error
		if conf, err = handler.listing.GetConfig(grp, id); err != nil {
			log.Printf("Error: %v", err)
			if err == listing.ErrGroupNotFound || err == listing.ErrConfigNotFound {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
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
