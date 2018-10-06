package crud

import (
	"bytes"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// RequestLoggerWrapper wraps a http.Request. All writes to the http.Request gets written to a buffer which can be used to Log the
// content written to it by calling the Log function.
type RequestLoggerWrapper struct {
	breadcrumb string
	req        *http.Request
	reader     io.ReadCloser
	buffer     io.ReadWriter
}

// NewRequestLogger returns a new RequestLogger
func NewRequestLogger(req *http.Request, breadcrumb string) *RequestLoggerWrapper {
	buffer := new(bytes.Buffer)
	return &RequestLoggerWrapper{
		breadcrumb: breadcrumb,
		req:        req,
		reader:     ioutil.NopCloser(io.TeeReader(req.Body, buffer)),
		buffer:     buffer,
	}
}

// Log prints the content of the buffer
func (r *RequestLoggerWrapper) Log() {
	reqBytes, err := ioutil.ReadAll(r.buffer)
	if err != nil {
		log.Printf("Error logging request: %v\n", err)
	}

	log.Printf("Inbound message:\nBreadcrumb: %s\nHost: %s\nRemoteAddr: %s\nMethod: %s\nProto: %s\nPath: %s\nPayload: %s",
		r.breadcrumb, r.req.Host, r.req.RemoteAddr, r.req.Method, r.req.Proto, r.req.URL.Path, string(reqBytes))
}

func (r *RequestLoggerWrapper) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *RequestLoggerWrapper) Close() error {
	return r.reader.Close()
}

// ResponseLoggerWrapper wraps a http.ResponseWriter. All writes to the http.ResponseWriter gets written to a buffer which can be
// used to Log the content written to it by calling the Log function.
type ResponseLoggerWrapper struct {
	breadcrumb string
	status     int
	resWriter  http.ResponseWriter
	writer     io.Writer
	buffer     io.ReadWriter
}

// NewResponseLoggerWrapper returns a new ResponseLoggerWrapper
func NewResponseLoggerWrapper(responseWriter http.ResponseWriter, breadcrumb string) *ResponseLoggerWrapper {
	buffer := new(bytes.Buffer)
	return &ResponseLoggerWrapper{
		breadcrumb: breadcrumb,
		writer:     io.MultiWriter(responseWriter, buffer),
		resWriter:  responseWriter,
		buffer:     buffer,
	}
}

// Log prints the content of the buffer
func (r *ResponseLoggerWrapper) Log() {
	resBytes, err := ioutil.ReadAll(r.buffer)
	if err != nil {
		log.Printf("Error logging request: %v\n", err)
	}

	log.Printf("Outbound Response:\nBreadcrumb: %s\nResponse-Code: %d\nHeaders: %v\nPayload: %s",
		r.breadcrumb, r.status, r.resWriter.Header(), string(resBytes))
}

func (r *ResponseLoggerWrapper) Header() http.Header {
	return r.resWriter.Header()
}

func (r *ResponseLoggerWrapper) WriteHeader(status int) {
	r.status = status
	r.resWriter.WriteHeader(status)
}

func (r *ResponseLoggerWrapper) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.status = 200
	}

	return r.writer.Write(p)
}

// TODO: Still not happy with this. What if something is logged inside another chain link? The innOut Log wil probably come at the end
// TODO: Coudl just go for a simple
func inOutLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		breadCrumb := uuid.New().String()
		requestLogger := NewRequestLogger(req, breadCrumb)
		req.Body = requestLogger
		resWriter := NewResponseLoggerWrapper(res, breadCrumb)

		h.ServeHTTP(resWriter, req)

		requestLogger.Log()
		resWriter.Log()
	})
}
