package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/spaceuptech/helpers"

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

	helpers.Logger.LogInfo(helpers.GetRequestID(nil), fmt.Sprintf("Creating a new server with id %s", nodeID), nil)

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
		handler := corsObj.Handler(loggerMiddleWare(s.routes(profiler, staticPath, restrictedHosts)))
		handler = s.modules.LetsEncrypt().LetsEncryptHTTPChallengeHandler(handler)

		// Add existing certificates if any
		if s.ssl.Key != "none" && s.ssl.Crt != "none" {
			helpers.Logger.LogDebug(helpers.GetRequestID(nil), "Adding existing certificates", nil)
			if err := s.modules.LetsEncrypt().AddExistingCertificate(s.ssl.Crt, s.ssl.Key); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(nil), "Could not log existing certificates", err, nil)
			}
		}

		go func() {
			// Start the server
			helpers.Logger.LogInfo(helpers.GetRequestID(nil), "Starting https server on port: "+strconv.Itoa(port+4), nil)
			httpsServer := &http.Server{Addr: ":" + strconv.Itoa(port+4), Handler: handler, TLSConfig: s.modules.LetsEncrypt().TLSConfig()}
			if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
				log.Fatalln("Error starting https server:", err)
			}
		}()
	}

	handler := corsObj.Handler(loggerMiddleWare(s.routes(profiler, staticPath, restrictedHosts)))
	handler = s.modules.LetsEncrypt().LetsEncryptHTTPChallengeHandler(handler)

	helpers.Logger.LogInfo(helpers.GetRequestID(nil), "Starting http server on port: "+strconv.Itoa(port), nil)

	if staticPath != "" {
		helpers.Logger.LogInfo(helpers.GetRequestID(nil), "Hosting mission control on http://localhost:"+strconv.Itoa(port)+"/mission-control/", nil)
	}

	helpers.Logger.LogInfo(helpers.GetRequestID(nil), fmt.Sprintf("Space cloud is running on the specified ports :%v", port), nil)
	return http.ListenAndServe(":"+strconv.Itoa(port), handler)
}
