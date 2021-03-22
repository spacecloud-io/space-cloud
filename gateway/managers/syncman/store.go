package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// Store abstracts the implementation of letsencrypt storage operations
type Store interface {
	WatchServices(cb func(projects scServices)) error
	WatchResources(cb func(eventType, resourceId string, resourceType config.Resource, resource interface{})) error

	Register()

	SetResource(ctx context.Context, resourceID string, resource interface{}) error
	DeleteResource(ctx context.Context, resourceID string) error

	GetGlobalConfig() (*config.Config, error)
}
