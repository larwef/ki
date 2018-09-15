package controller

import (
	"encoding/json"
	"github.com/larwef/ki/group"
	"github.com/larwef/ki/repository"
	"log"
	"net/http"
)

type groupHandler struct {
	repository repository.Repository
}

func (g *groupHandler) handleGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPut:
			newHandlerChain(h).
				//add(configPathValidator).
				add(setCommonHeaders).
				add(g.handlePut).
				ServeHTTP(res, req)
		case http.MethodGet:
			newHandlerChain(h).
				//add(configPathValidator).
				add(setCommonHeaders).
				add(g.handleGet).
				ServeHTTP(res, req)
		default:
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (g *groupHandler) handlePut(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newHandlerChain(h).
			add(g.storeGroup).
			add(g.retrieveGroup).
			ServeHTTP(res, req)
	})
}

func (g *groupHandler) handleGet(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newHandlerChain(h).
			add(g.retrieveGroup).
			ServeHTTP(res, req)
	})
}

func (g *groupHandler) storeGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var grp group.Group

		if err := json.NewDecoder(req.Body).Decode(&grp); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Unable to unmarshal request object", http.StatusInternalServerError)
			return
		}

		defer req.Body.Close()

		grp.ID, _, _ = getPathVariables(req.URL.Path)

		if err := g.repository.StoreGroup(grp); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Error persisting request object", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (g *groupHandler) retrieveGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		grp, _, _ := getPathVariables(req.URL.Path)

		var conf *group.Group
		var err error
		if conf, err = g.repository.RetrieveGroup(grp); err != nil {
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
