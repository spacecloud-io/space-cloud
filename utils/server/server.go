package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
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
)

type Server struct {
	Lock      sync.Mutex
	Router    *mux.Router
	Auth      *auth.Module
	Crud      *crud.Module
	User      *userman.Module
	File      *filestore.Module
	Functions *functions.Module
	Realtime  *realtime.Module
	Static    *static.Module
	IsProd    bool
	Config    *config.Project
}

func InitServer(isProd bool) *Server {
	r := mux.NewRouter()
	c := crud.Init()
	a := auth.Init(c)
	u := userman.Init(c, a)
	f := filestore.Init()
	rl := realtime.Init()
	s := static.Init()
	fn := functions.Init()

	return &Server{Router: r, Auth: a, Crud: c, User: u, File: f, Static: s, Functions: fn, Realtime: rl, IsProd: isProd}
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

	handler := corsObj.Handler(s.Router)

	fmt.Println("Starting http Server on port " + port)

	if s.Config.SSL != nil {
		return http.ListenAndServeTLS(":"+port, s.Config.SSL.Crt, s.Config.SSL.Key, handler)
	}

	return http.ListenAndServe(":"+port, handler)
}

func (s *Server) LoadConfig(config *config.Project) error {
	s.Lock.Lock()
	s.Config = config
	s.Lock.Unlock()

	// Set the configuration for the Auth module
	s.Auth.SetConfig(config.Secret, config.Modules.Crud, config.Modules.FileStore)

	// Set the configuration for the User management module
	s.User.SetConfig(config.Modules.Auth)

	// Set the configuration for the File storage module
	if err := s.File.SetConfig(config.Modules.FileStore); err != nil {
		return err
	}

	// Set the configuration for the Functions module
	if err := s.Functions.SetConfig(config.Modules.Functions); err != nil {
		return err
	}

	// Set the configuration for the Realtime module
	if err := s.Realtime.SetConfig(config.Modules.Realtime); err != nil {
		return err
	}

	// Set the configuration for Static module
	if err := s.Static.SetConfig(config.Modules.Static); err != nil {
		return err
	}

	// Set the configuration for the Crud module
	return s.Crud.SetConfig(config.Modules.Crud)
}

func (s *Server) initGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	var options []grpc.ServerOption
	if s.Config.SSL != nil {
		creds, err := credentials.NewServerTLSFromFile(s.Config.SSL.Crt, s.Config.SSL.Key)
		if err != nil {
			log.Fatalln("Error -", err)
		}
		options = append(options, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(options...)
	pb.RegisterSpaceCloudServer(grpcServer, s)

	fmt.Println("Starting gRPC Server on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
