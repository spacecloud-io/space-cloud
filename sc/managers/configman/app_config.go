package configman

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/model"
	"go.uber.org/zap"
)

// ConfigMan manages all the operation and config type definitions
type ConfigMan struct {
	// For internal usage
	logger                         *zap.Logger
	configControllerDefinitions    map[string]model.ConfigTypes    // Key = moduleName
	operationControllerDefinitions map[string]model.OperationTypes // Key = moduleName
}

// CaddyModule returns the Caddy module information.
func (ConfigMan) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "configman",
		New: func() caddy.Module { return new(ConfigMan) },
	}
}

// Provision sets up the store module.
func (c *ConfigMan) Provision(ctx caddy.Context) error {
	c.logger = ctx.Logger(c)

	// Acquire a read lock
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	// Load all the operation type definitions
	c.operationControllerDefinitions = make(map[string]model.OperationTypes, len(registeredOperationControllers))
	for _, module := range registeredOperationControllers {
		app, err := ctx.App(module)
		if err != nil {
			c.logger.Error("Unable to load operation controller application", zap.Error(err))
			return err
		}

		c.operationControllerDefinitions[module] = app.(model.OperationCtrl).GetOperationTypes()
	}

	// Load all the operation type definitions
	c.configControllerDefinitions = make(map[string]model.ConfigTypes, len(registeredConfigControllers))
	for _, module := range registeredConfigControllers {
		app, err := ctx.App(module)
		if err != nil {
			c.logger.Error("Unable to load operation controller application", zap.Error(err))
			return err
		}

		c.configControllerDefinitions[module] = app.(model.ConfigCtrl).GetConfigTypes()
	}

	return nil
}

// Start begins the app's operations
func (c *ConfigMan) Start() error {
	return nil
}

// Stop ends the app's operations
func (c *ConfigMan) Stop() error {
	return nil
}

// GetOperationTypes returns the operation type definitions accross all modules
func (c *ConfigMan) GetOperationTypes() map[string]model.OperationTypes {
	return c.operationControllerDefinitions
}

// GetConfigTypes returns the config type definitions accross all modules
func (c *ConfigMan) GetConfigTypes() map[string]model.ConfigTypes {
	return c.configControllerDefinitions
}

// Interface guards
var (
	_ caddy.Provisioner = (*ConfigMan)(nil)
	_ caddy.App         = (*ConfigMan)(nil)
)
