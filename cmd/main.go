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

// PersistenceType defines available storage types
type PersistenceType string

const (
	// CRUD uses json over http
	CRUD APIType = 1 << iota
	// GRPC uses rpc with protobuf
	GRPC

	// Memory will store data in memory
	Memory PersistenceType = "memory"
	// JSON will store data in JSON files saved on disk
	JSON PersistenceType = "json"
)

// Version specifies application version. Set at build time.
var Version = "No version provided"

var appConfig = struct {
	disableTLS *bool

	host             string
	certCache        string
	acmeDirectoryURL string
	acmeClientEmail  string

	persistenceType     PersistenceType
	persistenceLocation string

	apiTypeCrudEnabled bool
	apiTypeCrudAddr    string
	apiTypeGrpcEnabled bool
	apiTypeGrpcAddr    string
}{
	// Use this for local testing
	disableTLS: flag.Bool("disable-tls", false, "Set TLS config on server objects"),

	host:             os.Getenv("tls.acme.host"),
	certCache:        os.Getenv("tls.acme.certCache"),
	acmeDirectoryURL: os.Getenv("tls.acme.directoryUrl"),
	acmeClientEmail:  os.Getenv("tls.acme.client.email"),

	persistenceType:     PersistenceType(os.Getenv("persistence.type")),
	persistenceLocation: os.Getenv("persistence.location"),

	apiTypeCrudEnabled: parseEnvBool(os.Getenv("apiType.crud.enabled")),
	apiTypeCrudAddr:    os.Getenv("apiType.crud.address"),
	apiTypeGrpcEnabled: parseEnvBool(os.Getenv("apiType.grpc.enabled")),
	apiTypeGrpcAddr:    os.Getenv("apiType.grpc.address"),
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting Ki. Version: %s\n", Version)

	flag.Parse()

	// Setting bits since we want to be able to run multiple api types
	var apiType APIType
	if appConfig.apiTypeCrudEnabled {
		apiType = apiType | CRUD
	}
	if appConfig.apiTypeGrpcEnabled {
		apiType = apiType | GRPC
	}

	var add adding.Service
	var lst listing.Service
	switch appConfig.persistenceType {
	case Memory:
		repo := memory.NewRepository()
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		log.Println("Using in memory storage")
		break
	case JSON:
		repo := local.NewRepository(appConfig.persistenceLocation)
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		log.Println("Using JSON storage")
		break
	default:
		log.Fatal("Unsupported storage type")
	}

	var tlsConfig *tls.Config
	if !*appConfig.disableTLS {

		acmeClient := &acme.Client{
			DirectoryURL: appConfig.acmeDirectoryURL,
		}

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(appConfig.certCache),
			HostPolicy: autocert.HostWhitelist(appConfig.host),
			Client:     acmeClient,
			Email:      appConfig.acmeClientEmail,
		}

		tlsConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			// Listens for challenges from Let's encrypt
			if err := http.ListenAndServe(":http", m.HTTPHandler(nil)); err != nil {
				// TODO: If this listener fails the app will no longer be able to get new certificates from ACME provider. But
				// TODO as long as theres a valid cached certificate this shouldnt be a problem
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
				Addr:         appConfig.apiTypeCrudAddr,
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
		listener, err := net.Listen("tcp", appConfig.apiTypeGrpcAddr)
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

func parseEnvBool(variable string) bool {
	b, err := strconv.ParseBool(variable)

	if err != nil {
		log.Printf("Error parsing environment variable: %q. Returning false", variable)
	}

	return b
}
