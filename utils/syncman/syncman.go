package syncman

import (
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"

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
	nodeID    string
	clusterID string
	port      int

	// Configuration for clustering
	isConsulEnabled bool
	consulClient    *api.Client
	services        []*api.ServiceEntry
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID string, isConsulEnabled bool) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID}

	// Initialise the consul client if enabled
	if isConsulEnabled {
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			return nil, err
		}

		m.isConsulEnabled = true
		m.consulClient = client

		// Set the node id if not already exists
		if m.nodeID == "" || strings.HasPrefix(m.nodeID, "auto") {
			info, err := client.Agent().Self()
			if err != nil {
				return nil, err
			}

			m.nodeID = info["Config"]["NodeID"].(string)
		}

		return m, nil
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

	if s.isConsulEnabled {
		// Start routine to observe active space-cloud services
		if err := s.watchService(); err != nil {
			return err
		}

		// Start routine to observe space cloud projects
		if err := s.watchProjects(); err != nil {
			return err
		}
	}

	return nil
}

//func (s *Manager) StartConnectServer(port int, handler http.Handler) error {
//	if !s.isConsulEnabled {
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
//}

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
