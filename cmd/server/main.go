package main

import (
	"crypto/tls"
	"flag"
	"github.com/larwef/ki/internal/app"
	"github.com/larwef/ki/internal/config"
	"github.com/larwef/ki/internal/http/auth"
	"github.com/larwef/ki/internal/repository"
	"github.com/larwef/ki/internal/repository/local"
	"github.com/larwef/ki/internal/repository/memory"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

// PersistenceType defines available storage types
type PersistenceType string

const (
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

	config.Init(true, *propertyFile)

	// Setting bits since we want to be able to run multiple api types
	var apiType app.APIType
	crudEnabled, _ := config.GetBool("apiType.crud.enabled", false, false)
	if crudEnabled {
		apiType = apiType | app.CRUD
	}

	grpcEnabled, _ := config.GetBool("apiType.grpc.enabled", false, false)
	if grpcEnabled {
		apiType = apiType | app.GRPC
	}

	persistenceType, _ := config.GetString("persistence.type", true)
	var repo repository.Repository
	switch PersistenceType(persistenceType) {
	case Memory:
		repo = memory.NewRepository()
		log.Println("Using in memory storage")
		break
	case JSON:
		persistenceLocation, _ := config.GetString("persistence.location", true)
		repo = local.NewRepository(persistenceLocation)
		log.Println("Using JSON storage")
		break
	default:
		log.Fatal("Unsupported storage type")
	}

	var tlsConfig *tls.Config
	if !*disableTLS {
		certManager := getAutoCertManager()
		tlsConfig = getTLSConfig(certManager)
	} else {
		log.Println("TLS disabled")
	}

	options := []app.Option{
		app.TLSConfig(tlsConfig),
		app.Repository(repo),
		app.APITypes(apiType),
		app.Auth(getBasicAuth()),
	}

	app.NewApp(options...).Run()

	log.Println("Exiting application.")
}

func getTLSConfig(autoCertManager *autocert.Manager) *tls.Config {
	return &tls.Config{GetCertificate: autoCertManager.GetCertificate}
}

func getAutoCertManager() *autocert.Manager {
	acmeDirectoryURL, _ := config.GetString("tls.acme.directoryUrl", false)
	acmeClient := &acme.Client{
		DirectoryURL: acmeDirectoryURL,
	}

	certCache, _ := config.GetString("tls.acme.certCache", true)
	acmeHost, _ := config.GetString("tls.acme.host", false)
	acmeClientEmail, _ := config.GetString("tls.acme.client.email", false)

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(certCache),
		HostPolicy: autocert.HostWhitelist(acmeHost),
		Client:     acmeClient,
		Email:      acmeClientEmail,
	}

	go func() {
		// Listens for challenges from ACME provider
		if err := http.ListenAndServe(":http", certManager.HTTPHandler(nil)); err != nil {
			// TODO: If this listener fails the app will no longer be able to get new certificates from ACME provider. But
			// TODO as long as theres a valid cached certificate this shouldnt be a problem
			log.Fatalf("Error listening to port http: %v", err)
		}
	}()

	return certManager
}

func getBasicAuth() *auth.Basic {
	var basic *auth.Basic
	if basiAuthEnabled, err := config.GetBool("auth.basic.enabled", true, false); basiAuthEnabled && err == nil {
		basic = auth.NewBasic()
		adminUsername, _ := config.GetString("auth.admin.username", true)
		adminPasswordHash, _ := config.GetString("auth.admin.password", true)

		admin := auth.User{
			Username:     adminUsername,
			PasswordHash: adminPasswordHash,
			Role:         auth.ADMIN,
		}

		if err := basic.RegisterUser(admin); err != nil {
			log.Fatalf("Error adding admin user to user basic: %v", err)
		}

		clientUsername, _ := config.GetString("auth.client.username", true)
		clientPasswordHash, _ := config.GetString("auth.client.password", true)

		client := auth.User{
			Username:     clientUsername,
			PasswordHash: clientPasswordHash,
			Role:         auth.CLIENT,
		}

		if err := basic.RegisterUser(client); err != nil {
			log.Printf("Error adding client user to user basic: %v", err)
		}
	}

	return basic
}
