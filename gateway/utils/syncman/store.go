package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// Store abstracts the implementation of letsencrypt storage operations
type Store interface {
	WatchProjects(cb func(projects []*config.Project)) error
	WatchServices(cb func(projects scServices)) error
	WatchAdminConfig(cb func(clusters []*config.Admin)) error

	Register()

	SetAdminConfig(ctx context.Context, adminConfig *config.Admin) error
	SetProject(ctx context.Context, project *config.Project) error
	DeleteProject(ctx context.Context, projectID string) error
}
