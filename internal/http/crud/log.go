package crud

import (
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type requestLogger struct {
	breadCrumb string
	req        *http.Request
	io.ReadCloser
}

func (r *requestLogger) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	log.Printf("Inbound message:\nBreadcrumb: %s\nHost: %s\nRemoteAddr: %s\nMethod: %s\nProto: %s\nPath: %s\nPayload: %s",
		r.breadCrumb, r.req.Host, r.req.RemoteAddr, r.req.Method, r.req.Proto, r.req.URL.Path, string(p))

	return n, err
}

func (r *requestLogger) Close() error {
	return r.ReadCloser.Close()
}

type responseLogger struct {
	breadcrumb string
	http.ResponseWriter
	status int
}

func (r *responseLogger) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseLogger) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.status = 200
	}

	log.Printf("Outbound Response:\nBreadcrumb: %s\nResponse-Code: %d\nHeaders: %v\nPayload: %s",
		r.breadcrumb, r.status, r.ResponseWriter.Header(), string(p))

	return r.ResponseWriter.Write(p)
}

func inOutLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		breadCrumb := uuid.New().String()

		// Workaround for logging GET request
		if req.Method == http.MethodGet {
			log.Printf("Inbound Request:\nBreadcrumb: %s\nHost: %s\nRemoteAddr: %s\nMethod: %s\nProto: %s\nPath: %s\n",
				breadCrumb, req.Host, req.RemoteAddr, req.Method, req.Proto, req.URL.Path)
		} else {
			req.Body = &requestLogger{
				breadCrumb: breadCrumb,
				req:        req,
				ReadCloser: req.Body,
			}
		}

		responseWriter := responseLogger{
			breadcrumb:     breadCrumb,
			ResponseWriter: res,
		}

		h.ServeHTTP(&responseWriter, req)
	})
}
