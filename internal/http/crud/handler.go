package crud

import (
	"encoding/json"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	healthPath = "health"
	configPath = "config"

	contentType = "application/json; charset=utf-8"
)

// Handler handles is the entry point for requests and handles routing and processing
type Handler struct {
	adding  adding.Service
	listing listing.Service
}

// NewHandler returns a new Handler
func NewHandler(adding adding.Service, listing listing.Service) *Handler {
	return &Handler{
		adding:  adding,
		listing: listing,
	}
}

func (handler *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	newHandlerChain(emptyHandler()).
		add(inOutLog).
		add(setCommonHeaders).
		add(handler.route).
		ServeHTTP(res, req)
}

func (handler *Handler) route(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		service, _, _, _ := getPathVariables(req.URL.Path)

		switch service {
		case healthPath:
			newHandlerChain(emptyHandler()).
				add(handler.handleHealth).
				ServeHTTP(res, req)

		case configPath:
			newHandlerChain(emptyHandler()).
				add(handler.handleConfig).
				ServeHTTP(res, req)
		default:
			log.Printf("Invalid path %q called\n", req.URL.Path)
			http.Error(res, "Not Found", http.StatusNotFound)
		}
	})
}

func (handler *Handler) handleHealth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		healthCheck := &health{Status: "OK"}
		if err := json.NewEncoder(res).Encode(healthCheck); err != nil {
			log.Printf("Health check failed. Error: %v", err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (handler *Handler) handleConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		chain := newHandlerChain(h)
		_, grpID, confID, remainder := getPathVariables(req.URL.Path)

		if remainder != "/" {
			log.Printf("Invalid path %q called", req.URL.Path)
			http.Error(res, "Invalid Path", http.StatusBadRequest)
		} else if grpID != "" && confID == "" && remainder == "/" {
			chain.add(handler.handleGroupAction)
		} else if grpID != "" && confID != "" && remainder == "/" {
			chain.add(handler.handleConfigAction)
		} else {
			log.Printf("Unexpected state when processing path %q", req.URL.Path)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}

		chain.ServeHTTP(res, req)
	})
}

func (handler *Handler) handleGroupAction(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPut:
			newHandlerChain(h).
				add(handler.storeGroup).
				add(handler.retrieveGroup).
				ServeHTTP(res, req)
			break
		case http.MethodGet:
			newHandlerChain(h).
				add(handler.retrieveGroup).
				ServeHTTP(res, req)
			break
		default:
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (handler *Handler) handleConfigAction(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPut:
			newHandlerChain(h).
				add(handler.storeConfig).
				add(handler.retrieveConfig).
				ServeHTTP(res, req)
			break
		case http.MethodGet:
			newHandlerChain(h).
				add(handler.retrieveConfig).
				ServeHTTP(res, req)
			break
		default:
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (handler *Handler) storeGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var grp adding.Group

		if err := json.NewDecoder(req.Body).Decode(&grp); err != nil {
			log.Printf("Error: %v", err)
			http.Error(res, "Unable to unmarshal request object", http.StatusBadRequest)
			return
		}

		defer req.Body.Close()

		_, grp.ID, _, _ = getPathVariables(req.URL.Path)

		if err := handler.adding.AddGroup(grp); err != nil {
			log.Printf("Error: %v", err)
			writeServiceError(res, err)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (handler *Handler) retrieveGroup(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, grpID, _, _ := getPathVariables(req.URL.Path)

		var conf *listing.Group
		var err error
		if conf, err = handler.listing.GetGroup(grpID); err != nil {
			writeServiceError(res, err)
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

		_, conf.Group, conf.ID, _ = getPathVariables(req.URL.Path)
		conf.LastModified = time.Now()

		if err := handler.adding.AddConfig(conf); err != nil {
			writeServiceError(res, err)
			return
		}

		h.ServeHTTP(res, req)
	})
}

func (handler *Handler) retrieveConfig(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, grp, id, _ := getPathVariables(req.URL.Path)

		var conf *listing.Config
		var err error
		if conf, err = handler.listing.GetConfig(grp, id); err != nil {
			writeServiceError(res, err)
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

func writeServiceError(res http.ResponseWriter, err error) {
	log.Printf("Error: %v", err)
	switch err {
	case adding.ErrGroupConflict:
		fallthrough
	case listing.ErrGroupNotFound:
		fallthrough
	case listing.ErrConfigNotFound:
		http.Error(res, err.Error(), http.StatusNotFound)
	default:
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	}
}

func setCommonHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", contentType)
		h.ServeHTTP(res, req)
	})
}

func getPathVariables(url string) (string, string, string, string) {
	var serivce, grp, id string
	serivce, url = shiftPath(url)
	grp, url = shiftPath(url)
	id, url = shiftPath(url)

	return serivce, grp, id, url
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
