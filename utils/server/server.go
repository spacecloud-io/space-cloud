package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/nats-server/v2/server"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/filestore"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/modules/static"
	"github.com/spaceuptech/space-cloud/modules/userman"
	pb "github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	nodeID         string
	router         *mux.Router
	routerSecure   *mux.Router
	auth           *auth.Module
	crud           *crud.Module
	user           *userman.Module
	file           *filestore.Module
	functions      *functions.Module
	realtime       *realtime.Module
	static         *static.Module
	adminMan       *admin.Manager
	nats           *nats.Server
	configFilePath string
	syncMan        *syncman.SyncManager
	ssl            *config.SSL
}

// New creates a new server instance
func New(nodeID string) *Server {
	r := mux.NewRouter()
	r2 := mux.NewRouter()
	c := crud.Init()
	rt := realtime.Init(c)
	s := static.Init()
	fn := functions.Init()
	a := auth.Init(c, fn)
	u := userman.Init(c, a)
	f := filestore.Init(a)
	adminMan := admin.New()
	syncMan := syncman.New(adminMan)

	fmt.Println("Creating a new server with id", nodeID)

	return &Server{nodeID: nodeID, router: r, routerSecure: r2, auth: a, crud: c,
		user: u, file: f, static: s, syncMan: syncMan, adminMan: adminMan,
		functions: fn, realtime: rt, configFilePath: utils.DefaultConfigFilePath}
}

// Start begins the server operations
func (s *Server) Start(seeds string, disableMetrics bool) error {
	// Start gRPC server in a separate goroutine
	go s.initGRPCServer()

	// Start the sync manager
	if seeds == "" {
		seeds = "127.0.0.1"
	}
	array := strings.Split(seeds, ",")
	if err := s.syncMan.Start(s.nodeID, s.configFilePath, utils.PortGossip, utils.PortRaft, array, s.LoadConfig); err != nil {
		return err
	}

	// Anonymously collect usage metrics if not explicitly disabled
	if !disableMetrics {
		go s.RoutineMetrics()
	}

	// Allow cors
	corsObj := cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc: func(s string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		ExposedHeaders: []string{"Authorization", "Content-Type"},
	})

	fmt.Println("Starting http server on port: " + utils.PortHTTP)

	if s.ssl != nil && s.ssl.Enabled {
		handler := corsObj.Handler(s.routerSecure)
		fmt.Println("Starting https server on port: " + utils.PortHTTPSecure)
		go func() {

			if err := http.ListenAndServeTLS(":"+utils.PortHTTPSecure, s.ssl.Crt, s.ssl.Key, handler); err != nil {
				fmt.Println("Error starting https server:", err)
			}
		}()
	}

	handler := corsObj.Handler(s.router)

	fmt.Println()
	fmt.Println("\t Hosting mission control on http://localhost:" + utils.PortHTTP + "/mission-control/")
	fmt.Println()

	fmt.Println("Space cloud is running on the specified ports :D")
	return http.ListenAndServe(":"+utils.PortHTTP, handler)
}

// SetConfig sets the config
func (s *Server) SetConfig(c *config.Config, isProd bool) {
	s.ssl = c.SSL
	s.syncMan.SetGlobalConfig(c)
	s.adminMan.SetEnv(isProd)
	s.adminMan.SetConfig(c.Admin)
}

// LoadConfig configures each module to to use the provided config
func (s *Server) LoadConfig(config *config.Config) error {

	if config.Projects == nil {
		return errors.New("No projects provided")
	}

	p := config.Projects[0]

	// Set the configuration for the auth module
	s.auth.SetConfig(p.ID, p.Secret, p.Modules.Crud, p.Modules.FileStore, p.Modules.Functions)

	// Set the configuration for the user management module
	s.user.SetConfig(p.Modules.Auth)

	// Set the configuration for the file storage module
	if err := s.file.SetConfig(p.Modules.FileStore); err != nil {
		return err
	}

	// Set the configuration for the functions module
	if err := s.functions.SetConfig(p.Modules.Functions); err != nil {
		return err
	}

	// Set the configuration for the realtime module
	if err := s.realtime.SetConfig(p.ID, p.Modules.Realtime); err != nil {
		return err
	}

	// Set the configuration for static module
	if err := s.static.SetConfig(config.Static); err != nil {
		return err
	}

	// Set the configuration for the crud module
	return s.crud.SetConfig(p.Modules.Crud)
}

func (s *Server) initGRPCServer() {

	if s.ssl != nil && s.ssl.Enabled {
		lis, err := net.Listen("tcp", ":"+utils.PortGRPCSecure)
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}
		creds, err := credentials.NewServerTLSFromFile(s.ssl.Crt, s.ssl.Key)
		if err != nil {
			log.Fatalln("Error: ", err)
		}
		options := []grpc.ServerOption{grpc.Creds(creds)}

		grpcServer := grpc.NewServer(options...)
		pb.RegisterSpaceCloudServer(grpcServer, s)

		fmt.Println("Starting grpc secure server on port: " + utils.PortGRPCSecure)
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatal("Error starting grpc secure server:", err)
			}
		}()
	}

	lis, err := net.Listen("tcp", ":"+utils.PortGRPC)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	options := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(options...)
	pb.RegisterSpaceCloudServer(grpcServer, s)

	fmt.Println("Starting grpc server on port: " + utils.PortGRPC)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Error starting grpc server:", err)
	}
}

// SetConfigFilePath sets the config file path
func (s *Server) SetConfigFilePath(configFilePath string) {
	s.configFilePath = configFilePath
}
