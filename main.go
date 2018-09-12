package main

import (
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/controller"
	"log"
	"net/http"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	server := http.Server{
		Addr:         ":8080",
		Handler:      controller.NewBaseHTTPHandler(config.NewLocal("testDir")),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Exiting application.")
}
