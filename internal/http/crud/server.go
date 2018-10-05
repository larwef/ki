package crud

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server represents a crud server object
type Server struct {
	Server *http.Server
}

// Serve starts listening on the crud server. Sends a signal on error.
// TODO: Improvement: When Shutdown is called ListenAndServe will return, hopefully without error. Should it return with an error
// TODO: the program will panic (since this will call close on a closed channel). So this could be executed cleaner.
func (s *Server) Serve(signal chan bool) {
	log.Printf("Starting crud server on %s\n", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		log.Printf("crud server error: %v\n", err)
		signal <- true
	}
}

// GracefulShutdown provides a shutdown function for the Server. Will try to shut down for 15s before returning.
func (s *Server) GracefulShutdown() {
	log.Printf("Shutting down crud server on %s\n", s.Server.Addr)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()

	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Printf("error shuting down crud server: %v\n", err)
	}
}
