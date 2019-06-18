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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/static"
	"github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

type Server struct {
	lock     sync.RWMutex
	router   *mux.Router
	isProd   bool
	nats     *nats.Server
	static   *static.Module
	projects *projects.Projects
	ssl      *config.SSL
}

func New(isProd bool) *Server {
	r := mux.NewRouter()
	s := static.Init()
	projects := projects.New()
	return &Server{router: r, static: s, projects: projects, isProd: isProd}
}

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

// Projects returns a copy of the projects
func (s *Server) Projects() *projects.Projects {
	return s.projects
}
