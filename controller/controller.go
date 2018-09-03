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

type BaseHttpHandler struct {
	configHandler *configHandler
}

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
