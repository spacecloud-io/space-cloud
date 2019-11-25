package server

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/spaceuptech/space-cloud/modules/crud/driver"
	"github.com/spaceuptech/space-cloud/utils/handlers"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/metrics"
	"github.com/spaceuptech/space-cloud/utils/projects"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	lock sync.RWMutex

	nodeID         string
	configFilePath string

	router       *mux.Router
	routerSecure *mux.Router

	adminMan *admin.Manager
	syncMan  *syncman.Manager
	metrics  *metrics.Module

	projects *projects.Projects
	ssl      *config.SSL
}

// New creates a new server instance
func New(nodeID, clusterID string, isConsulEnabled, removeProjectScope bool, metricsConfig *metrics.Config) (*Server, error) {
	r := mux.NewRouter()
	r2 := mux.NewRouter()

	m, err := metrics.New(nodeID, metricsConfig)
	if err != nil {
		return nil, err
	}

	adminMan := admin.New(nodeID)
	syncMan, err := syncman.New(nodeID, clusterID, isConsulEnabled, adminMan)
	if err != nil {
		return nil, err
	}

	p := projects.New(nodeID, removeProjectScope, driver.New(removeProjectScope), adminMan, syncMan, m)
	syncMan.SetProjectCallbacks(&model.ProjectCallbacks{
		Store: p.StoreProject,

		SetGlobalConfig:      p.SetGlobalConfig,
		SetCrudConfig:        p.SetCrudConfig,
		SetServicesConfig:    p.SetServicesConfig,
		SetFileStorageConfig: p.SetFileStoreConfig,
		SetEventingConfig:    p.SetEventingConfig,
		SetUserManConfig:     p.SetUserManConfig,

		Delete:     p.DeleteProject,
		ProjectIDs: p.GetProjectIDs,
	})

	fmt.Println("Creating a new server with id", nodeID)

	return &Server{nodeID: nodeID, router: r, routerSecure: r2, projects: p,
		syncMan: syncMan, adminMan: adminMan, configFilePath: utils.DefaultConfigFilePath,
	}, nil
}

// Start begins the server operations
func (s *Server) Start(disableMetrics bool, port int) error {

	// Start the sync manager
	if err := s.syncMan.Start(s.configFilePath); err != nil {
		return err
	}

	// Anonymously collect usage metrics if not explicitly disabled
	if !disableMetrics {
		go s.RoutineMetrics()
	}

	// Allow cors
	corsObj := utils.CreateCorsObject()

	fmt.Println("Starting http server on port: " + strconv.Itoa(port))

	if s.ssl != nil && s.ssl.Enabled {
		handler := corsObj.Handler(s.routerSecure)
		fmt.Println("Starting https server on port: " + strconv.Itoa(port+4))
		go func() {

			if err := http.ListenAndServeTLS(":"+strconv.Itoa(port+4), s.ssl.Crt, s.ssl.Key, handlers.HandleMetricMiddleWare(handler, s.metrics)); err != nil {
				fmt.Println("Error starting https server:", err)
			}
		}()
	}

	handler := corsObj.Handler(s.router)

	fmt.Println()
	fmt.Println("\t Hosting mission control on http://localhost:" + strconv.Itoa(port) + "/mission-control/")
	fmt.Println()

	fmt.Println("Space cloud is running on the specified ports :D")
	return http.ListenAndServe(":"+strconv.Itoa(port), handlers.HandleMetricMiddleWare(handler, s.metrics))
}

// SetConfig sets the config
func (s *Server) SetConfig(c *config.Config, isProd bool) {
	s.ssl = c.SSL
	s.syncMan.SetGlobalConfig(c)
	s.adminMan.SetEnv(isProd)
	s.adminMan.SetConfig(c.Admin)
}

// SetConfigFilePath sets the config file path
func (s *Server) SetConfigFilePath(configFilePath string) {
	s.configFilePath = configFilePath
}
