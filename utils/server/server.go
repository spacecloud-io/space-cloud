package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/nats-server/v2/server"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud/driver"
	"github.com/spaceuptech/space-cloud/modules/deploy"
	"github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/projects"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// Server is the main server object
type Server struct {
	nodeID         string
	lock           sync.RWMutex
	router         *mux.Router
	routerSecure   *mux.Router
	isProd         bool
	nats           *nats.Server
	projects       *projects.Projects
	ssl            *config.SSL
	syncMan        *syncman.SyncManager
	configFilePath string
	adminMan       *admin.Manager
	deploy         *deploy.Module
}

// New creates a new server instance
func New(nodeID string, isProd bool) *Server {
	r := mux.NewRouter()
	r2 := mux.NewRouter()
	d := deploy.New()
	adminMan := admin.New(nodeID)
	projects := projects.New(driver.New())
	syncMan := syncman.New(projects, d, adminMan)
	return &Server{nodeID: nodeID, router: r, routerSecure: r2, projects: projects, isProd: isProd,
		syncMan: syncMan, adminMan: adminMan, configFilePath: utils.DefaultConfigFilePath,
		deploy: d,
	}
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
	if err := s.syncMan.Start(s.nodeID, s.configFilePath, utils.PortGossip, utils.PortRaft, array); err != nil {
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
func (s *Server) SetConfig(c *config.Config) {
	s.ssl = c.SSL
	s.syncMan.SetGlobalConfig(c)
	s.adminMan.SetConfig(c.Admin)
	s.deploy.SetConfig(&c.Deploy)
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
		proto.RegisterSpaceCloudServer(grpcServer, s)

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
	proto.RegisterSpaceCloudServer(grpcServer, s)

	fmt.Println("Starting grpc server on port: " + utils.PortGRPC)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Error starting grpc server:", err)
	}
}

// GetProjects returns a copy of the projects
func (s *Server) GetProjects() *projects.Projects {
	return s.projects
}

// GetID returns the server id
func (s *Server) GetID() string {
	return s.nodeID
}

// SetConfigFilePath sets the config file path
func (s *Server) SetConfigFilePath(configFilePath string) {
	s.configFilePath = configFilePath
}
