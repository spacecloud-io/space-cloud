package main

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/faas"
	"github.com/spaceuptech/space-cloud/modules/filestore"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/modules/userman"
	pb "github.com/spaceuptech/space-cloud/proto"
)

type server struct {
	router   *mux.Router
	auth     *auth.Module
	crud     *crud.Module
	user     *userman.Module
	file     *filestore.Module
	faas     *faas.Module
	realtime *realtime.Module
	isProd   bool
}

func initServer(isProd bool) *server {
	r := mux.NewRouter()
	c := crud.Init()
	a := auth.Init(c)
	u := userman.Init(c, a)
	f := filestore.Init()
	realtime := realtime.Init()
	faas := faas.Init()
	return &server{r, a, c, u, f, faas, realtime, isProd}
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
	return http.ListenAndServe(":"+port, handler)
}

func (s *server) loadConfig(config *config.Project) error {
	// Set the configuration for the auth module
	s.auth.SetConfig(config.Secret, config.Modules.Crud, config.Modules.FileStore)

	// Set the configuration for the user management module
	s.user.SetConfig(config.Modules.Auth)

	// Set the configuration for the file storage module
	s.file.SetConfig(config.Modules.FileStore)

	// Set the configuration for the FaaS module
	err := s.faas.SetConfig(config.Modules.FaaS)
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
	grpcServer := grpc.NewServer()
	pb.RegisterSpaceCloudServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
