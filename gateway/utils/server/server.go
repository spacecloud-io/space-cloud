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
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr string, disableMetrics, isDev bool, adminUserInfo *config.AdminUser, ssl *config.SSL) (*Server, error) {

	// Create the fundamental modules
	adminMan := admin.New(nodeID, clusterID, isDev, adminUserInfo)
	syncMan, err := syncman.New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, adminMan, ssl)
	if err != nil {
		return nil, err
	}
	m, err := metrics.New(clusterID, nodeID, disableMetrics, adminMan, syncMan, !isDev)
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

	modules, err := modules.New(nodeID, syncMan, adminMan, m)
	if err != nil {
		return nil, err
	}

	syncMan.SetModules(modules, le, r)

	logrus.Infoln("Creating a new server with id", nodeID)

	return &Server{nodeID: nodeID, syncMan: syncMan, adminMan: adminMan, letsencrypt: le, routing: r, metrics: m, configFilePath: utils.DefaultConfigFilePath, modules: modules, ssl: ssl}, nil
}

// Start begins the server operations
func (s *Server) Start(profiler bool, staticPath string, port int, restrictedHosts []string) error {
	// Start the sync manager
	if err := s.syncMan.Start(port); err != nil {
		return err
	}

	// Allow cors
	corsObj := utils.CreateCorsObject()

	if s.ssl != nil && s.ssl.Enabled {

		// Setup the handler
		handler := corsObj.Handler(s.routes(profiler, staticPath, restrictedHosts))
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

	handler := corsObj.Handler(s.routes(profiler, staticPath, restrictedHosts))
	handler = s.letsencrypt.LetsEncryptHTTPChallengeHandler(handler)

	logrus.Infoln("Starting http server on port: " + strconv.Itoa(port))

	fmt.Println()
	logrus.Infoln("\t Hosting mission control on http://localhost:" + strconv.Itoa(port) + "/mission-control/")
	fmt.Println()

	logrus.Infoln("Space cloud is running on the specified ports :D")
	return http.ListenAndServe(":"+strconv.Itoa(port), handler)
}
