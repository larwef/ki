package crud

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

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
