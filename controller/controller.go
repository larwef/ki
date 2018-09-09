package controller

import (
	"github.com/larwef/ki/config/persistence"
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	statusPath = "status"
	configPath = "config"
)

type baseHttpHandler struct {
	configHandler *configHandler
}

func NewBaseHttpHandler(persistence persistence.Persistence) *baseHttpHandler {
	return &baseHttpHandler{
		configHandler: &configHandler{persistence: persistence},
	}
}

func (b *baseHttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = shiftPath(req.URL.Path)

	switch head {
	case configPath:
		newHandlerChain(emptyHandler()).
			add(b.configHandler.handleConfig).
			ServeHTTP(res, req)
	default:
		log.Printf("Invalid path %s called", head)
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}

type handlerChain struct {
	handlers []func(handler http.Handler) http.Handler
	chained  http.Handler
}

func newHandlerChain(h http.Handler) *handlerChain {
	return &handlerChain{chained: h}
}

func (hc *handlerChain) add(h func(http.Handler) http.Handler) *handlerChain {
	// Prepend handler function
	hc.handlers = append([]func(http.Handler) http.Handler{h}, hc.handlers...)

	return hc
}

func (hc *handlerChain) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	hc.buildChain().chained.ServeHTTP(res, req)
}

func (hc *handlerChain) buildChain() *handlerChain {
	for _, handlerFunc := range hc.handlers {
		hc.chained = handlerFunc(hc.chained)
	}

	return hc
}

func emptyHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {})
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
