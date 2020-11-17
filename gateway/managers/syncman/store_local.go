package syncman

import (
	"os"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// LocalStore is an object for storing localstore information
type LocalStore struct {
	configPath   string
	globalConfig *config.Config
	services     model.ScServices

	// Callbacks
	watchAdminCB func(clusters []*config.Admin)
}

// NewLocalStore creates a new local store
func NewLocalStore(nodeID string, ssl *config.SSL) (*LocalStore, error) {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}
	// Load the configFile from path if provided
	conf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		conf = config.GenerateEmptyConfig()
	}

	// For compatibility with v18
	if conf.ClusterConfig == nil {
		conf.ClusterConfig = &config.ClusterConfig{EnableTelemetry: true}
	}

	if ssl.Enabled {
		conf.SSL = ssl
	}
	services := model.ScServices{}
	return &LocalStore{configPath: configPath, globalConfig: conf, services: append(services, &model.Service{ID: "single-node-cluster"})}, nil
}

// Register registers space cloud to the local store
func (s *LocalStore) Register() {}

// WatchResources maintains consistency over all projects
func (s *LocalStore) WatchResources(cb func(eventType, resourceId string, resourceType config.Resource, resource interface{})) error {
	return nil
}

// WatchServices maintains consistency over all services
func (s *LocalStore) WatchServices(cb func(string, string, model.ScServices)) error {
	cb(config.ResourceAddEvent, s.services[0].ID, s.services)
	return nil
}

// WatchLicense watches over changes in license secret
func (s *LocalStore) WatchLicense(cb func(eventType, resourceID string, resourceType config.Resource, resource *config.License)) {
	cb(config.ResourceAddEvent, config.GenerateResourceID("", "noProject", config.ResourceLicense, "license"), config.ResourceLicense, s.globalConfig.License)
}

func (s *LocalStore) SetLicense(ctx context.Context, resourceID string, resource *config.License) error {
	s.globalConfig.License = resource
	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// SetResource sets the project of the local globalConfig
func (s *LocalStore) SetResource(ctx context.Context, resourceID string, resource interface{}) error {
	if err := validateResource(ctx, config.ResourceAddEvent, s.globalConfig, resourceID, "", resource); err != nil {
		return err
	}
	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// DeleteResource deletes the project from the local gloablConfig
func (s *LocalStore) DeleteResource(ctx context.Context, resourceID string) error {
	if err := validateResource(ctx, config.ResourceDeleteEvent, s.globalConfig, resourceID, "", nil); err != nil {
		return err
	}
	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// DeleteProject deletes all the config resources which matches label projectId
func (s *LocalStore) DeleteProject(ctx context.Context, projectID string) error {
	delete(s.globalConfig.Projects, projectID)
	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// GetGlobalConfig gets config all projects
func (s *LocalStore) GetGlobalConfig() (*config.Config, error) {
	return s.globalConfig, nil
}
