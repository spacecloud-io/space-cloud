package main

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
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/filestore"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/modules/static"
	"github.com/spaceuptech/space-cloud/modules/userman"
	pb "github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
)

type server struct {
	lock           sync.Mutex
	router         *mux.Router
	auth           *auth.Module
	crud           *crud.Module
	user           *userman.Module
	file           *filestore.Module
	functions      *functions.Module
	realtime       *realtime.Module
	static         *static.Module
	isProd         bool
	config         *config.Project
	nats           *nats.Server
	configFilePath string
}

func initServer(isProd bool) *server {
	r := mux.NewRouter()
	c := crud.Init()
	f := filestore.Init()
	realtime := realtime.Init(c)
	s := static.Init()
	functions := functions.Init()
	a := auth.Init(c, functions)
	u := userman.Init(c, a)
	return &server{router: r, auth: a, crud: c, user: u, file: f, static: s, functions: functions, realtime: realtime, isProd: isProd, configFilePath: utils.DefaultConfigFilePath}
}

func (s *server) start(port, grpcPort string) error {

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

	fmt.Println("Starting http server on port " + port)

	if s.config.SSL != nil {
		return http.ListenAndServeTLS(":"+port, s.config.SSL.Crt, s.config.SSL.Key, handler)
	}

	return http.ListenAndServe(":"+port, handler)
}

func (s *server) loadConfig(config *config.Project) error {
	s.lock.Lock()
	s.config = config
	s.lock.Unlock()

	// Set the configuration for the auth module
	s.auth.SetConfig(config.ID, config.Secret, config.Modules.Crud, config.Modules.FileStore, config.Modules.Functions)

	// Set the configuration for the user management module
	s.user.SetConfig(config.Modules.Auth)

	// Set the configuration for the file storage module
	if err := s.file.SetConfig(config.Modules.FileStore); err != nil {
		return err
	}

	// Set the configuration for the functions module
	if err := s.functions.SetConfig(config.Modules.Functions); err != nil {
		return err
	}

	// Set the configuration for the Realtime module
	if err := s.realtime.SetConfig(config.ID, config.Modules.Realtime); err != nil {
		return err
	}

	// Set the configuration for Static module
	if err := s.static.SetConfig(config.Modules.Static); err != nil {
		return err
	}

	// Set the configuration for the crud module
	return s.crud.SetConfig(config.Modules.Crud)
}

func (s *server) initGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	options := []grpc.ServerOption{}
	if s.config.SSL != nil {
		creds, err := credentials.NewServerTLSFromFile(s.config.SSL.Crt, s.config.SSL.Key)
		if err != nil {
			log.Fatalln("Error -", err)
		}
		options = append(options, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(options...)
	pb.RegisterSpaceCloudServer(grpcServer, s)

	fmt.Println("Starting gRPC server on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
