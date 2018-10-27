package app

import (
	"crypto/tls"
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/config"
	"github.com/larwef/ki/internal/http/auth"
	"github.com/larwef/ki/internal/http/crud"
	"github.com/larwef/ki/internal/http/grpc"
	"github.com/larwef/ki/internal/listing"
	"github.com/larwef/ki/internal/repository"
	"github.com/larwef/ki/internal/repository/local"
	"github.com/larwef/ki/internal/runner"
	goGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"time"
)

// APIType defines available api types
type APIType int

const (
	// CRUD uses json over http
	CRUD APIType = 1 << iota
	// GRPC uses rpc with protobuf
	GRPC
)

// App handles some setup and running the app
type App struct {
	opts options
}

type options struct {
	tlsConfig  *tls.Config
	repository repository.Repository
	apiType    APIType
	aut        auth.Auth
}

var defaultAppOptions = options{
	repository: local.NewRepository("persistence"),
}

// Option sets options on the App object.
type Option func(*options)

// TLSConfig returns an Option that sets tls config for apis.
func TLSConfig(tc *tls.Config) Option { return func(o *options) { o.tlsConfig = tc } }

// Repository returns an Option that sets the storage.
func Repository(r repository.Repository) Option { return func(o *options) { o.repository = r } }

// APITypes returns an Option that sets which api types that should be active.
func APITypes(at APIType) Option { return func(o *options) { o.apiType = at } }

// Auth returns an Option that sets authentication implementation to be used.
func Auth(a auth.Auth) Option { return func(o *options) { o.aut = a } }

// NewApp returns a new app object
func NewApp(opt ...Option) *App {
	opts := defaultAppOptions
	for _, o := range opt {
		o(&opts)
	}

	return &App{
		opts: opts,
	}
}

// Run runs the application. Will start server objects for active APIs. CRUD and GRPC will use the same TLS config and the same
// persistence. So both can be used at the same time towards different clients and still provide the same resources to both.
func (a *App) Run() {
	add := adding.NewService(a.opts.repository)
	lst := listing.NewService(a.opts.repository)

	rnr := runner.NewRunner()

	// CRUD
	if a.opts.apiType&CRUD != 0 {
		crudAddress, _ := config.GetString("apiType.crud.address", true)
		crudServer := &crud.Server{
			Server: &http.Server{
				Addr:         crudAddress,
				Handler:      crud.NewHandler(a.opts.aut, add, lst),
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				TLSConfig:    a.opts.tlsConfig,
			},
		}
		rnr.Add(crudServer)
	}

	// gRPC
	if a.opts.apiType&GRPC != 0 {
		grpcAddress, _ := config.GetString("apiType.grpc.address", true)
		listener, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		var opts []goGrpc.ServerOption
		opts = append(opts, goGrpc.UnaryInterceptor(grpc.InOutLoggingUnaryInterceptor))

		if a.opts.tlsConfig != nil {
			log.Println("Enabling tls for gRPC")
			opts = append(opts, goGrpc.Creds(credentials.NewTLS(a.opts.tlsConfig)))
		}

		grpcServer := &grpc.Server{
			Server:   goGrpc.NewServer(opts...),
			Listener: listener,
			Handler:  grpc.NewHandler(add, lst),
		}

		rnr.Add(grpcServer)
	}

	rnr.Run()
}
