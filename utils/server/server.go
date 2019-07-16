package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/nats-server/server"
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
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
	auth           *auth.Module
	crud           *crud.Module
	user           *userman.Module
	file           *filestore.Module
	functions      *functions.Module
	realtime       *realtime.Module
	static         *static.Module
	adminMan       *admin.Manager
	isProd         bool
	nats           *nats.Server
	configFilePath string
	syncMan        *syncman.SyncManager
	ssl            *config.SSL
}

// New creates a new server instance
func New(isProd bool) *Server {
	r := mux.NewRouter()
	c := crud.Init()
	rt := realtime.Init(c)
	s := static.Init()
	fn := functions.Init()
	a := auth.Init(c, fn)
	u := userman.Init(c, a)
	f := filestore.Init(a)
	adminMan := admin.New()
	syncMan := syncman.New(adminMan)

	return &Server{nodeID: uuid.NewV1().String(), router: r, auth: a, crud: c,
		user: u, file: f, static: s, syncMan: syncMan, adminMan: adminMan,
		functions: fn, realtime: rt, isProd: isProd, configFilePath: utils.DefaultConfigFilePath}
}

// Start begins the server operations
func (s *Server) Start(port, grpcPort string, seeds string) error {
	// Start gRPC server in a separate goroutine
	go s.initGRPCServer(grpcPort)

	// Start the sync manager
	if seeds == "" {
		seeds = "127.0.0.1"
	}
	array := strings.Split(seeds, ",")
	if err := s.syncMan.Start(s.nodeID, s.configFilePath, "4232", "4234", array, s.LoadConfig); err != nil {
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
	if err := s.static.SetConfig(p.Modules.Static); err != nil {
		return err
	}

	// Set the configuration for the crud module
	return s.crud.SetConfig(p.Modules.Crud)
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
	pb.RegisterSpaceCloudServer(grpcServer, s)

	fmt.Println("Starting gRPC Server on port: " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}

// SetConfigFilePath sets the config file path
func (s *Server) SetConfigFilePath(configFilePath string) {
	s.configFilePath = configFilePath
}
