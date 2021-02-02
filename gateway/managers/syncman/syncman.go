package syncman

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
)

// Manager syncs the project config between folders
type Manager struct {
	lock         sync.RWMutex
	lockServices sync.RWMutex

	// Config related to cluster config
	projectConfig *config.Config

	// Configuration for cluster information
	nodeID     string
	clusterID  string
	runnerAddr string
	port       int

	// Configuration for clustering
	storeType string
	store     Store
	services  []*service

	// For authentication
	adminMan AdminSyncmanInterface

	// Modules
	modules       ModulesInterface
	globalModules GlobalModulesInterface
}

type service struct {
	id string
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID, storeType, runnerAddr string, adminMan *admin.Manager, ssl *config.SSL) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID, storeType: storeType, runnerAddr: runnerAddr, adminMan: adminMan}

	// Initialise the consul client if enabled
	var s Store
	var err error
	switch storeType {
	case "local":
		s, err = NewLocalStore(nodeID, ssl)
	case "kube":
		s, err = NewKubeStore(clusterID)
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

	// Set global config
	globalConfig, err := s.store.GetGlobalConfig()
	if err != nil {
		return err
	}

	// Set metric config
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Successfully loaded initial copy of config file", map[string]interface{}{})
	s.globalModules.SetMetricsConfig(globalConfig.ClusterConfig.EnableTelemetry)

	// Set letsencrypt config
	if globalConfig.ClusterConfig.LetsEncryptEmail != "" {
		s.modules.LetsEncrypt().SetLetsEncryptEmail(globalConfig.ClusterConfig.LetsEncryptEmail)
	}

	s.projectConfig = globalConfig

	// Set initial project config
	if err := s.modules.SetInitialProjectConfig(context.TODO(), globalConfig.Projects); err != nil {
		return err
	}

	// Start routine to observe space cloud project level resources
	if err := s.store.WatchResources(func(eventType, resourceID string, resourceType config.Resource, resource interface{}) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		isSkip, err := s.validateResource(ctx, eventType, resourceID, resourceType, resource)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to update resources", err, nil)
			return
		}
		if isSkip {
			helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Found duplicate resource, skipping the resource", map[string]interface{}{"event": eventType, "resourceId": resourceID, "resource": resource, "resourceType": resourceType})
			return
		}

		_, projectID, _, err := splitResourceID(ctx, resourceID)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to split resource id in watch resources", err, nil)
			return
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating resources", map[string]interface{}{"event": eventType, "resourceId": resourceID, "resource": resource, "projectId": projectID, "resourceType": resourceType})

		s.lock.Lock()
		defer s.lock.Unlock()
		if err := updateResource(ctx, eventType, s.projectConfig, resourceID, resourceType, resource); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to update resources", err, nil)
			return
		}

		switch resourceType {
		case config.ResourceProject:
			_ = s.modules.SetProjectConfig(ctx, s.projectConfig.Projects[projectID].ProjectConfig)

		case config.ResourceAuthProvider:
			_ = s.modules.SetUsermanConfig(ctx, projectID, s.projectConfig.Projects[projectID].Auths)

		case config.ResourceSecurityFunction:
			_ = s.modules.SetSecurityFunctionConfig(ctx, projectID, s.projectConfig.Projects[projectID].SecurityFunctions)

		case config.ResourceDatabaseConfig:
			p := s.projectConfig.Projects[projectID]
			_ = s.modules.SetDatabaseConfig(ctx, projectID, p.DatabaseConfigs, p.DatabaseSchemas, p.DatabaseRules, p.DatabasePreparedQueries)

		case config.ResourceDatabaseSchema:
			_ = s.modules.SetDatabaseSchemaConfig(ctx, projectID, s.projectConfig.Projects[projectID].DatabaseSchemas)

		case config.ResourceDatabaseRule:
			_ = s.modules.SetDatabaseRulesConfig(ctx, s.projectConfig.Projects[projectID].DatabaseRules)

		case config.ResourceDatabasePreparedQuery:
			_ = s.modules.SetDatabasePreparedQueryConfig(ctx, s.projectConfig.Projects[projectID].DatabasePreparedQueries)

		case config.ResourceEventingConfig:
			p := s.projectConfig.Projects[projectID]
			_ = s.modules.SetEventingConfig(ctx, projectID, p.EventingConfig, p.EventingRules, p.EventingSchemas, p.EventingTriggers)

		case config.ResourceEventingSchema:
			_ = s.modules.SetEventingSchemaConfig(ctx, s.projectConfig.Projects[projectID].EventingSchemas)

		case config.ResourceEventingRule:
			_ = s.modules.SetEventingRuleConfig(ctx, s.projectConfig.Projects[projectID].EventingRules)

		case config.ResourceEventingTrigger:
			_ = s.modules.SetEventingTriggerConfig(ctx, s.projectConfig.Projects[projectID].EventingTriggers)

		case config.ResourceFileStoreConfig:
			_ = s.modules.SetFileStoreConfig(ctx, projectID, s.projectConfig.Projects[projectID].FileStoreConfig)

		case config.ResourceFileStoreRule:
			s.modules.SetFileStoreSecurityRuleConfig(ctx, projectID, s.projectConfig.Projects[projectID].FileStoreRules)

		case config.ResourceProjectLetsEncrypt:
			_ = s.modules.SetLetsencryptConfig(ctx, projectID, s.projectConfig.Projects[projectID].LetsEncrypt)

		case config.ResourceIngressRoute:
			_ = s.modules.SetIngressRouteConfig(ctx, projectID, s.projectConfig.Projects[projectID].IngressRoutes)

		case config.ResourceIngressGlobal:
			_ = s.modules.SetIngressGlobalRouteConfig(ctx, projectID, s.projectConfig.Projects[projectID].IngressGlobal)

		case config.ResourceRemoteService:
			_ = s.modules.SetRemoteServiceConfig(ctx, projectID, s.projectConfig.Projects[projectID].RemoteService)

		case config.ResourceCluster:
			s.globalModules.SetMetricsConfig(s.projectConfig.ClusterConfig.EnableTelemetry)
			s.modules.LetsEncrypt().SetLetsEncryptEmail(s.projectConfig.ClusterConfig.LetsEncryptEmail)

		case config.ResourceIntegration:
		case config.ResourceIntegrationHook:

		default:
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unknown resource type provided", err, map[string]interface{}{"resourceType": resourceType})
			return
		}
	}); err != nil {
		return err
	}

	// Start routine to observe active space-cloud services
	if err := s.store.WatchServices(func(services scServices) {
		s.lockServices.Lock()
		defer s.lockServices.Unlock()
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating services", map[string]interface{}{"services": services})

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
func (s *Manager) SetModules(modulesInterface ModulesInterface) {
	s.modules = modulesInterface
}

// SetGlobalModules sets all the modules
func (s *Manager) SetGlobalModules(a GlobalModulesInterface) {
	s.globalModules = a
}
