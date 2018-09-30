package main

import (
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/http/crud"
	"github.com/larwef/ki/internal/http/grpc"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/internal/repository/local"
	"github.com/larwef/ki/internal/repository/memory"
	"github.com/larwef/ki/internal/runner"
	goGrpc "google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

// APIType defines available api types
type APIType int

// StorageType defines available storage types
type StorageType int

const (
	// CRUD uses json over http
	CRUD APIType = 1 << iota
	// GRPC uses rpc with protobuf
	GRPC

	// Memory will store data in memory
	Memory StorageType = iota
	// JSON will store data in JSON files saved on disk
	JSON
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	// TODO: Make configurable
	apiType := CRUD | GRPC
	storageType := Memory
	path := "testDir"

	var add adding.Service
	var lst listing.Service
	switch storageType {
	case Memory:
		repo := memory.NewRepository()
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		log.Println("Using in memory storage")
		break
	case JSON:
		repo := local.NewRepository(path)
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		log.Println("Using JSON storage")
		break
	default:
		log.Fatal("Unsupported storage type")
	}

	rnr := runner.NewRunner()

	// CRUD
	if apiType&CRUD != 0 {
		crudServer := &crud.Server{
			Server: &http.Server{
				Addr:         ":8080",
				Handler:      crud.NewHandler(add, lst),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			},
		}
		rnr.Add(crudServer)
	}

	// gRPC
	if apiType&GRPC != 0 {
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := &grpc.Server{
			Server:   goGrpc.NewServer(),
			Listener: listener,
			Handler:  grpc.NewHandler(add, lst),
		}

		rnr.Add(grpcServer)
	}

	rnr.Run()

	log.Println("Exiting application.")
}
