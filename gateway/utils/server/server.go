package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore"
	"github.com/spaceuptech/space-cloud/gateway/modules/functions"
	"github.com/spaceuptech/space-cloud/gateway/modules/realtime"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/modules/userman"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
	"github.com/spaceuptech/space-cloud/gateway/utils/handlers"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Server is the object which sets up the server and handles all server operations
type Server struct {
	nodeID         string
	auth           *auth.Module
	crud           *crud.Module
	user           *userman.Module
	file           *filestore.Module
	functions      *functions.Module
	realtime       *realtime.Module
	eventing       *eventing.Module
	configFilePath string
	adminMan       *admin.Manager
	syncMan        *syncman.Manager
	letsencrypt    *letsencrypt.LetsEncrypt
	routing        *routing.Routing
	metrics        *metrics.Module
	ssl            *config.SSL
	graphql        *graphql.Module
	schema         *schema.Schema
}

// New creates a new server instance
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, artifactAddr string, removeProjectScope bool, metricsConfig *metrics.Config) (*Server, error) {

	// Create the fundamental modules
	c := crud.Init(removeProjectScope)

	m, err := metrics.New(nodeID, metricsConfig)
	if err != nil {
		return nil, err
	}

	adminMan := admin.New()
	syncMan, err := syncman.New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, artifactAddr, adminMan)
	if err != nil {
		return nil, err
	}

	s := schema.Init(c, removeProjectScope)
	a := auth.Init(nodeID, c, s, removeProjectScope)
	a.SetMakeHTTPRequest(syncMan.MakeHTTPRequest)

	fn := functions.Init(a, syncMan)

	f := filestore.Init(a)

	// Initialise the eventing module and set the crud module hooks
	e := eventing.New(a, c, s, fn, adminMan, syncMan, f)
	f.SetEventingModule(e)

	c.SetHooks(&model.CrudHooks{
		Create: e.HookDBCreateIntent,
		Update: e.HookDBUpdateIntent,
		Delete: e.HookDBDeleteIntent,
		Batch:  e.HookDBBatchIntent,
		Stage:  e.HookStage,
	}, m.AddDBOperation)

	rt, err := realtime.Init(nodeID, e, a, c, s, m, syncMan)
	if err != nil {
		return nil, err
	}

	u := userman.Init(c, a)
	graphqlMan := graphql.New(a, c, fn, s)

	// Initialise a lets encrypt client
	le, err := letsencrypt.New()
	if err != nil {
		return nil, err
	}

	// Initialise the routing module
	r := routing.New()

	logrus.Infoln("Creating a new server with id", nodeID)

	return &Server{nodeID: nodeID, auth: a, crud: c,
		user: u, file: f, syncMan: syncMan, adminMan: adminMan, letsencrypt: le, routing: r, metrics: m,
		functions: fn, realtime: rt, configFilePath: utils.DefaultConfigFilePath,
		eventing: e, graphql: graphqlMan, schema: s}, nil
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
		if err := s.crud.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorln("Error in crud module config: ", err)
			return err
		}

		if err := s.schema.SetConfig(p.Modules.Crud, p.ID); err != nil {
			logrus.Errorln("Error in schema module config: ", err)
			return err
		}

		// Set the configuration for the auth module
		if err := s.auth.SetConfig(p.ID, p.Secret, p.AESkey, p.Modules.Crud, p.Modules.FileStore, p.Modules.Services, &p.Modules.Eventing); err != nil {
			logrus.Errorln("Error in auth module config: ", err)
			return err
		}

		// Set the configuration for the functions module
		s.functions.SetConfig(p.ID, p.Modules.Services)

		// Set the configuration for the user management module
		s.user.SetConfig(p.Modules.Auth)

		// Set the configuration for the file storage module
		if err := s.file.SetConfig(p.Modules.FileStore); err != nil {
			logrus.Errorln("Error in files module config: ", err)
			return err
		}

		if err := s.eventing.SetConfig(p.ID, &p.Modules.Eventing); err != nil {
			logrus.Errorln("Error in eventing module config: ", err)
			return err
		}

		// Set the configuration for the realtime module
		if err := s.realtime.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorln("Error in realtime module config: ", err)
			return err
		}

		// Set the configuration for the graphql module
		s.graphql.SetConfig(p.ID)

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
