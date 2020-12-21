package syncman

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/leader"
	"github.com/spaceuptech/space-cloud/gateway/utils/pubsub"
)

const pubSubOperationRenew = "renew"
const pubSubOperationUpgrade = "upgrade"

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

	leader       *leader.Module
	pubsubClient *pubsub.Module
	// Configuration for clustering
	storeType string
	store     Store
	services  model.ScServices

	// For authentication
	adminMan       AdminSyncmanInterface
	integrationMan integrationInterface

	// Modules
	modules       ModulesInterface
	globalModules GlobalModulesInterface
}

// New creates a new instance of the sync manager
func New(nodeID, clusterID, storeType, runnerAddr string, adminMan AdminSyncmanInterface, integrationMan integrationInterface, ssl *config.SSL) (*Manager, error) {

	// Create a new manager instance
	m := &Manager{nodeID: nodeID, clusterID: clusterID, storeType: storeType, runnerAddr: runnerAddr, adminMan: adminMan, integrationMan: integrationMan}

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

	pubsubClient, err := pubsub.New("license-manager", os.Getenv("REDIS_CONN"))
	if err != nil {
		return nil, helpers.Logger.LogError("syncman-new", "Unable to initialize pub sub client required for sync module, ensure that redis database is running", err, nil)
	}
	m.pubsubClient = pubsubClient
	m.leader = leader.New(nodeID, pubsubClient)

	if err := m.SetPubSubRoutines(nodeID); err != nil {
		return nil, err
	}

	return m, nil
}

// Start begins the sync manager operations
func (s *Manager) Start(port int) error {

	// Start routine to observe space cloud projects
	if err := s.store.WatchServices(func(eventType, serviceID string, services model.ScServices) {
		s.lockServices.Lock()
		defer s.lockServices.Unlock()
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating services", map[string]interface{}{"services": services, "eventType": eventType, "id": serviceID})

		s.adminMan.SetServices(eventType, services)
		s.services = services
	}); err != nil {
		return err
	}

	count := 1
	maxRetryCount := 6
	incrementBy := 10
	for len(s.services) == 0 {
		helpers.Logger.LogDebug("syncman-start", fmt.Sprintf("Waiting for gateway services to register - retry count (%d)", count), nil)
		time.Sleep(time.Duration(count*incrementBy) * time.Second)
		if count == maxRetryCount {
			return helpers.Logger.LogError("syncman-start", "Cannot start gateway, gateway service not registered", nil, nil)
		}
		count++
	}

	// Save the ports
	s.port = port
	// NOTE: SSL is not set in config
	s.projectConfig = &config.Config{}

	// Set global config
	globalConfig, err := s.store.GetGlobalConfig()
	if err != nil {
		return err
	}

	_ = s.adminMan.SetConfig(globalConfig.License)
	s.adminMan.SetServices(config.ResourceAddEvent, s.services)
	s.adminMan.SetIntegrationConfig(globalConfig.Integrations)
	_ = s.integrationMan.SetConfig(globalConfig.Integrations, globalConfig.IntegrationHooks)

	s.leader.AddCallBack("admin-set-service", func() {
		s.lockServices.RLock()
		s.adminMan.SetServices(config.ResourceDeleteEvent, s.services)
		s.lockServices.RUnlock()
	})

	// Set caching config
	if err := s.modules.Caching().SetCachingConfig(context.TODO(), globalConfig.CacheConfig); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set config of caching module, ensure redis instance is running", err, nil)
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
			if eventType == config.ResourceDeleteEvent {
				s.modules.Delete(projectID)
				return
			}
			_ = s.modules.SetProjectConfig(ctx, s.projectConfig.Projects[projectID].ProjectConfig)

		case config.ResourceAuthProvider:
			_ = s.modules.SetUsermanConfig(ctx, projectID, s.projectConfig.Projects[projectID].Auths)

		case config.ResourceDatabaseConfig:
			p := s.projectConfig.Projects[projectID]
			_ = s.modules.SetDatabaseConfig(ctx, projectID, p.DatabaseConfigs, p.DatabaseSchemas, p.DatabaseRules, p.DatabasePreparedQueries)

		case config.ResourceDatabaseSchema:
			_ = s.modules.SetDatabaseSchemaConfig(ctx, projectID, s.projectConfig.Projects[projectID].DatabaseSchemas)

		case config.ResourceDatabaseRule:
			_ = s.modules.SetDatabaseRulesConfig(ctx, projectID, s.projectConfig.Projects[projectID].DatabaseRules)

		case config.ResourceDatabasePreparedQuery:
			_ = s.modules.SetDatabasePreparedQueryConfig(ctx, projectID, s.projectConfig.Projects[projectID].DatabasePreparedQueries)

		case config.ResourceEventingConfig:
			p := s.projectConfig.Projects[projectID]
			_ = s.modules.SetEventingConfig(ctx, projectID, p.EventingConfig, p.EventingRules, p.EventingSchemas, p.EventingTriggers)

		case config.ResourceEventingSchema:
			_ = s.modules.SetEventingSchemaConfig(ctx, projectID, s.projectConfig.Projects[projectID].EventingSchemas)

		case config.ResourceEventingRule:
			_ = s.modules.SetEventingRuleConfig(ctx, projectID, s.projectConfig.Projects[projectID].EventingRules)

		case config.ResourceEventingTrigger:
			_ = s.modules.SetEventingTriggerConfig(ctx, projectID, s.projectConfig.Projects[projectID].EventingTriggers)

		case config.ResourceFileStoreConfig:
			_ = s.modules.SetFileStoreConfig(ctx, projectID, s.projectConfig.Projects[projectID].FileStoreConfig)

		case config.ResourceFileStoreRule:
			_ = s.modules.SetFileStoreSecurityRuleConfig(ctx, projectID, s.projectConfig.Projects[projectID].FileStoreRules)

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
			if err := s.integrationMan.SetIntegrations(s.projectConfig.Integrations); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to apply integration config", err, nil)
				return
			}
			s.adminMan.SetIntegrationConfig(s.projectConfig.Integrations)

		case config.ResourceIntegrationHook:
			s.integrationMan.SetIntegrationHooks(s.projectConfig.IntegrationHooks)

		case config.ResourceLicense:
			if err := s.adminMan.SetConfig(s.projectConfig.License); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to apply admin config provided by other space cloud service", err, map[string]interface{}{})
				return
			}
		case config.ResourceCacheConfig:
			if err := s.modules.Caching().SetCachingConfig(ctx, s.projectConfig.CacheConfig); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to apply admin config provided by other space cloud service", err, map[string]interface{}{})
				return
			}
		default:
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unknown resource type provided", err, map[string]interface{}{"resourceType": resourceType})
			return
		}
	}); err != nil {
		return err
	}
	s.store.WatchLicense(func(eventType, resourceID string, resourceType config.Resource, resource *config.License) {

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Updating license", map[string]interface{}{"event": eventType, "resourceId": resourceID, "resource": resource, "resourceType": resourceType})

		if resourceType == config.ResourceLicense {
			if err := s.adminMan.SetConfig(resource); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to apply admin config provided by other space cloud service", err, map[string]interface{}{})
				return
			}
			s.lock.Lock()
			s.projectConfig.License = resource
			s.lock.Unlock()
		}

	})

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
