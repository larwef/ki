package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// Server represents a grpc server object
type Server struct {
	Server   *grpc.Server
	Listener net.Listener
	Handler  *Handler
}

// Serve starts listening on the grpc server. Sends a signal on error.
// TODO: Improvement: When GracefulStop is called Serve will return, hopefully without error. Should it return with an error
// TODO: the program will panic (since this will call close on a closed channel). So this could be executed cleaner.
func (s *Server) Serve(signal chan bool) {
	RegisterGroupServiceServer(s.Server, s.Handler)
	RegisterConfigServiceServer(s.Server, s.Handler)
	reflection.Register(s.Server)

	log.Printf("Starting grpc server on %s\n", s.Listener.Addr().String())
	if err := s.Server.Serve(s.Listener); err != nil {
		log.Printf("grpc error: %v\n", err)
		close(signal)
	}
}

// GracefulShutdown provides a shutdown function for the Server.
func (s *Server) GracefulShutdown() {
	log.Printf("Shutting down grpc server on %s\n", s.Listener.Addr().String())
	s.Server.GracefulStop()
}
