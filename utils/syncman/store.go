package syncman

import "github.com/spaceuptech/space-cloud/config"

type Store interface {
	WatchProjects(cb func(projects []*config.Project)) error
	WatchServices(cb func(projects scServices)) error

	Register()

	SetProject(project *config.Project) error
	DeleteProject(projectID string) error
}
