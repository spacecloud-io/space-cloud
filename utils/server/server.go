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
	uuid "github.com/satori/go.uuid"
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
func New(isProd bool) *Server {
	nodeID := uuid.NewV1().String()
	r := mux.NewRouter()
	d := deploy.New()
	adminMan := admin.New(nodeID)
	projects := projects.New(driver.New())
	syncMan := syncman.New(projects, d, adminMan)
	return &Server{nodeID: nodeID, router: r, projects: projects, isProd: isProd,
		syncMan: syncMan, adminMan: adminMan, configFilePath: utils.DefaultConfigFilePath,
		deploy: d,
	}
}

// Start begins the server operations
func (s *Server) Start(port, grpcPort, seeds string) error {

	go s.initGRPCServer(grpcPort)

	// Start the sync manager
	if seeds == "" {
		seeds = "127.0.0.1"
	}
	array := strings.Split(seeds, ",")
	if err := s.syncMan.Start(s.nodeID, s.configFilePath, "4232", "4234", array); err != nil {
		return err
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

	handler := corsObj.Handler(s.router)

	fmt.Println("Starting HTTP Server on port: " + port)

	fmt.Println("Space Cloud is running on the specified ports :D")
	if s.ssl != nil && s.ssl.Enabled {
		return http.ListenAndServeTLS(":"+port, s.ssl.Crt, s.ssl.Key, handler)
	}

	return http.ListenAndServe(":"+port, handler)
}

// SetConfig sets the config
func (s *Server) SetConfig(c *config.Config) {
	s.ssl = c.SSL
	s.syncMan.SetGlobalConfig(c)
	s.adminMan.SetConfig(c.Admin)
	s.deploy.SetConfig(&c.Deploy)
}

func (s *Server) initGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	options := []grpc.ServerOption{}
	if s.ssl != nil && s.ssl.Enabled {
		creds, err := credentials.NewServerTLSFromFile(s.ssl.Crt, s.ssl.Key)
		if err != nil {
			log.Fatalln("Error: ", err)
		}
		options = append(options, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(options...)
	proto.RegisterSpaceCloudServer(grpcServer, s)

	fmt.Println("Starting gRPC Server on port: " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
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
