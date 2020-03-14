package syncman

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
)

// Manager syncs the project config between folders
type Manager struct {
	lock sync.RWMutex

	// Config related to cluster config
	projectConfig *config.Config
	configFile    string

	// Configuration for cluster information
	nodeID        string
	clusterID     string
	advertiseAddr string
	runnerAddr    string
	artifactAddr  string
	port          int

	// Configuration for clustering
	storeType string
	store     Store
	services  []*service

	// For authentication
	adminMan *admin.Manager

	// Modules
	modules     model.ModulesInterface
	letsencrypt *letsencrypt.LetsEncrypt
	routing     *routing.Routing
}

type service struct {
	id   string
	addr string
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, artifactAddr string, adminMan *admin.Manager) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr, storeType: storeType, runnerAddr: runnerAddr, adminMan: adminMan, artifactAddr: artifactAddr}

	// Initialise the consul client if enabled
	switch storeType {
	case "none":
		m.services = []*service{{id: nodeID, addr: advertiseAddr}}
		return m, nil
	case "kube":
		s, err := NewKubeStore(clusterID)
		if err != nil {
			return nil, err
		}
		m.store = s
		m.store.Register()
	case "consul":
		s, err := NewConsulStore(nodeID, clusterID, advertiseAddr)
		if err != nil {
			return nil, err
		}
		m.store = s
		m.store.Register()
	case "etcd":
		s, err := NewETCDStore(nodeID, clusterID, advertiseAddr)
		if err != nil {
			return nil, err
		}
		m.store = s
		m.store.Register()
	}

	return m, nil
}

// Start begins the sync manager operations
func (s *Manager) Start(configFilePath string, projectConfig *config.Config, port int) error {
	// Save the ports
	s.lock.Lock()
	defer s.lock.Unlock()

	// Set the callback
	_ = s.modules.SetProjectConfig(projectConfig)
	s.port = port

	s.configFile = configFilePath

	// Write the config to file
	_ = config.StoreConfigToFile(s.projectConfig, s.configFile)

	if len(s.projectConfig.Projects) > 0 {
		_ = s.modules.SetProjectConfig(s.projectConfig)
	}

	if s.storeType != "none" {
		// Start routine to observe active space-cloud services
		if err := s.store.WatchProjects(func(projects []*config.Project) {
			s.lock.Lock()
			defer s.lock.Unlock()

			logrus.WithFields(logrus.Fields{"projects": projects}).Debugln("Updating projects")
			s.projectConfig.Projects = projects
			_ = config.StoreConfigToFile(s.projectConfig, s.configFile)

			if s.projectConfig.Projects != nil && len(s.projectConfig.Projects) > 0 {
				_ = s.modules.SetProjectConfig(s.projectConfig)
			}
		}); err != nil {
			return err
		}

		// Start routine to observe space cloud projects
		if err := s.store.WatchServices(func(services scServices) {
			s.lock.Lock()
			defer s.lock.Unlock()
			logrus.WithFields(logrus.Fields{"services": services}).Debugln("Updating services")

			s.services = services
		}); err != nil {
			return err
		}
	}

	return nil
}

// func (s *Manager) StartConnectServer(port int, handler http.Handler) error {
//	if !s.storeType {
//		return errors.New("consul is not enabled")
//	}
//
//	s.port = port
//
//	// Creating an HTTP server that serves via Connect
//	server := &http.Server{
//		Addr:      ":" + strconv.Itoa(s.port+2),
//		TLSConfig: s.consulService.ServerTLSConfig(),
//		Handler:   handler,
//	}
//
//	fmt.Println("Starting https server (consul connect) on port: " + strconv.Itoa(s.port+2))
//	return server.ListenAndServeTLS("", "")
// }

// SetGlobalConfig sets the global config. This must be called before the Start command.
func (s *Manager) SetGlobalConfig(c *config.Config) {
	s.lock.Lock()
	s.projectConfig = c
	s.lock.Unlock()
}

// GetGlobalConfig gets the global config
func (s *Manager) GetGlobalConfig() *config.Config {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.projectConfig
}

// SetModules sets all the modules
func (s *Manager) SetModules(modulesInterface model.ModulesInterface, letsEncrypt *letsencrypt.LetsEncrypt, routing *routing.Routing) {
	s.modules = modulesInterface
	s.letsencrypt = letsEncrypt
	s.routing = routing
}
