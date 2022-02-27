package routing

import (
	"context"
	"sync"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Routing manages the routing functionality of space cloud
type Routing struct {
	lock sync.RWMutex

	routes       config.Routes
	globalConfig *config.GlobalRoutesConfig
	caching      cachingInterface
	goTemplates  map[string]*template.Template
}

// New creates a new instance of the routing module
func New() *Routing {
	return &Routing{routes: make(config.Routes, 0), goTemplates: map[string]*template.Template{}, globalConfig: new(config.GlobalRoutesConfig)}
}

// SetCachingModule sets caching module
func (r *Routing) SetCachingModule(c cachingInterface) {
	r.caching = c
}

type cachingInterface interface {
	SetIngressRouteKey(ctx context.Context, redisKey string, cache *config.ReadCacheOptions, result *model.CacheIngressRoute) error
	GetIngressRoute(ctx context.Context, routeID string, cacheOptions []interface{}) (string, bool, *model.CacheIngressRoute, error)
}
