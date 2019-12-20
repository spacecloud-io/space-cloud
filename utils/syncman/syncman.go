package syncman

import (
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

// Manager syncs the project config between folders
type Manager struct {
	lock sync.RWMutex

	// Config related to cluster config
	projectConfig *config.Config
	configFile    string
	cb            func(*config.Config) error

	// Configuration for cluster information
	nodeID        string
	clusterID     string
	advertiseAddr string
	port          int

	// Configuration for clustering
	storeType string
	store     Store
	services  []*service
}

type service struct {
	id   string
	addr string
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID, storeAddr, advertiseAddr, storeType string) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr, storeType: storeType}

	// Initialise the consul client if enabled
	switch storeType {
	case "none":
		return m, nil
	case "consul":
		s, err := NewConsulStore(nodeID, clusterID, advertiseAddr)
		if err != nil {
			return nil, err
		}
		m.store = s
		m.store.Register()
	case "etcd":
		s, err := NewETCDStore(nodeID, clusterID, advertiseAddr, storeAddr)
		if err != nil {
			return nil, err
		}
		m.store = s
		m.store.Register()
	}

	return m, nil
}

// Start begins the sync manager operations
func (s *Manager) Start(configFilePath string, cb func(*config.Config) error) error {
	// Save the ports
	s.lock.Lock()
	defer s.lock.Unlock()

	// Set the callback
	s.cb = cb

	s.configFile = configFilePath

	// Write the config to file
	config.StoreConfigToFile(s.projectConfig, s.configFile)

	if len(s.projectConfig.Projects) > 0 {
		cb(s.projectConfig)
	}

	if s.storeType != "none" {
		// Start routine to observe active space-cloud services
		if err := s.store.WatchProjects(func(projects []*config.Project) {
			s.lock.Lock()
			defer s.lock.Unlock()

			s.projectConfig.Projects = projects
			config.StoreConfigToFile(s.projectConfig, s.configFile)

			if s.projectConfig.Projects != nil && len(s.projectConfig.Projects) > 0 {
				s.cb(s.projectConfig)
			}
		}); err != nil {
			return err
		}

		// Start routine to observe space cloud projects
		if err := s.store.WatchServices(func(services scServices) {
			s.lock.Lock()
			defer s.lock.Unlock()

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
