package main

import (
	"crypto/tls"
	"flag"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/config"
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting Ki. Version: %s\n", Version)

	// Config
	disableTLS := flag.Bool("disable-tls", false, "Set TLS config on server objects")
	propertyFile := flag.String("property-file", "", "Set if properties should be gotten from file")
	flag.Parse()

	// ReadEnv will overwrite the properties from file
	if *propertyFile != "" {
		config.ReadPorpertyFile(*propertyFile)
	}
	config.ReadEnv()

	// Setting bits since we want to be able to run multiple api types
	var apiType APIType
	crudEnabled, _ := config.GetBool("apiType.crud.enabled", false, false)
	if crudEnabled {
		apiType = apiType | CRUD
	}

	grpcEnabled, _ := config.GetBool("apiType.grpc.enabled", false, false)
	if grpcEnabled {
		apiType = apiType | GRPC
	}

	var add adding.Service
	var lst listing.Service
	persistenceType, _ := config.GetString("persistence.type", true)
	switch PersistenceType(persistenceType) {
	case Memory:
		repo := memory.NewRepository()
		add = adding.NewService(repo)
		lst = listing.NewService(repo)
		log.Println("Using in memory storage")
		break
	case JSON:
		persistenceLocation, _ := config.GetString("persistence.location", true)
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

		acmeDirectoryURL, _ := config.GetString("tls.acme.directoryUrl", false)
		acmeClient := &acme.Client{
			DirectoryURL: acmeDirectoryURL,
		}

		certCache, _ := config.GetString("tls.acme.certCache", true)
		acmeHost, _ := config.GetString("tls.acme.host", false)
		acmeClientEmail, _ := config.GetString("tls.acme.client.email", false)
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(certCache),
			HostPolicy: autocert.HostWhitelist(acmeHost),
			Client:     acmeClient,
			Email:      acmeClientEmail,
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
		log.Println("TLS disabled")
	}

	rnr := runner.NewRunner()

	// CRUD
	if apiType&CRUD != 0 {
		crudAddress, _ := config.GetString("apiType.crud.address", true)
		crudServer := &crud.Server{
			Server: &http.Server{
				Addr:         crudAddress,
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
		grpcAddress, _ := config.GetString("apiType.grpc.address", true)
		listener, err := net.Listen("tcp", grpcAddress)
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
