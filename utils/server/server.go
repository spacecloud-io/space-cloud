package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/nats-server/server"
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud/driver"
	"github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

// Server is the main server object
type Server struct {
	id       string
	lock     sync.RWMutex
	router   *mux.Router
	isProd   bool
	nats     *nats.Server
	projects *projects.Projects
	ssl      *config.SSL
}

// New creates a new server instance
func New(isProd bool) *Server {
	r := mux.NewRouter()
	projects := projects.New(driver.New())
	id := uuid.NewV1().String()
	return &Server{id: id, router: r, projects: projects, isProd: isProd}
}

// Start begins the server operations
func (s *Server) Start(port, grpcPort string) error {

	go s.initGRPCServer(grpcPort)

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

	fmt.Println("Starting http Server on port " + port)

	if s.ssl != nil {
		return http.ListenAndServeTLS(":"+port, s.ssl.Crt, s.ssl.Key, handler)
	}

	return http.ListenAndServe(":"+port, handler)
}

func (s *Server) initGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	options := []grpc.ServerOption{}
	if s.ssl != nil {
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
	return s.id
}
