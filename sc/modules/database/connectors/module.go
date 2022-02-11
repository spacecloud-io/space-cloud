package connectors

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/spaceuptech/helpers"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/modules/database/connectors/sql"
)

// Module is the root block providing convenient wrappers
type Module struct {
	lock sync.RWMutex

	// batch operation
	batchMapTableToChan batchMap // every table gets mapped to group of channels

	// Config objects
	project           string
	dbConfig          *config.DatabaseConfig
	dbSchemas         config.DatabaseSchemas
	dbPreparedQueries config.DatabasePreparedQueries

	dataLoader loader
	// Variables to store the hooks
	metricHook model.MetricCrudHook

	connector Connector
	// function to get secrets from runner
	// TODO: Fix secrets
	//getSecrets utils.GetSecrets

	// Schema module
	// schemaDoc model.Type
}

// New create a new instance of the Module object
func New(projectID string, dbConfig *config.DatabaseConfig, dbSchemas config.DatabaseSchemas, dbPreparedQueries config.DatabasePreparedQueries) (*Module, error) {
	m := &Module{
		batchMapTableToChan: make(batchMap),
		project:             projectID,
		dbConfig:            dbConfig,
		dbSchemas:           dbSchemas,
		dbPreparedQueries:   sanitizePrepareQueries(dbPreparedQueries),
		dataLoader:          loader{loaderMap: map[string]*dataloader.Loader{}},
	}

	// Set the database type to alias if it isn't provided
	if dbConfig.Type == "" {
		dbConfig.Type = dbConfig.DbAlias
	}
	dbConfig.Type = strings.TrimPrefix(dbConfig.Type, "sql-")

	// Set default database name to project id
	if dbConfig.DBName == "" {
		dbConfig.DBName = projectID
	}

	// Set the limit if not provided by end user
	if dbConfig.Limit == 0 {
		dbConfig.Limit = model.DefaultFetchLimit
	}

	// TODO: Load the connection string from a secret if required
	connectionString := dbConfig.Conn

	// Create a new connector object
	c, err := m.initConnector(model.DBType(dbConfig.Type), connectionString, dbConfig.DBName, dbConfig.DriverConf)
	if err != nil {
		return nil, err
	}
	m.connector = c

	// Update changable config
	m.UpdateConfig(dbConfig, dbSchemas, dbPreparedQueries)

	return m, nil
}

// UpdateConfig updates a connectors config
func (m *Module) UpdateConfig(dbConfig *config.DatabaseConfig, dbSchemas config.DatabaseSchemas, dbPreparedQueries config.DatabasePreparedQueries) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Set the limit if not provided by end user
	if dbConfig.Limit == 0 {
		dbConfig.Limit = model.DefaultFetchLimit
	}

	// Update the connectors config
	m.connector.SetQueryFetchLimit(dbConfig.Limit)
	m.dbConfig.BatchTime = dbConfig.BatchTime
	m.dbConfig.BatchRecords = dbConfig.BatchRecords

	// Restart the batching operations if the configured tables have changed
	if !areSchemasSimilar(m.dbSchemas, dbSchemas) {
		m.closeBatchOperation()
		m.initBatchOperation()

		// Clear previous data loaders as well
		m.dataLoader = loader{loaderMap: map[string]*dataloader.Loader{}}
	}

	// Update all config objects
	m.dbConfig = dbConfig
	m.dbSchemas = dbSchemas
	m.dbPreparedQueries = dbPreparedQueries
}

func (m *Module) initConnector(dbType model.DBType, connection, dbName string, driverConf config.DriverConfig) (Connector, error) {
	switch dbType {
	// TODO: Add support for the remaining connectors soon
	// case model.Mongo:
	// 	return mgo.Init(enabled, connection, dbName, driverConf)
	// case model.EmbeddedDB:
	// 	return bolt.Init(enabled, connection, dbName)
	case model.MySQL, model.Postgres, model.SQLServer:
		// Attempt to initialse the sql connector
		c, err := sql.Init(dbType, connection, dbName, driverConf)
		if err != nil {
			return nil, err
		}

		// Create a database for the user is it doesn't already exists
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := c.CreateDatabaseIfNotExist(ctx, dbName); err != nil {
			_ = c.Close()
			return nil, err
		}

		// For mysql database, create a new dbname specific connection string
		if dbType == model.MySQL {
			_ = c.Close()
			return sql.Init(dbType, fmt.Sprintf("%s%s", connection, dbName), dbName, driverConf)
		}
		return c, err
	default:
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unsupported database (%s) provided", dbType), nil, map[string]interface{}{})
	}
}

// Destruct destroys the database module
func (m *Module) Destruct() error {
	// Acquire a lock
	m.lock.Lock()
	defer m.lock.Unlock()

	m.dbConfig = nil
	m.dbSchemas = nil
	m.dbPreparedQueries = nil

	for k := range m.dataLoader.loaderMap {
		delete(m.dataLoader.loaderMap, k)
	}

	// Close the batching goroutine
	m.closeBatchOperation()

	// Close the connector
	return m.connector.Close()
}
