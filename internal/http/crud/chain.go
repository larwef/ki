package crud

import "net/http"

// Used to chain http handler functions.
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

// Returns a handler comprised of all the handler functions added to the chain.
func (hc *handlerChain) buildChain() *handlerChain {
	for _, handlerFunc := range hc.handlers {
		hc.chained = handlerFunc(hc.chained)
	}

	return hc
}

func emptyHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {})
}
