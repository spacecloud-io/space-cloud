package database

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/modules/database/connectors"
)

func init() {
	caddy.RegisterModule(App{})
	configman.RegisterConfigController("database")
}

var connectorPool = caddy.NewUsagePool()

// App manages all the database modules
type App struct {
	// The config this app needs
	DBConfigs map[string]*Config `json:"dbConfigs,omitempty"`

	// For internal usage
	logger     *zap.Logger
	connectors map[string]*connectors.Module
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "database",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the file loader module.
func (l *App) Provision(ctx caddy.Context) error {
	l.logger = ctx.Logger(l)

	// Create an empty connectors map
	if l.connectors == nil {
		l.connectors = make(map[string]*connectors.Module)
	}

	// Iterate over all database configs
	for k, dbConfig := range l.DBConfigs {
		if !dbConfig.Connector.Enabled {
			continue
		}

		// Split the config key
		projectID, dbAlias := SplitDBConfigKey(k)

		// TODO: Update the connection string with the secret value if provided

		// Load the module. We will use a special poolkey which uniquely identifies a database
		// connection. If any paramter of this connection changes we want to create a new
		// connector instance and destroy the old one.
		poolKey := generateUniqueDBKey(projectID, dbConfig.Connector)
		val, loaded, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
			return connectors.New(l.logger, projectID, dbConfig.Connector, dbConfig.Schemas, dbConfig.PreparedQueries)
		})
		if err != nil {
			l.logger.Error("Unable to open database connector",
				zap.String("project", projectID), zap.String("dbAlias", dbAlias), zap.Error(err))
			continue
		}
		module := val.(*connectors.Module)

		// Update config if it was already present
		if loaded {
			module.UpdateConfig(l.logger, dbConfig.Connector, dbConfig.Schemas, dbConfig.PreparedQueries)
		}

		// Store the connector in the map for future reference
		l.connectors[k] = module
	}

	return nil
}

// Start begins the database app operations
func (l *App) Start() error {
	return nil
}

// Stop ends the database app operations
func (l *App) Stop() error {
	return nil
}

// Cleanup clean up the app
func (l *App) Cleanup() error {
	// Iterate over all database configs
	for k, dbConfig := range l.DBConfigs {
		// Split the config key
		projectID, dbAlias := SplitDBConfigKey(k)

		// Delete the connector from pool. Note, pool will actually delete the config only if
		// it isn't referenced in the newer config provided
		poolKey := generateUniqueDBKey(projectID, dbConfig.Connector)
		_, err := connectorPool.Delete(poolKey)
		if err != nil {
			l.logger.Error("Unable to gracefully close database connector",
				zap.String("project", projectID), zap.String("dbAlias", dbAlias), zap.Error(err))
		}
	}

	return nil
}

// Interface guards
var (
	_ caddy.Provisioner  = (*App)(nil)
	_ caddy.CleanerUpper = (*App)(nil)
	_ caddy.App          = (*App)(nil)
	_ model.ConfigCtrl   = (*App)(nil)
)
