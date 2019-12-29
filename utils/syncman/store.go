package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/config"
)

type Store interface {
	WatchProjects(cb func(projects []*config.Project)) error
	WatchServices(cb func(projects scServices)) error

	Register()

	SetProject(ctx context.Context, project *config.Project) error
	DeleteProject(ctx context.Context, projectID string) error
}
