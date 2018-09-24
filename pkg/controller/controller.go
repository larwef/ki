package controller

import (
	"bytes"
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/listing"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	statusPath = "status"
	configPath = "config"
)

// Handler handles is the entry point for requests and handles initial routing
type Handler struct {
	adding  adding.Service
	listing listing.Service
}

// NewBaseHTTPHandler returns a new Handler
func NewBaseHTTPHandler(adding adding.Service, listing listing.Service) *Handler {
	return &Handler{
		adding:  adding,
		listing: listing,
	}
}

func (handler *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = shiftPath(req.URL.Path)

	switch head {
	case configPath:
		newHandlerChain(emptyHandler()).
			add(inOutLog).
			add(handler.handleConfig).
			//add(handler.groupHandler.handleGroup).
			ServeHTTP(res, req)
	default:
		log.Printf("Invalid path %q called", head)
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}

// TODO: Check if this could be done more elegantly
type requestLogger struct{}

func (il *requestLogger) logRequest(req *http.Request) {
	var bodyString string

	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("Error reading incomming request")
		}
		bodyString = string(b)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}

	log.Printf("Inbound message:\nHost: %s\nRemoteAddr: %s\nMethod: %s\nProto: %s\nPayload: %s", req.Host, req.RemoteAddr, req.Method, req.Proto, bodyString)
}

type responseLogger struct {
	http.ResponseWriter
	status int
}

func (w *responseLogger) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseLogger) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}

	log.Printf("Outbound message:\nResponse-Code: %d\nHeaders: %v\nPayload: %s", w.status, w.ResponseWriter.Header(), string(b))

	return w.ResponseWriter.Write(b)
}

func inOutLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		rl := requestLogger{}
		rl.logRequest(req)
		responseWriter := responseLogger{ResponseWriter: res}
		h.ServeHTTP(&responseWriter, req)
	})
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