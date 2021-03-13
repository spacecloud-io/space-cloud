package syncman

import (
	"os"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// LocalStore is an object for storing localstore information
type LocalStore struct {
	configPath   string
	globalConfig *config.Config
}

// NewLocalStore creates a new local store
func NewLocalStore(ssl *config.SSL) (*LocalStore, error) {
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
	return &LocalStore{configPath: configPath, globalConfig: conf}, nil
}

// Register registers space cloud to the local store
func (s *LocalStore) Register() {}

// WatchResources maintains consistency over all projects
func (s *LocalStore) WatchResources(cb func(eventType, resourceId string, resourceType config.Resource, resource interface{})) error {
	return nil
}

// SetResource sets the project of the local globalConfig
func (s *LocalStore) SetResource(ctx context.Context, resourceID string, resource interface{}) error {
	if err := updateResource(ctx, config.ResourceAddEvent, s.globalConfig, resourceID, "", resource); err != nil {
		return err
	}
	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// DeleteResource deletes the project from the local gloablConfig
func (s *LocalStore) DeleteResource(ctx context.Context, resourceID string) error {
	if err := updateResource(ctx, config.ResourceDeleteEvent, s.globalConfig, resourceID, "", nil); err != nil {
		return err
	}
	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// GetGlobalConfig gets config all projects
func (s *LocalStore) GetGlobalConfig() (*config.Config, error) {
	return s.globalConfig, nil
}
