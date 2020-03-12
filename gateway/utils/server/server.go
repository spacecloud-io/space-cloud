package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/handlers"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	nodeID         string
	configFilePath string
	adminMan       *admin.Manager
	syncMan        *syncman.Manager
	letsencrypt    *letsencrypt.LetsEncrypt
	routing        *routing.Routing
	metrics        *metrics.Module
	ssl            *config.SSL
	modules        *modules.Modules
}

// New creates a new server instance
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, artifactAddr string, removeProjectScope bool, metricsConfig *metrics.Config) (*Server, error) {

	// Create the fundamental modules
	m, err := metrics.New(nodeID, metricsConfig)
	if err != nil {
		return nil, err
	}

	adminMan := admin.New()
	syncMan, err := syncman.New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, artifactAddr, adminMan)
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

	modules, err := modules.New(nodeID, removeProjectScope, syncMan, adminMan, m)
	if err != nil {
		return nil, err
	}

	syncMan.SetModules(modules, le, r)

	logrus.Infoln("Creating a new server with id", nodeID)

	return &Server{nodeID: nodeID, syncMan: syncMan, adminMan: adminMan, letsencrypt: le, routing: r, metrics: m, configFilePath: utils.DefaultConfigFilePath, modules: modules}, nil
}

// Start begins the server operations
func (s *Server) Start(profiler, disableMetrics bool, staticPath string, port int, restrictedHosts []string) error {

	// Start the sync manager
	if err := s.syncMan.Start(s.configFilePath, s.LoadConfig, port); err != nil {
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
		handler := corsObj.Handler(s.routes(profiler, staticPath, restrictedHosts))
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

	handler := corsObj.Handler(s.routes(profiler, staticPath, restrictedHosts))
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

// LoadConfig configures each module to to use the provided config
func (s *Server) LoadConfig(config *config.Config) error {

	if config.Projects != nil {

		p := config.Projects[0]

		// Always set the config of the crud module first
		// Set the configuration for the crud module
		if err := s.modules.Crud.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorln("Error in crud module config: ", err)
			return err
		}

		if err := s.modules.Schema.SetConfig(p.Modules.Crud, p.ID); err != nil {
			logrus.Errorln("Error in schema module config: ", err)
			return err
		}

		// Set the configuration for the auth module
		if err := s.modules.Auth.SetConfig(p.ID, p.Secret, p.AESkey, p.Modules.Crud, p.Modules.FileStore, p.Modules.Services, &p.Modules.Eventing); err != nil {
			logrus.Errorln("Error in auth module config: ", err)
			return err
		}

		// Set the configuration for the functions module
		s.modules.Functions.SetConfig(p.ID, p.Modules.Services)

		// Set the configuration for the user management module
		s.modules.User.SetConfig(p.Modules.Auth)

		// Set the configuration for the file storage module
		if err := s.modules.File.SetConfig(p.Modules.FileStore); err != nil {
			logrus.Errorln("Error in files module config: ", err)
			return err
		}

		if err := s.modules.Eventing.SetConfig(p.ID, &p.Modules.Eventing); err != nil {
			logrus.Errorln("Error in eventing module config: ", err)
			return err
		}

		// Set the configuration for the realtime module
		if err := s.modules.Realtime.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorln("Error in realtime module config: ", err)
			return err
		}

		// Set the configuration for the graphql module
		s.modules.Graphql.SetConfig(p.ID)

		// Set the configuration for the letsencrypt module
		if err := s.letsencrypt.SetProjectDomains(p.ID, p.Modules.LetsEncrypt); err != nil {
			logrus.Errorln("Error in letsencrypt module config: ", err)
			return err
		}

		// Set the configuration for the routing module
		s.routing.SetProjectRoutes(p.ID, p.Modules.Routes)
	}

	return nil
}

// SetConfigFilePath sets the config file path
func (s *Server) SetConfigFilePath(configFilePath string) {
	s.configFilePath = configFilePath
}
