package syncman

import (
	"context"
	"fmt"
	"sync"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
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
	adminMan       AdminSyncmanInterface
	integrationMan integrationInterface

	// Modules
	modules       ModulesInterface
	globalModules GlobalModulesInterface
}

type service struct {
	id   string
	addr string
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr string, adminMan AdminSyncmanInterface, integrationMan integrationInterface, ssl *config.SSL) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr, storeType: storeType, runnerAddr: runnerAddr, adminMan: adminMan, integrationMan: integrationMan}

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
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Cannot initialize syncaman as invalid store type (%v) provided", storeType), nil, nil)
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

	// Fetch initial version of admin config. This must be called before watch admin config callback is invoked
	adminConfig, err := s.store.GetAdminConfig(context.Background())
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to fetch initial copy of admin config", err, map[string]interface{}{})
	}
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Successfully loaded initial copy of config file", map[string]interface{}{})
	s.globalModules.SetMetricsConfig(adminConfig.ClusterConfig.EnableTelemetry)
	if adminConfig.ClusterConfig.LetsEncryptEmail != "" {
		s.modules.LetsEncrypt().SetLetsEncryptEmail(adminConfig.ClusterConfig.LetsEncryptEmail)
	}
	_ = s.adminMan.SetConfig(adminConfig, true)
	_ = s.integrationMan.SetConfig(adminConfig.Integrations)
	s.projectConfig.Admin = adminConfig

	// Start routine to observe active space-cloud services
	if err := s.store.WatchProjects(func(projects []*config.Project) {
		s.lock.Lock()
		defer s.lock.Unlock()
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating projects", map[string]interface{}{"projects": projects})
		for _, p := range s.projectConfig.Projects {
			doesNotExist := true
			for _, q := range projects {
				if p.ID == q.ID {
					doesNotExist = false
					break
				}
			}
			if doesNotExist {
				err := s.store.DeleteProject(context.Background(), p.ID)
				if err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to delete project", err, map[string]interface{}{"project": p.ID})
				}
				s.modules.Delete(p.ID)
			}
		}
		s.projectConfig.Projects = projects

		if s.projectConfig.Projects != nil && len(s.projectConfig.Projects) > 0 {
			for _, p := range s.projectConfig.Projects {
				if err := s.modules.SetProjectConfig(p); err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project config", err, nil)
					break
				}
			}
		}
	}); err != nil {
		return err
	}

	// Start routine to observe space cloud projects
	if err := s.store.WatchAdminConfig(func(clusters []*config.Admin) {
		if len(clusters) == 0 {
			return
		}
		cluster := clusters[0]

		s.lock.Lock()
		s.projectConfig.Admin = cluster
		s.lock.Unlock()

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating admin config", nil)
		if err := s.adminMan.SetConfig(cluster, false); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to apply admin config provided by other space cloud service", err, map[string]interface{}{})
		}

		if err := s.integrationMan.SetConfig(cluster.Integrations); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to apply integration config", err, nil)
		}

		s.globalModules.SetMetricsConfig(cluster.ClusterConfig.EnableTelemetry)
		s.modules.LetsEncrypt().SetLetsEncryptEmail(cluster.ClusterConfig.LetsEncryptEmail)

	}); err != nil {
		return err
	}

	// Start routine to observe space cloud projects
	if err := s.store.WatchServices(func(services scServices) {
		s.lock.Lock()
		defer s.lock.Unlock()
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating services", map[string]interface{}{"services": services})

		s.services = services
	}); err != nil {
		return err
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Exiting syncman start", nil)
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
func (s *Manager) SetModules(modulesInterface ModulesInterface) {
	s.modules = modulesInterface
}

// SetGlobalModules sets all the modules
func (s *Manager) SetGlobalModules(a GlobalModulesInterface) {
	s.globalModules = a
}
