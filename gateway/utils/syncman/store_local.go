package syncman

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"golang.org/x/net/context"
)

// LocalStore is an object for storing localstore information
type LocalStore struct {
}

// NewLocalStore creates a new local store
func NewLocalStore() (*LocalStore, error) {
	return &LocalStore{}, nil
}

// Register registers space cloud to the local store
func (s *LocalStore) Register() {}

// WatchProjects maintains consistency over all projects
func (s *LocalStore) WatchProjects(cb func(projects []*config.Project)) error {
	return nil
}

// WatchServices maintains consistency over all services
func (s *LocalStore) WatchServices(cb func(scServices)) error {
	return nil
}

//WatchGlobalConfig maintains consistency between all instances of sc
func (s *LocalStore) WatchGlobalConfig(cb func(projects *config.GlobalConfig)) error {
	return nil
}

// SetProject sets the project of the local store
func (s *LocalStore) SetProject(ctx context.Context, project *config.Project) error {
	return nil
}

// DeleteProject deletes the project from the local store
func (s *LocalStore) DeleteProject(ctx context.Context, projectID string) error {
	return nil
}
