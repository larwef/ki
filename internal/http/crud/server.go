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
func (s *Server) Serve(signal chan bool) {
	log.Printf("Starting crud server on %s\n", s.Server.Addr)
	if s.Server.TLSConfig != nil {
		log.Println("Listening and serving tls")
		if err := s.Server.ListenAndServeTLS("", ""); err != nil {
			log.Printf("crud server error: %v\n", err)
			signal <- true
		}
	} else {
		if err := s.Server.ListenAndServe(); err != nil {
			log.Printf("crud server error: %v\n", err)
			signal <- true
		}
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
