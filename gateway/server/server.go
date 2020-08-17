package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/modules/global"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	nodeID   string
	ssl      *config.SSL
	modules  *modules.Modules
	managers *managers.Managers
}

// New creates a new server instance
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr string, isDev bool, adminUserInfo *config.AdminUser, ssl *config.SSL) (*Server, error) {

	managers, err := managers.New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, isDev, adminUserInfo, ssl)
	if err != nil {
		return nil, err
	}

	globalMods, err := global.New(clusterID, nodeID, isDev, managers)
	if err != nil {
		return nil, err
	}

	modules, err := modules.New(nodeID, managers, globalMods)
	if err != nil {
		return nil, err
	}

	managers.Sync().SetModules(modules)
	managers.Sync().SetGlobalModules(globalMods.Metrics())

	logrus.Infoln("Creating a new server with id", nodeID)

	return &Server{nodeID: nodeID, managers: managers, modules: modules, ssl: ssl}, nil
}

// Start begins the server operations
func (s *Server) Start(profiler bool, staticPath string, port int, restrictedHosts []string) error {
	// Start the sync manager
	if err := s.managers.Sync().Start(port); err != nil {
		return err
	}

	// Allow cors
	corsObj := utils.CreateCorsObject()

	if s.ssl != nil && s.ssl.Enabled {

		// Setup the handler
		handler := corsObj.Handler(s.routes(profiler, staticPath, restrictedHosts))
		handler = s.modules.LetsEncrypt().LetsEncryptHTTPChallengeHandler(handler)

		// Add existing certificates if any
		if s.ssl.Key != "none" && s.ssl.Crt != "none" {
			logrus.Debugln("Adding existing certificates")
			if err := s.modules.LetsEncrypt().AddExistingCertificate(s.ssl.Crt, s.ssl.Key); err != nil {
				logrus.Errorf("Could not log existing certificates")
				return err
			}
		}

		go func() {
			// Start the server
			logrus.Info("Starting https server on port: " + strconv.Itoa(port+4))
			httpsServer := &http.Server{Addr: ":" + strconv.Itoa(port+4), Handler: handler, TLSConfig: s.modules.LetsEncrypt().TLSConfig()}
			if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
				logrus.Fatalln("Error starting https server:", err)
			}
		}()
	}

	handler := corsObj.Handler(s.routes(profiler, staticPath, restrictedHosts))
	handler = s.modules.LetsEncrypt().LetsEncryptHTTPChallengeHandler(handler)

	logrus.Infoln("Starting http server on port: " + strconv.Itoa(port))

	fmt.Println()
	logrus.Infoln("\t Hosting mission control on http://localhost:" + strconv.Itoa(port) + "/mission-control/")
	fmt.Println()

	logrus.Infoln("Space cloud is running on the specified ports :D")
	return http.ListenAndServe(":"+strconv.Itoa(port), handler)
}
