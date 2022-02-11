package database

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/modules/database/connectors"
)

func init() {
	caddy.RegisterModule(DatabaseApp{})
}

var connectorPool = caddy.NewUsagePool()

// DatabaseApp manages all the database modules
type DatabaseApp struct {
	// The config this app needs
	DBConfigs map[string]*Config `json:"dbConfigs,omitempty"`

	// For internal usage
	logger     *zap.Logger
	connectors map[string]*connectors.Module
}

// CaddyModule returns the Caddy module information.
func (DatabaseApp) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "database",
		New: func() caddy.Module { return new(DatabaseApp) },
	}
}

// Provision sets up the file loader module.
func (l *DatabaseApp) Provision(ctx caddy.Context) error {
	l.logger = ctx.Logger(l)

	return nil
}

// Start begins the database app operations
func (l *DatabaseApp) Start() error {
	// Create an empty connectors map
	if l.connectors == nil {
		l.connectors = make(map[string]*connectors.Module, 0)
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
			return connectors.New(projectID, dbConfig.Connector, dbConfig.Schemas, dbConfig.PreparedQueries)
		})
		if err != nil {
			l.logger.Error("Unable to open database connector",
				zap.String("project", projectID), zap.String("dbAlias", dbAlias), zap.Error(err))
			continue
		}
		module := val.(*connectors.Module)

		// Update config if it was already present
		if loaded {
			module.UpdateConfig(dbConfig.Connector, dbConfig.Schemas, dbConfig.PreparedQueries)
		}

		// Store the connector in the map for future reference
		l.connectors[k] = module
	}
	return nil
}

// Stop ends the database app operations
func (l *DatabaseApp) Stop() error {
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
	_ caddy.Provisioner = (*DatabaseApp)(nil)
	_ caddy.App         = (*DatabaseApp)(nil)
)
