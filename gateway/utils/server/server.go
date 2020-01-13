package server

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/driver"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/handlers"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	lock sync.RWMutex

	nodeID         string
	configFilePath string

	adminMan    *admin.Manager
	syncMan     *syncman.Manager
	metrics     *metrics.Module
	letsencrypt *letsencrypt.LetsEncrypt
	routing     *routing.Routing

	projects *projects.Projects
	ssl      *config.SSL
}

// New creates a new server instance
func New(nodeID, clusterID, advertiseAddr, storeType string, removeProjectScope bool, metricsConfig *metrics.Config) (*Server, error) {

	m, err := metrics.New(nodeID, metricsConfig)
	if err != nil {
		return nil, err
	}

	adminMan := admin.New(nodeID)
	syncMan, err := syncman.New(nodeID, clusterID, advertiseAddr, storeType, adminMan)
	if err != nil {
		return nil, err
	}

	// Initialise a lets encrypt client
	le, err := letsencrypt.New()
	if err != nil {
		return nil, err
	}

	// Initialise the routing module
	r := routing.New()

	p := projects.New(nodeID, removeProjectScope, driver.New(removeProjectScope), adminMan, syncMan, m, le, r)
	syncMan.SetProjectCallbacks(&model.ProjectCallbacks{
		Store:            p.StoreProject,
		StoreIgnoreError: p.StoreIgnoreErrors,

		SetGlobalConfig:      p.SetGlobalConfig,
		SetCrudConfig:        p.SetCrudConfig,
		SetServicesConfig:    p.SetServicesConfig,
		SetFileStorageConfig: p.SetFileStoreConfig,
		SetEventingConfig:    p.SetEventingConfig,
		SetUserManConfig:     p.SetUserManConfig,
		SetLetsencryptConfig: p.SetLetsencryptDomains,
		SetRoutingConfig:     p.SetProjectRoutes,

		Delete:     p.DeleteProject,
		ProjectIDs: p.GetProjectIDs,
	})

	logrus.Infoln("Creating a new server with id", nodeID)

	return &Server{nodeID: nodeID, projects: p, letsencrypt: le, routing: r,
		syncMan: syncMan, adminMan: adminMan, configFilePath: utils.DefaultConfigFilePath,
	}, nil
}

// Start begins the server operations
func (s *Server) Start(profiler, disableMetrics bool, staticPath string, port int) error {

	// Start the sync manager
	if err := s.syncMan.Start(s.configFilePath, port); err != nil {
		return err
	}

	// Anonymously collect usage metrics if not explicitly disabled
	if !disableMetrics {
		go s.RoutineMetrics()
	}

	// Allow cors
	corsObj := utils.CreateCorsObject()

	if s.ssl != nil && s.ssl.Enabled {

		// Setup the handler
		handler := corsObj.Handler(s.routes(profiler, staticPath))
		handler = handlers.HandleMetricMiddleWare(handler, s.metrics)
		handler = s.letsencrypt.LetsEncryptHTTPChallengeHandler(handler)

		// Add existing certificates if any
		if s.ssl.Key != "none" && s.ssl.Crt != "none" {
			logrus.Debugln("Adding existing certificates")
			if err := s.letsencrypt.AddExistingCertificate(s.ssl.Crt, s.ssl.Key); err != nil {
				logrus.Errorf("Could not log existing certificates")
				return err
			}
		}

		go func() {
			// Start the server
			logrus.Info("Starting https server on port: " + strconv.Itoa(port+4))
			httpsServer := &http.Server{Addr: ":" + strconv.Itoa(port+4), Handler: handler, TLSConfig: s.letsencrypt.TLSConfig()}
			if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
				logrus.Fatalln("Error starting https server:", err)
			}
		}()
	}

	// go s.syncMan.StartConnectServer(port, handlers.HandleMetricMiddleWare(corsObj.Handler(s.routerConnect), s.metrics))

	handler := corsObj.Handler(s.routes(profiler, staticPath))
	handler = handlers.HandleMetricMiddleWare(handler, s.metrics)
	handler = s.letsencrypt.LetsEncryptHTTPChallengeHandler(handler)

	logrus.Infoln("Starting http server on port: " + strconv.Itoa(port))

	fmt.Println()
	logrus.Infoln("\t Hosting mission control on http://localhost:" + strconv.Itoa(port) + "/mission-control/")
	fmt.Println()

	logrus.Infoln("Space cloud is running on the specified ports :D")
	return http.ListenAndServe(":"+strconv.Itoa(port), handler)
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
