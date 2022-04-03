package configman

import (
	"context"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/configman/connector"
	"github.com/spacecloud-io/space-cloud/model"
	"go.uber.org/zap"
)

var connectorPool = caddy.NewUsagePool()
var poolKey = "configman-poolkey"

// ConfigMan manages all the database modules
type ConfigMan struct {
	// The config this app needs
	StoreType string `json:"type,omitempty"`
	Path      string `json:"path,omitempty"`

	// For internal usage
	logger     *zap.Logger
	Connectors connector.ConfigManConnector
}

// CaddyModule returns the Caddy module information.
func (ConfigMan) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "configman",
		New: func() caddy.Module { return new(ConfigMan) },
	}
}

// Provision sets up the file loader module.
func (l *ConfigMan) Provision(ctx caddy.Context) error {
	l.logger = ctx.Logger(l)

	return nil
}

// Start begins the configman app operations
func (l *ConfigMan) Start() error {

	val, _, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
		return connector.New(l.logger, l.StoreType, l.Path)
	})
	if err != nil {
		l.logger.Error("Unable to open configman connector", zap.String("store-type", l.StoreType), zap.String("path", l.Path), zap.Error(err))

	}
	connector := val.(*connector.Connector)
	l.Connectors = connector.Connector

	return nil
}

// ApplyResource applies resource in the store
func (l *ConfigMan) ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error {
	return l.Connectors.ApplyResource(ctx, resourceObj)
}

// GetResource gets resource from the store
func (l *ConfigMan) GetResource(ctx context.Context, resourceMeta *model.ResourceMeta) (*model.ResourceObject, error) {
	return l.Connectors.GetResource(ctx, resourceMeta)
}

// GetResources gets resources from the store
func (l *ConfigMan) GetResources(ctx context.Context, resourceMeta *model.ResourceMeta) (*model.ListResourceObjects, error) {
	return l.Connectors.GetResources(ctx, resourceMeta)
}

// DeleteResource delete resource from the store
func (l *ConfigMan) DeleteResource(ctx context.Context, resourceMeta *model.ResourceMeta) error {
	return l.Connectors.DeleteResource(ctx, resourceMeta)
}

// DeleteResources delete resources from the store
func (l *ConfigMan) DeleteResources(ctx context.Context, resourceMeta *model.ResourceMeta) error {
	return l.Connectors.DeleteResource(ctx, resourceMeta)
}

// Stop ends the app operations
func (l *ConfigMan) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*ConfigMan)(nil)
	_ caddy.App         = (*ConfigMan)(nil)
)
