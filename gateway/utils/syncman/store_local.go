package syncman

import (
	"os"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"golang.org/x/net/context"
)

// LocalStore is an object for storing localstore information
type LocalStore struct {
	configPath string
	store      *config.Config
	service    *service
}

// NewLocalStore creates a new local store
func NewLocalStore(nodeID, advertiseAddr string) (*LocalStore, error) {
	configPath := os.Getenv("CONFIG")

	// Load the configFile from path if provided
	conf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		conf = config.GenerateEmptyConfig()
	}
	return &LocalStore{configPath: configPath, store: conf, service: &service{id: nodeID, addr: advertiseAddr}}, nil
}

// Register registers space cloud to the local store
func (s *LocalStore) Register() {}

// WatchProjects maintains consistency over all projects
func (s *LocalStore) WatchProjects(cb func(projects []*config.Project)) error {
	cb(s.store.Projects)
	return nil
}

// WatchServices maintains consistency over all services
func (s *LocalStore) WatchServices(cb func(scServices)) error {
	services := scServices{}
	doesExist := false
	for _, service := range services {
		if service.id == s.service.id {
			doesExist = true
			service.addr = s.service.addr
			break
		}
	}

	// add service if it doesn't exist
	if !doesExist {
		services = append(services, &service{id: s.service.id, addr: s.service.addr})
	}
	cb(services)
	return nil
}

//WatchGlobalConfig maintains consistency between all instances of sc
func (s *LocalStore) WatchGlobalConfig(cb func(projects *config.GlobalConfig)) error {
	cb(s.store.GlobalConfig)
	return nil
}

// SetProject sets the project of the local store
func (s *LocalStore) SetProject(ctx context.Context, project *config.Project) error {
	doesExist := false
	for i, v := range s.store.Projects {
		if v.ID == project.ID {
			doesExist = true
			s.store.Projects[i] = project
		}
	}
	if !doesExist {
		s.store.Projects = append(s.store.Projects, project)
	}

	config.StoreConfigToFile(s.store, s.configPath)
	return nil
}

// DeleteProject deletes the project from the local store
func (s *LocalStore) DeleteProject(ctx context.Context, projectID string) error {
	for i, v := range s.store.Projects {
		if v.ID == projectID {
			s.store.Projects[i] = s.store.Projects[len(s.store.Projects)-1]
			s.store.Projects = s.store.Projects[:len(s.store.Projects)-1]
			break
		}
	}
	config.StoreConfigToFile(s.store, s.configPath)
	return nil
}
