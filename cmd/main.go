package main

import (
	"crypto/tls"
	"flag"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/http/crud"
	"github.com/larwef/ki/internal/http/grpc"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/internal/repository/local"
	"github.com/larwef/ki/internal/repository/memory"
	"github.com/larwef/ki/internal/runner"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	goGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"time"
)

// APIType defines available api types
type APIType int

// StorageType defines available storage types
type StorageType int

const stagingURL = "https://acme-staging.api.letsencrypt.org/directory"

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

	// TODO: Do better configuration
	// Use this for local testing
	disableTLS := flag.Bool("disable-tls", false, "Set TLS config on server objects")
	// Use acme-staging.api. Use this when testing setup to not risking hitting rate limit in prod.
	useStaging := flag.Bool("use-staging", false, "Use Let's encrypt staging api")
	flag.Parse()
	log.Printf("TLS Enabled: %t", !*disableTLS)

	apiType := CRUD | GRPC
	storageType := Memory
	path := "testDir"
	host := "tlstest.wefald.no"

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

	// TLS stuff
	var tlsConfig *tls.Config
	if !*disableTLS {
		// Cached certificates are stored here
		dataDir := "certCache"

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(host),
			Cache:      autocert.DirCache(dataDir),
		}

		if *useStaging {
			m.Client = &acme.Client{DirectoryURL: stagingURL}
		}

		tlsConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			// Listens for challenges from Let's encrypt
			if err := http.ListenAndServe(":http", m.HTTPHandler(nil)); err != nil {
				log.Fatalf("Error listening to port http: %v", err)
			}
		}()
	}
	// End TLS stuff

	rnr := runner.NewRunner()

	// CRUD
	if apiType&CRUD != 0 {
		crudServer := &crud.Server{
			Server: &http.Server{
				Addr:         ":8080",
				Handler:      crud.NewHandler(add, lst),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				TLSConfig:    tlsConfig,
			},
		}
		rnr.Add(crudServer)
	}

	// gRPC
	if apiType&GRPC != 0 {
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		var opts []goGrpc.ServerOption
		opts = append(opts, goGrpc.UnaryInterceptor(grpc.InOutLoggingUnaryInterceptor))

		if tlsConfig != nil {
			log.Println("Enabling tls for gRPC")
			opts = append(opts, goGrpc.Creds(credentials.NewTLS(tlsConfig)))
		}

		grpcServer := &grpc.Server{
			Server:   goGrpc.NewServer(opts...),
			Listener: listener,
			Handler:  grpc.NewHandler(add, lst),
		}

		rnr.Add(grpcServer)
	}

	rnr.Run()

	log.Println("Exiting application.")
}
