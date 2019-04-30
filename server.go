package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/modules/filestore"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/modules/userman"
	pb "github.com/spaceuptech/space-cloud/proto"
)

type server struct {
	lock     sync.Mutex
	router   *mux.Router
	auth     *auth.Module
	crud     *crud.Module
	user     *userman.Module
	file     *filestore.Module
	functions     *functions.Module
	realtime *realtime.Module
	isProd   bool
	config   *config.Project
}

func initServer(isProd bool) *server {
	r := mux.NewRouter()
	c := crud.Init()
	a := auth.Init(c)
	u := userman.Init(c, a)
	f := filestore.Init()
	realtime := realtime.Init()
	functions := functions.Init()
	return &server{router: r, auth: a, crud: c, user: u, file: f, functions: functions, realtime: realtime, isProd: isProd}
}

func (s *server) start(port string) error {

	portInt, _ := strconv.Atoi(port)
	go s.initGRPCServer(strconv.Itoa(portInt + 1))

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
	s.auth.SetConfig(config.Secret, config.Modules.Crud, config.Modules.FileStore)

	// Set the configuration for the user management module
	s.user.SetConfig(config.Modules.Auth)

	// Set the configuration for the file storage module
	s.file.SetConfig(config.Modules.FileStore)

	// Set the configuration for the Functions module
	err := s.functions.SetConfig(config.Modules.Functions)
	if err != nil {
		return err
	}

	// Set the configuration for the Realtime module
	s.realtime.SetConfig(config.Modules.Realtime)

	// Set the configuration for the curd module
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
