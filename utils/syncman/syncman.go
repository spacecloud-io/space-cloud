package syncman

import (
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

// Manager syncs the project config between folders
type Manager struct {
	lock          sync.RWMutex
	projectConfig *config.Config
	configFile    string
	cb            func(*config.Config) error
}

// New creates a new instance of the sync manager
func New() *Manager {
	// Create a SyncManger instance
	return &Manager{}
}

// Start begins the sync manager operations
func (s *Manager) Start(nodeID, configFilePath string, cb func(*config.Config) error) error {
	// Save the ports
	s.lock.Lock()
	defer s.lock.Unlock()

	// Set the callback
	s.cb = cb

	s.configFile = configFilePath
	if s.projectConfig.NodeID == "" {
		s.projectConfig.NodeID = nodeID
	}

	// Write the config to file
	config.StoreConfigToFile(s.projectConfig, s.configFile)

	if len(s.projectConfig.Projects) > 0 {
		cb(s.projectConfig)
	}

	return nil
}

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
