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
	"os"
	"strconv"
	"time"
)

// APIType defines available api types
type APIType int

// StorageType defines available storage types
type StorageType string

const (
	// CRUD uses json over http
	CRUD APIType = 1 << iota
	// GRPC uses rpc with protobuf
	GRPC

	// Memory will store data in memory
	Memory StorageType = "memory"
	// JSON will store data in JSON files saved on disk
	JSON StorageType = "json"
)

// TODO: Clean main function
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	// Flags
	// Use this for local testing
	disableTLS := flag.Bool("disable-tls", false, "Set TLS config on server objects")
	flag.Parse()

	// Environment variables
	host := os.Getenv("tls.acme.host")
	certCache := os.Getenv("tls.acme.certCache")
	acmeDirectoryURL := os.Getenv("tls.acme.directoryUrl")
	acmeClientEmail := os.Getenv("tls.acme.client.email")

	persistenceType := os.Getenv("persistence.type")
	persistenceLocation := os.Getenv("persistence.location")

	apiTypeCrudEnabled := os.Getenv("apiType.crud.enabled")
	apiTypeCrudAddr := os.Getenv("apiType.crud.address")
	apiTypeGrpcEnabled := os.Getenv("apiType.grpc.enabled")
	apiTypeGrpcAddr := os.Getenv("apiType.grpc.address")

	// Setting bits since we want to be able to run multiple api types
	var apiType APIType
	if b, _ := strconv.ParseBool(apiTypeCrudEnabled); b {
		apiType = apiType | CRUD
	}
	if b, _ := strconv.ParseBool(apiTypeGrpcEnabled); b {
		apiType = apiType | GRPC
	}

	var add adding.Service
	var lst listing.Service
	switch StorageType(persistenceType) {
	case Memory:
		repo := memory.NewRepository()
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		log.Println("Using in memory storage")
		break
	case JSON:
		repo := local.NewRepository(persistenceLocation)
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		log.Println("Using JSON storage")
		break
	default:
		log.Fatal("Unsupported storage type")
	}

	var tlsConfig *tls.Config
	if !*disableTLS {

		acmeClient := &acme.Client{
			DirectoryURL: acmeDirectoryURL,
		}

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(certCache),
			HostPolicy: autocert.HostWhitelist(host),
			Client:     acmeClient,
			Email:      acmeClientEmail,
		}

		tlsConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			// Listens for challenges from Let's encrypt
			if err := http.ListenAndServe(":http", m.HTTPHandler(nil)); err != nil {
				log.Fatalf("Error listening to port http: %v", err)
			}
		}()
	} else {
		log.Println("TLS diabled")
	}

	rnr := runner.NewRunner()

	// CRUD
	if apiType&CRUD != 0 {
		crudServer := &crud.Server{
			Server: &http.Server{
				Addr:         apiTypeCrudAddr,
				Handler:      crud.NewHandler(add, lst),
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				TLSConfig:    tlsConfig,
			},
		}
		rnr.Add(crudServer)
	}

	// gRPC
	if apiType&GRPC != 0 {
		listener, err := net.Listen("tcp", apiTypeGrpcAddr)
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
