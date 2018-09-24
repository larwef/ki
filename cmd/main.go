package main

import (
	"github.com/larwef/ki/pkg/adding"
	"github.com/larwef/ki/pkg/controller"
	"github.com/larwef/ki/pkg/listing"
	"github.com/larwef/ki/pkg/repository/local"
	"github.com/larwef/ki/pkg/repository/memory"
	"log"
	"net/http"
	"time"
)

// StorageType defines available storage types
type StorageType int

const (
	// Memory will store data in memory
	Memory StorageType = iota
	// JSON will store data in JSON files saved on disk
	JSON
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	// TODO: Make configurable
	storageType := Memory
	path := "testDir"

	var add adding.Service
	var lst listing.Service
	switch storageType {
	case JSON:
		repo := local.NewRepository(path)
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		break
	case Memory:
		repo := memory.NewRepository()
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		break
	default:
		log.Fatal("Unsupported storage type")
	}

	server := http.Server{
		Addr:         ":8080",
		Handler:      controller.NewBaseHTTPHandler(add, lst),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Exiting application.")
}
