package syncman

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
	"github.com/spaceuptech/space-cloud/gateway/utils/types"
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
	port          int

	// Configuration for clustering
	storeType string
	store     Store
	services  []*service

	// For authentication
	adminMan *admin.Manager

	// Modules
	modules     types.ModulesInterface
	letsencrypt *letsencrypt.LetsEncrypt
	routing     *routing.Routing
}

type service struct {
	id   string
	addr string
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, configFile string, adminMan *admin.Manager) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr, storeType: storeType, runnerAddr: runnerAddr, configFile: configFile, adminMan: adminMan}

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
func (s *Manager) Start(projectConfig *config.Config, port int) error {
	// Save the ports

	s.port = port

	// Write the config to file
	_ = config.StoreConfigToFile(s.projectConfig, s.configFile)

	if len(s.projectConfig.Projects) > 0 {
		for _, p := range s.projectConfig.Projects {
			if err := s.modules.SetProjectConfig(p, s.letsencrypt, s.routing); err != nil {
				logrus.Errorf("Unable to apply project (%s). Upgrade your plan.", p.ID)
				return err
			}
		}
	}

	if s.storeType != "none" {
		// Fetch initial version of admin config. This must be called before watch admin config callback is invoked
		adminConfig, err := s.store.GetAdminConfig(context.Background())
		if err != nil {
			return utils.LogError("Unable to fetch initial copy of admin config", "syncman", "Start", err)
		}
		utils.LogDebug("Successfully loaded initial copy of config file", "syncman", "Start", nil)
		_ = s.adminMan.SetConfig(adminConfig, true)

		// Now lets store the config as well
		if s.checkIfLeaderGateway(s.nodeID) {
			s.projectConfig.Admin = s.adminMan.GetConfig()
			if err := s.store.SetAdminConfig(context.Background(), s.projectConfig.Admin); err != nil {
				return utils.LogError("Unable to save initial license copy", "syncman", "Start", err)
			}
			utils.LogDebug("Successfully stored initial copy of config file", "syncman", "Start", nil)
		}

		// Start routine to observe active space-cloud services
		if err := s.store.WatchProjects(func(projects []*config.Project) {
			s.lock.Lock()
			defer s.lock.Unlock()

			logrus.WithFields(logrus.Fields{"projects": projects}).Debugln("Updating projects")
			s.projectConfig.Projects = projects
			_ = config.StoreConfigToFile(s.projectConfig, s.configFile)

			if s.projectConfig.Projects != nil && len(s.projectConfig.Projects) > 0 {
				for _, p := range s.projectConfig.Projects {
					if err := s.modules.SetProjectConfig(p, s.letsencrypt, s.routing); err != nil {
						_ = utils.LogError("Unable to set project config", "syncman", "watch-projects", err)
						break
					}
				}
			}
		}); err != nil {
			return err
		}

		// Start routine to observe space cloud projects
		if err := s.store.WatchServices(func(services scServices) {
			s.lock.Lock()
			defer s.lock.Unlock()

			data, _ := json.Marshal(services)
			logrus.WithFields(logrus.Fields{"services": string(data)}).Debugln("Updating services")

			s.services = services
		}); err != nil {
			return err
		}

		// Start routine to observe space cloud projects
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
			_ = config.StoreConfigToFile(s.projectConfig, s.configFile)
			s.lock.Unlock()

			logrus.WithFields(logrus.Fields{"admin config": clusters}).Debugln("Updating admin config")
			if err := s.adminMan.SetConfig(cluster, false); err != nil {
				_ = utils.LogError("Unable to apply admin config", "syncman", "Start", err)
				return
			}
		}); err != nil {
			return err
		}
	}

	utils.LogDebug("Exiting syncman start", "syncman", "Start", nil)
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
func (s *Manager) SetModules(modulesInterface types.ModulesInterface, letsEncrypt *letsencrypt.LetsEncrypt, routing *routing.Routing) {
	s.modules = modulesInterface
	s.letsencrypt = letsEncrypt
	s.routing = routing
}
