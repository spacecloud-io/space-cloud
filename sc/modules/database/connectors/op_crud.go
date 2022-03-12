package connectors

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/modules/database/connectors/schema"
	"github.com/spacecloud-io/space-cloud/utils"
)

// Create inserts a document (or multiple when op is "all") into the database based on dbType
func (m *Module) Create(ctx context.Context, col string, req *model.CreateRequest, params model.RequestParams) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Validate the documents to be inserted
	if err := schema.ValidateCreateOperation(ctx, m.dbConfig.DbAlias, m.dbConfig.Type, col, m.schemaDoc, req); err != nil {
		return err
	}

	// TODO: fix integration hooks logic
	// params.Payload = req
	// hookResponse := m.integrationMan.InvokeHook(ctx, params)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return err
	// 	}

	// 	// Gracefully return
	// 	return nil
	// }

	// Check if the client is safe
	if err := m.connector.IsClientSafe(ctx); err != nil {
		return err
	}

	// var n int64
	var err error
	if req.IsBatch {
		// add the request for batch operation
		_, err = m.createBatch(ctx, col, req.Document)
	} else {
		// Perform the create operation
		_, err = m.connector.Create(ctx, col, req)
	}

	// TODO: Fix the metric hook logic
	// // Invoke the metric hook if the operation was successful
	// if err == nil {
	// 	m.metricHook(m.project, m.dbConfig.DbAlias, col, n, model.Create)
	// }

	return err
}

// Read returns the documents(s) which match a query from the database based on dbType
func (m *Module) Read(ctx context.Context, col string, req *model.ReadRequest, params model.RequestParams) (interface{}, *model.SQLMetaData, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Adjust where clause
	if err := schema.AdjustWhereClause(ctx, m.dbConfig.DbAlias, model.DBType(m.dbConfig.Type), col, m.schemaDoc, req.Find); err != nil {
		return nil, nil, err
	}

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return nil, nil, err
	}

	// TODO: Fix integration hook logic
	// params.Payload = req
	// hookResponse := m.integrationMan.InvokeHook(ctx, params)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return nil, nil, err
	// 	}

	// 	// Gracefully return
	// 	return hookResponse.Result(), nil, nil
	// }

	// Check if we want to use the dataloader
	if req.IsBatch {
		key := model.ReadRequestKey{DBType: m.dbConfig.Type, DBAlias: m.dbConfig.DbAlias, Col: col, HasOptions: req.Options.HasOptions, Req: *req, ReqParams: params}
		dataLoader, ok := m.getLoader(col)
		if !ok {
			dataLoader = m.createLoader(col)
		}
		data, err := dataLoader.Load(ctx, key)()
		if err != nil {
			return nil, nil, err
		}
		res := data.(queryResult)
		if res.metaData != nil {
			res.metaData.DbAlias = m.dbConfig.DbAlias
			res.metaData.Col = col
		}
		return res.doc, res.metaData, err
	}

	// TODO: Fix caching logic
	// dbCacheOptions, err := m.caching.GetDatabaseKey(ctx, m.project, m.dbConfig.DbAlias, col, req)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// TODO: Add metric hook for cache

	// See if result is present in cache
	var metaData *model.SQLMetaData
	var result interface{}
	var err error
	// if !dbCacheOptions.IsCacheHit() {
	// Perform the read operation
	// var n int64
	// var cacheJoinInfo map[string]map[string]string
	_, result, _, metaData, err = m.connector.Read(ctx, col, req)

	// // Set result in cache & invoke the metric hook if the operation was successful
	// if err == nil {
	// 	if err := m.caching.SetDatabaseKey(ctx, m.project, m.dbConfig.DbAlias, col, &model.CacheDatabaseResult{MetricCount: n, Result: result}, dbCacheOptions, req.Cache, cacheJoinInfo); err != nil {
	// 		return nil, nil, err
	// 	}

	// 	m.metricHook(m.project, m.dbConfig.DbAlias, col, n, model.Read)
	// }
	// } else {
	// 	// Make a metadata object for cached results
	// 	metaData = &model.SQLMetaData{QueryTime: "0s", SQL: "fetched from cache"}

	// 	cacheResult := dbCacheOptions.GetDatabaseResult()
	// 	result = cacheResult.Result
	// 	m.metricHook(m.project, m.dbConfig.DbAlias, col, cacheResult.MetricCount, model.Read)
	// }

	// Process the response
	if err := schema.CrudPostProcess(ctx, m.dbConfig.DbAlias, m.dbConfig.Type, col, m.schemaDoc, result); err != nil {
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error executing read request in crud module unable to perform schema post process for un marshalling json for project (%s) col (%s)", m.project, col), err, nil)
	}

	if metaData != nil {
		metaData.DbAlias = m.dbConfig.DbAlias
		metaData.Col = col
	}

	return result, metaData, err
}

// Update updates the documents(s) which match a query from the database based on dbType
func (m *Module) Update(ctx context.Context, col string, req *model.UpdateRequest, params model.RequestParams) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Validate the update operation that needs to be performed
	if err := schema.ValidateUpdateOperation(ctx, m.dbConfig.DbAlias, m.dbConfig.Type, col, req.Operation, req.Update, req.Find, m.schemaDoc); err != nil {
		return err
	}

	// TODO: Fix integration hooks logic
	// params.Payload = req
	// hookResponse := m.integrationMan.InvokeHook(ctx, params)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return err
	// 	}

	// 	// Gracefully return
	// 	return nil
	// }

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return err
	}

	// Adjust where clause
	if err := schema.AdjustWhereClause(ctx, m.dbConfig.DbAlias, model.DBType(m.dbConfig.Type), col, m.schemaDoc, req.Find); err != nil {
		return err
	}

	// Perform the update operation
	_, err := m.connector.Update(ctx, col, req)

	// TODO: Fix metric hook logic
	// // Invoke the metric hook if the operation was successful
	// if err == nil {
	// 	m.metricHook(m.project, m.dbConfig.DbAlias, col, n, model.Update)
	// }

	return err
}

// Delete removes the documents(s) which match a query from the database based on dbType
func (m *Module) Delete(ctx context.Context, col string, req *model.DeleteRequest, params model.RequestParams) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return err
	}

	// Adjust where clause
	if err := schema.AdjustWhereClause(ctx, m.dbConfig.DbAlias, model.DBType(m.dbConfig.Type), col, m.schemaDoc, req.Find); err != nil {
		return err
	}

	// TODO: Fix integration hook logic
	// params.Payload = req
	// hookResponse := m.integrationMan.InvokeHook(ctx, params)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return err
	// 	}

	// 	// Gracefully return
	// 	return nil
	// }

	// Perform the delete operation
	_, err := m.connector.Delete(ctx, col, req)

	// TODO: Fix metric hook logic
	// // Invoke the metric hook if the operation was successful
	// if err == nil {
	// 	m.metricHook(m.project, m.dbConfig.DbAlias, col, n, model.Delete)
	// }

	return err
}

// ExecPreparedQuery executes PreparedQueries request
func (m *Module) ExecPreparedQuery(ctx context.Context, id string, req *model.PreparedQueryRequest, params model.RequestParams) (interface{}, *model.SQLMetaData, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// TODO: Fix integration hooks logic
	// params.Payload = req
	// hookResponse := m.integrationMan.InvokeHook(ctx, params)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return nil, nil, err
	// 	}

	// 	// Gracefully return
	// 	return hookResponse.Result(), nil, nil
	// }

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return nil, nil, err
	}

	// Check if prepared query exists
	preparedQuery, p := m.dbPreparedQueries[id]
	if !p {
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Prepared Query for given id (%s) does not exist", id), nil, nil)
	}

	// Load the arguments
	var args []interface{}
	for i := 0; i < len(preparedQuery.Arguments); i++ {
		arg, err := utils.LoadValue(preparedQuery.Arguments[i], map[string]interface{}{"args": req.Params, "auth": params})
		if err != nil {
			return nil, nil, err
		}
		args = append(args, arg)
	}

	// Fire the query and return the result
	_, b, metaData, err := m.connector.RawQuery(ctx, preparedQuery.SQL, req.Debug, args)
	if metaData != nil {
		metaData.DbAlias = m.dbConfig.DbAlias
		metaData.Col = id
	}
	return b, metaData, err
}

// Aggregate performs an aggregation defined via the pipeline
func (m *Module) Aggregate(ctx context.Context, col string, req *model.AggregateRequest, params model.RequestParams) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// TODO: Fix prepared query logic
	// params.Payload = req
	// hookResponse := m.integrationMan.InvokeHook(ctx, params)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return nil, err
	// 	}

	// 	// Gracefully return
	// 	return hookResponse.Result(), nil
	// }

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return nil, err
	}

	return m.connector.Aggregate(ctx, col, req)
}

// Batch performs a batch operation on the database
func (m *Module) Batch(ctx context.Context, req *model.BatchRequest, params model.RequestParams) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, r := range req.Requests {
		switch r.Type {
		case string(model.Create):
			v := &model.CreateRequest{Document: r.Document, Operation: r.Operation}
			if err := schema.ValidateCreateOperation(ctx, m.dbConfig.DbAlias, m.dbConfig.Type, r.Col, m.schemaDoc, v); err != nil {
				return err
			}
			r.Document = v.Document
			r.Operation = v.Operation
		case string(model.Update):
			if err := schema.ValidateUpdateOperation(ctx, m.dbConfig.DbAlias, m.dbConfig.Type, r.Col, r.Operation, r.Update, r.Find, m.schemaDoc); err != nil {
				return err
			}
		}
	}

	// TODO: Fix integration hooks logic
	// params.Payload = req
	// hookResponse := m.integrationMan.InvokeHook(ctx, params)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return err
	// 	}

	// 	// Gracefully return
	// 	return nil
	// }

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return err
	}

	// Perform the batch operation
	counts, err := m.connector.Batch(ctx, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		for i, r := range req.Requests {
			m.metricHook(m.project, m.dbConfig.DbAlias, r.Col, counts[i], model.OperationType(r.Type))
		}
	}

	return err
}

// RawBatch performs a db operation for schema creation
func (m *Module) RawBatch(ctx context.Context, batchedQueries []string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return err
	}

	return m.connector.RawBatch(ctx, batchedQueries)
}

// GetCollections returns collection / tables name of specified database
func (m *Module) GetCollections(ctx context.Context) ([]model.DatabaseCollections, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return nil, err
	}

	return m.connector.GetCollections(ctx)
}

// GetConnectionState gets the current state of client
func (m *Module) GetConnectionState(ctx context.Context) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return false
	}

	return m.connector.GetConnectionState(ctx)
}

// DeleteTable drop specified table from database
func (m *Module) DeleteTable(ctx context.Context, col string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.connector.IsClientSafe(ctx); err != nil {
		return err
	}

	return m.connector.DeleteCollection(ctx, col)
}
