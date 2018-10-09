package crud

import (
	"bytes"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// requestLoggerWrapper wraps a http.Request. All writes to the http.Request gets written to a buffer which can be used to log the
// content written to it by calling the log function.
type requestLoggerWrapper struct {
	breadcrumb string
	req        *http.Request
	reader     io.ReadCloser
	buffer     io.ReadWriter
}

// newRequestLogger returns a new RequestLogger
func newRequestLogger(req *http.Request, breadcrumb string) *requestLoggerWrapper {
	buffer := new(bytes.Buffer)
	return &requestLoggerWrapper{
		breadcrumb: breadcrumb,
		req:        req,
		reader:     ioutil.NopCloser(io.TeeReader(req.Body, buffer)),
		buffer:     buffer,
	}
}

// log prints the content of the buffer
func (r *requestLoggerWrapper) log() {
	reqBytes, err := ioutil.ReadAll(r.buffer)
	if err != nil {
		log.Printf("Error logging request: %v\n", err)
	}

	log.Printf("Inbound message:\nBreadcrumb: %s\nHost: %s\nRemoteAddr: %s\nMethod: %s\nProto: %s\nPath: %s\nPayload: %s",
		r.breadcrumb, r.req.Host, r.req.RemoteAddr, r.req.Method, r.req.Proto, r.req.URL.Path, string(reqBytes))
}

func (r *requestLoggerWrapper) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *requestLoggerWrapper) Close() error {
	return r.reader.Close()
}

// responseLoggerWrapper wraps a http.ResponseWriter. All writes to the http.ResponseWriter gets written to a buffer which can be
// used to log the content written to it by calling the log function.
type responseLoggerWrapper struct {
	breadcrumb string
	status     int
	resWriter  http.ResponseWriter
	writer     io.Writer
	buffer     io.ReadWriter
}

// newResponseLoggerWrapper returns a new responseLoggerWrapper
func newResponseLoggerWrapper(responseWriter http.ResponseWriter, breadcrumb string) *responseLoggerWrapper {
	buffer := new(bytes.Buffer)
	return &responseLoggerWrapper{
		breadcrumb: breadcrumb,
		writer:     io.MultiWriter(responseWriter, buffer),
		resWriter:  responseWriter,
		buffer:     buffer,
	}
}

// log prints the content of the buffer
func (r *responseLoggerWrapper) log() {
	resBytes, err := ioutil.ReadAll(r.buffer)
	if err != nil {
		log.Printf("Error logging request: %v\n", err)
	}

	log.Printf("Outbound Response:\nBreadcrumb: %s\nResponse-Code: %d\nHeaders: %v\nPayload: %s",
		r.breadcrumb, r.status, r.resWriter.Header(), string(resBytes))
}

func (r *responseLoggerWrapper) Header() http.Header {
	return r.resWriter.Header()
}

func (r *responseLoggerWrapper) WriteHeader(status int) {
	r.status = status
	r.resWriter.WriteHeader(status)
}

func (r *responseLoggerWrapper) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.status = 200
	}

	return r.writer.Write(p)
}

// TODO: Still not happy with this. What if something is logged inside another chain link? The innOut log wil probably come at the end
// TODO: Coudl just go for a simple
func inOutLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		breadCrumb := uuid.New().String()
		requestLogger := newRequestLogger(req, breadCrumb)
		req.Body = requestLogger
		resWriter := newResponseLoggerWrapper(res, breadCrumb)

		h.ServeHTTP(resWriter, req)

		requestLogger.log()
		resWriter.log()
	})
}
