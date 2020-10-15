package crud

import (
	"context"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, crud config.DatabaseConfigs) error {
	m.Lock()
	defer m.Unlock()

	if len(crud) > 1 {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Crud module cannot have more than 1 database", nil, map[string]interface{}{"project": project})
	}

	m.project = project

	// clear previous data loader1
	m.dataLoader = loader{loaderMap: map[string]*dataloader.Loader{}}

	// Create a new crud blocks
	for _, v := range crud {
		if v.Type == "" {
			v.Type = v.DbAlias
		}

		// set default database name to project id
		if v.DBName == "" {
			v.DBName = project
		}

		if v.Limit == 0 {
			v.Limit = model.DefaultFetchLimit
		}

		// check if connection string starts with secrets
		secretName, isSecretExists := splitConnectionString(v.Conn)
		connectionString := v.Conn
		if isSecretExists {
			var err error
			connectionString, err = m.getSecrets(project, secretName, "CONN")
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to fetch secret from runner", err, map[string]interface{}{"project": project})
			}
		}

		if m.block != nil {
			m.block.SetQueryFetchLimit(v.Limit)
			// Skip if the connection string is the same
			if m.block.IsSame(connectionString, v.DBName) {
				continue
			}
			// Close the previous database connection
			if err := m.block.Close(); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to close database connections", err, map[string]interface{}{"project": project})
			}
		}

		var c Crud
		var err error

		v.Type = strings.TrimPrefix(v.Type, "sql-")
		c, err = m.initBlock(model.DBType(v.Type), v.Enabled, connectionString, v.DBName)

		if v.Enabled {
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Cannot connect to database", err, map[string]interface{}{"project": project, "dbAlias": v.DbAlias, "dbType": v.Type, "conn": v.Conn, "logicalDbName": v.DBName})
			}
			helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), "Successfully connected to database", map[string]interface{}{"project": project, "dbAlias": v.DbAlias, "dbType": v.Type})
		}

		m.dbType = v.Type
		m.config = v
		m.block = c
		m.alias = strings.TrimPrefix(v.DbAlias, "sql-")
	}

	return nil
}

// SetPreparedQueryConfig set prepared query config of crud module
func (m *Module) SetPreparedQueryConfig(ctx context.Context, prepQueries config.DatabasePreparedQueries) error {
	m.Lock()
	defer m.Unlock()

	if m.block == nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Unable to get database connection, crud module not initialized", nil)
		return nil
	}

	temp := make(config.DatabasePreparedQueries)
	for _, preparedQuery := range prepQueries {
		if preparedQuery.DbAlias != m.alias {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unknown dbAlias (%s) provided in prepared query", preparedQuery.DbAlias), nil, map[string]interface{}{"queryId": preparedQuery.ID})
		}
		temp[getPreparedQueryKey(strings.TrimPrefix(preparedQuery.DbAlias, "sql-"), preparedQuery.ID)] = preparedQuery
	}
	m.queries = temp
	return nil
}

// SetSchemaConfig set schema config of crud module
func (m *Module) SetSchemaConfig(ctx context.Context, schemas config.DatabaseSchemas) error {
	m.Lock()
	defer m.Unlock()

	if m.block == nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Unable to get database connection, crud module not initialized", nil)
		return nil
	}

	m.closeBatchOperation()
	if err := m.initBatchOperation(m.project, schemas); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to initialized database batch operation", err, nil)
	}
	return nil
}

// SetGetSecrets sets the GetSecrets function
func (m *Module) SetGetSecrets(function utils.GetSecrets) {
	m.Lock()
	defer m.Unlock()

	m.getSecrets = function
}

// SetSchema sets the schema module
func (m *Module) SetSchema(s model.SchemaCrudInterface) {
	m.schema = s
}

// SetHooks sets the internal hooks
func (m *Module) SetHooks(hooks *model.CrudHooks, metricHook model.MetricCrudHook) {
	m.hooks = hooks
	m.metricHook = metricHook
}
