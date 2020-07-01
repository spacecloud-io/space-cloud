package syncman

import (
	"errors"
	"os"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// LocalStore is an object for storing localstore information
type LocalStore struct {
	configPath   string
	globalConfig *config.Config
	services     scServices
}

// NewLocalStore creates a new local store
func NewLocalStore(nodeID, advertiseAddr string, ssl *config.SSL) (Store, error) {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}
	// Load the configFile from path if provided
	conf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		conf = config.GenerateEmptyConfig()
	}
	if ssl.Enabled {
		conf.SSL = ssl
	}
	services := scServices{}
	return &LocalStore{configPath: configPath, globalConfig: conf, services: append(services, &service{id: nodeID, addr: advertiseAddr})}, nil
}

// Register registers space cloud to the local store
func (s *LocalStore) Register() {}

// WatchProjects maintains consistency over all projects
func (s *LocalStore) WatchProjects(cb func(projects []*config.Project)) error {
	cb(s.globalConfig.Projects)
	return nil
}

// WatchServices maintains consistency over all services
func (s *LocalStore) WatchServices(cb func(scServices)) error {
	cb(s.services)
	return nil
}

// WatchAdminConfig sets the admin config when the gateways is started
func (s *LocalStore) WatchAdminConfig(cb func(clusters []*config.Admin)) error {
	cb([]*config.Admin{s.globalConfig.Admin})
	return nil
}

// SetAdminConfig maintains consistency between all instances of sc
func (s *LocalStore) SetAdminConfig(ctx context.Context, adminConfig *config.Admin) error {
	s.globalConfig.Admin = adminConfig
	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// GetAdminConfig returns the admin config present in the store
func (s *LocalStore) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	return nil, errors.New("not implemented for local store")
}

// SetProject sets the project of the local globalConfig
func (s *LocalStore) SetProject(ctx context.Context, project *config.Project) error {
	doesExist := false
	for i, v := range s.globalConfig.Projects {
		if v.ID == project.ID {
			doesExist = true
			s.globalConfig.Projects[i] = project
		}
	}
	if !doesExist {
		s.globalConfig.Projects = append(s.globalConfig.Projects, project)
	}

	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}

// DeleteProject deletes the project from the local gloablConfig
func (s *LocalStore) DeleteProject(ctx context.Context, projectID string) error {
	for index, project := range s.globalConfig.Projects {
		if project.ID == projectID {
			s.globalConfig.Projects = append(s.globalConfig.Projects[:index], s.globalConfig.Projects[index+1:]...)
			break
		}
	}

	return config.StoreConfigToFile(s.globalConfig, s.configPath)
}
