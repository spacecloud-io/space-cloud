package configman

import (
	"context"

	"github.com/caddyserver/caddy/v2"
	connector "github.com/spacecloud-io/space-cloud/managers/configman/store"
	"github.com/spacecloud-io/space-cloud/model"
	"go.uber.org/zap"
)

var connectorPool = caddy.NewUsagePool()
var poolKey = "configman-poolkey"

// Store manages all the store modules
type Store struct {
	// The config this app needs
	StoreType string `json:"type,omitempty"`
	Path      string `json:"path,omitempty"`

	// For internal usage
	logger    *zap.Logger
	connector connector.Store
}

// CaddyModule returns the Caddy module information.
func (Store) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "config_store",
		New: func() caddy.Module { return new(Store) },
	}
}

// Provision sets up the store module.
func (s *Store) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger(s)

	val, _, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
		return connector.New(s.logger, s.StoreType, s.Path)
	})
	if err != nil {
		s.logger.Error("Unable to open configman connector", zap.String("store-type", s.StoreType), zap.String("path", s.Path), zap.Error(err))
		return nil
	}
	s.connector = val.(connector.Store)
	s.connector.SetLogger(s.logger)
	return nil
}

// Start begins the configman app operations
func (s *Store) Start() error {
	return nil
}

// ApplyResource applies resource in the store
func (s *Store) ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error {
	return s.connector.ApplyResource(ctx, resourceObj)
}

// GetResource gets resource from the store
func (s *Store) GetResource(ctx context.Context, resourceMeta *model.ResourceMeta) (*model.ResourceObject, error) {
	return s.connector.GetResource(ctx, resourceMeta)
}

// GetResources gets resources from the store
func (s *Store) GetResources(ctx context.Context, resourceMeta *model.ResourceMeta) (*model.ListResourceObjects, error) {
	return s.connector.GetResources(ctx, resourceMeta)
}

// DeleteResource delete resource from the store
func (s *Store) DeleteResource(ctx context.Context, resourceMeta *model.ResourceMeta) error {
	return s.connector.DeleteResource(ctx, resourceMeta)
}

// DeleteResources delete resources from the store
func (s *Store) DeleteResources(ctx context.Context, resourceMeta *model.ResourceMeta) error {
	return s.connector.DeleteResources(ctx, resourceMeta)
}

// Stop ends the app operations
func (s *Store) Stop() error {
	return nil
}

// Cleanup destroys the configman app
func (s *Store) Cleanup() error {
	_, err := connectorPool.Delete(poolKey)
	if err != nil {
		s.logger.Error("Unable to gracefully close store connector", zap.Error(err))
	}
	return err
}

// Interface guards
var (
	_ caddy.Provisioner  = (*Store)(nil)
	_ caddy.CleanerUpper = (*Store)(nil)
	_ caddy.App          = (*Store)(nil)
)
