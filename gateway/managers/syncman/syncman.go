package syncman

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/types"
)

// Manager syncs the project config between folders
type Manager struct {
	lock sync.RWMutex

	// Config related to cluster config
	projectConfig *config.Config

	// Configuration for cluster information
	nodeID        string
	clusterID     string
	advertiseAddr string
	runnerAddr    string
	port          int

	// Configuration for clustering
	storeType string
	store     Store
	services  []*service

	// For authentication
	adminMan AdminSyncmanInterface

	// Modules
	modules types.ModulesInterface
}

type service struct {
	id   string
	addr string
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr string, adminMan *admin.Manager, ssl *config.SSL) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr, storeType: storeType, runnerAddr: runnerAddr, adminMan: adminMan}

	// Initialise the consul client if enabled
	var s Store
	var err error
	switch storeType {
	case "local":
		s, err = NewLocalStore(nodeID, advertiseAddr, ssl)
	case "kube":
		s, err = NewKubeStore(clusterID)
	case "consul":
		s, err = NewConsulStore(nodeID, clusterID, advertiseAddr)
	case "etcd":
		s, err = NewETCDStore(nodeID, clusterID, advertiseAddr)
	default:
		return nil, fmt.Errorf("couldnt initialize syncaman, unknown store type (%v) provided", storeType)
	}

	if err != nil {
		return nil, err
	}
	m.store = s
	m.store.Register()

	return m, nil
}

// Start begins the sync manager operations
func (s *Manager) Start(port int) error {
	// Save the ports
	s.port = port
	// NOTE: SSL is not set in config
	s.projectConfig = &config.Config{}

	adminConfig, err := s.store.GetAdminConfig(context.Background())
	if err != nil {
		return utils.LogError("Unable to fetch initial copy of admin config", "syncman", "Start", err)
	}
	utils.LogDebug("Successfully loaded initial copy of config file", "syncman", "Start", nil)
	_ = s.adminMan.SetConfig(adminConfig)

	// Start routine to observe space cloud projects
	if err := s.store.WatchProjects(func(projects []*config.Project) {
		s.lock.Lock()
		defer s.lock.Unlock()
		utils.LogDebug("Updating projects", "syncman", "Start", map[string]interface{}{"projects": projects})
		s.projectConfig.Projects = projects

		if s.projectConfig.Projects != nil && len(s.projectConfig.Projects) > 0 {
			s.modules.SetProjectConfig(s.projectConfig.Projects[0])
		}
	}); err != nil {
		return err
	}

	// Start routine to admin config
	if err := s.store.WatchAdminConfig(func(clusters []*config.Admin) {
		if len(clusters) == 0 {
			return
		}
		cluster := clusters[0]

		if reflect.DeepEqual(cluster, s.adminMan.GetConfig()) {
			return
		}

		s.lock.Lock()
		s.projectConfig.Admin = cluster
		s.lock.Unlock()

		utils.LogDebug("Updating admin config", "syncman", "Start", map[string]interface{}{"admin config": clusters})
		if err := s.adminMan.SetConfig(cluster); err != nil {
			_ = utils.LogError("Unable to apply admin config", "syncman", "Start", err)
			return
		}

	}); err != nil {
		return err
	}

	// Start routine to observe active space-cloud services
	if err := s.store.WatchServices(func(services scServices) {
		s.lock.Lock()
		defer s.lock.Unlock()
		utils.LogDebug("Updating services", "syncman", "Start", map[string]interface{}{"services": services})

		s.services = services
	}); err != nil {
		return err
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

// SetModules sets all the modules
func (s *Manager) SetModules(modulesInterface types.ModulesInterface) {
	s.modules = modulesInterface
}
