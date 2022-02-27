package crud

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	schemaHelpers "github.com/spaceuptech/space-cloud/gateway/modules/schema/helpers"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Create inserts a documents (or multiple when op is "all") into the database based on dbType
func (m *Module) Create(ctx context.Context, dbAlias, col string, req *model.CreateRequest, params model.RequestParams) error {
	m.RLock()
	defer m.RUnlock()

	dbType, err := m.getDBType(dbAlias)
	if err != nil {
		return err
	}
	if err := schemaHelpers.ValidateCreateOperation(ctx, dbAlias, dbType, col, m.schemaDoc, req); err != nil {
		return err
	}

	params.Payload = req
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return err
		}

		// Gracefully return
		return nil
	}

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return err
	}

	var n int64
	if req.IsBatch {
		// add the request for batch operation
		n, err = m.createBatch(ctx, m.project, dbAlias, col, req.Document)
	} else {
		// Perform the create operation
		n, err = crud.Create(ctx, col, req)
	}

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbAlias, col, n, model.Create)
	}

	return err
}

// Read returns the documents(s) which match a query from the database based on dbType
func (m *Module) Read(ctx context.Context, dbAlias, col string, req *model.ReadRequest, params model.RequestParams) (interface{}, *model.SQLMetaData, error) {
	m.RLock()
	defer m.RUnlock()

	// Adjust where clause
	dbType, err := m.getDBType(dbAlias)
	if err != nil {
		return nil, nil, err
	}
	if err := schemaHelpers.AdjustWhereClause(ctx, dbAlias, model.DBType(dbType), col, m.schemaDoc, req.Find); err != nil {
		return nil, nil, err
	}

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return nil, nil, err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return nil, nil, err
	}

	params.Payload = req
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return nil, nil, err
		}

		// Gracefully return
		return hookResponse.Result(), nil, nil
	}

	if req.IsBatch {
		dbType, err := m.getDBType(dbAlias)
		if err != nil {
			return nil, nil, err
		}
		key := model.ReadRequestKey{DBType: dbType, DBAlias: dbAlias, Col: col, HasOptions: req.Options.HasOptions, Req: *req, ReqParams: params}
		dataLoader, ok := m.getLoader(fmt.Sprintf("%s-%s-%s", m.project, dbAlias, col))
		if !ok {
			dataLoader = m.createLoader(fmt.Sprintf("%s-%s-%s", m.project, dbAlias, col))
		}
		data, err := dataLoader.Load(ctx, key)()
		if err != nil {
			return nil, nil, err
		}
		res := data.(queryResult)
		if res.metaData != nil {
			res.metaData.DbAlias = dbAlias
			res.metaData.Col = col
		}
		return res.doc, res.metaData, err
	}

	dbCacheOptions, err := m.caching.GetDatabaseKey(ctx, m.project, dbAlias, col, req)
	if err != nil {
		return nil, nil, err
	}
	// TODO: Add metric hook for cache

	// See if result is present in cache
	var metaData *model.SQLMetaData
	var result interface{}
	if !dbCacheOptions.IsCacheHit() {
		// Perform the read operation
		var n int64
		var cacheJoinInfo map[string]map[string]string
		n, result, cacheJoinInfo, metaData, err = crud.Read(ctx, col, req)

		// Set result in cache & invoke the metric hook if the operation was successful
		if err == nil {
			if err := m.caching.SetDatabaseKey(ctx, m.project, dbAlias, col, &model.CacheDatabaseResult{MetricCount: n, Result: result}, dbCacheOptions, req.Cache, cacheJoinInfo); err != nil {
				return nil, nil, err
			}

			m.metricHook(m.project, dbAlias, col, n, model.Read)
		}
	} else {
		// Make a metadata object for cached results
		metaData = &model.SQLMetaData{QueryTime: "0s", SQL: "fetched from cache"}

		cacheResult := dbCacheOptions.GetDatabaseResult()
		result = cacheResult.Result
		m.metricHook(m.project, dbAlias, col, cacheResult.MetricCount, model.Read)
	}

	// Process the response
	if err := schemaHelpers.CrudPostProcess(ctx, dbAlias, dbType, col, m.schemaDoc, result); err != nil {
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error executing read request in crud module unable to perform schema post process for un marshalling json for project (%s) col (%s)", m.project, col), err, nil)
	}

	if metaData != nil {
		metaData.DbAlias = dbAlias
		metaData.Col = col
	}

	return result, metaData, err
}

// Update updates the documents(s) which match a query from the database based on dbType
func (m *Module) Update(ctx context.Context, dbAlias, col string, req *model.UpdateRequest, params model.RequestParams) error {
	m.RLock()
	defer m.RUnlock()

	dbType, err := m.getDBType(dbAlias)
	if err != nil {
		return err
	}
	if err := schemaHelpers.ValidateUpdateOperation(ctx, dbAlias, dbType, col, req.Operation, req.Update, req.Find, m.schemaDoc); err != nil {
		return err
	}

	params.Payload = req
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return err
		}

		// Gracefully return
		return nil
	}

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return err
	}

	// Adjust where clause
	if err := schemaHelpers.AdjustWhereClause(ctx, dbAlias, model.DBType(dbType), col, m.schemaDoc, req.Find); err != nil {
		return err
	}

	// Perform the update operation
	n, err := crud.Update(ctx, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbAlias, col, n, model.Update)
	}

	return err
}

// Delete removes the documents(s) which match a query from the database based on dbType
func (m *Module) Delete(ctx context.Context, dbAlias, col string, req *model.DeleteRequest, params model.RequestParams) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return err
	}

	// Adjust where clause
	dbType, err := m.getDBType(dbAlias)
	if err != nil {
		return err
	}
	if err := schemaHelpers.AdjustWhereClause(ctx, dbAlias, model.DBType(dbType), col, m.schemaDoc, req.Find); err != nil {
		return err
	}

	params.Payload = req
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return err
		}

		// Gracefully return
		return nil
	}

	// Perform the delete operation
	n, err := crud.Delete(ctx, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbAlias, col, n, model.Delete)
	}

	return err
}

// ExecPreparedQuery executes PreparedQueries request
func (m *Module) ExecPreparedQuery(ctx context.Context, dbAlias, id string, req *model.PreparedQueryRequest, params model.RequestParams) (interface{}, *model.SQLMetaData, error) {
	m.RLock()
	defer m.RUnlock()

	params.Payload = req
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return nil, nil, err
		}

		// Gracefully return
		return hookResponse.Result(), nil, nil
	}

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return nil, nil, err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return nil, nil, err
	}

	// Check if prepared query exists
	preparedQuery, p := m.queries[getPreparedQueryKey(dbAlias, id)]
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
	_, b, metaData, err := crud.RawQuery(ctx, preparedQuery.SQL, req.Debug, args)
	if metaData != nil {
		metaData.DbAlias = dbAlias
		metaData.Col = id
	}
	return b, metaData, err
}

// Aggregate performs an aggregation defined via the pipeline
func (m *Module) Aggregate(ctx context.Context, dbAlias, col string, req *model.AggregateRequest, params model.RequestParams) (interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	params.Payload = req
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return nil, err
		}

		// Gracefully return
		return hookResponse.Result(), nil
	}

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return nil, err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return nil, err
	}

	return crud.Aggregate(ctx, col, req)
}

// Batch performs a batch operation on the database
func (m *Module) Batch(ctx context.Context, dbAlias string, req *model.BatchRequest, params model.RequestParams) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	dbType, err := m.getDBType(dbAlias)
	if err != nil {
		return err
	}
	for _, r := range req.Requests {
		switch r.Type {
		case string(model.Create):
			v := &model.CreateRequest{Document: r.Document, Operation: r.Operation}
			if err := schemaHelpers.ValidateCreateOperation(ctx, dbAlias, dbType, r.Col, m.schemaDoc, v); err != nil {
				return err
			}
			r.Document = v.Document
			r.Operation = v.Operation
		case string(model.Update):
			if err := schemaHelpers.ValidateUpdateOperation(ctx, dbAlias, dbType, r.Col, r.Operation, r.Update, r.Find, m.schemaDoc); err != nil {
				return err
			}
		}
	}

	params.Payload = req
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return err
		}

		// Gracefully return
		return nil
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return err
	}

	// Perform the batch operation
	counts, err := crud.Batch(ctx, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		for i, r := range req.Requests {
			m.metricHook(m.project, dbAlias, r.Col, counts[i], model.OperationType(r.Type))
		}
	}

	return err
}

// DescribeTable performs a db operation for describing a table
func (m *Module) DescribeTable(ctx context.Context, dbAlias, col string) ([]model.InspectorFieldType, []model.IndexType, error) {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return nil, nil, err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return nil, nil, err
	}

	return crud.DescribeTable(ctx, col)
}

// RawBatch performs a db operation for schema creation
func (m *Module) RawBatch(ctx context.Context, dbAlias string, batchedQueries []string) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return err
	}

	return crud.RawBatch(ctx, batchedQueries)
}

// GetCollections returns collection / tables name of specified database
func (m *Module) GetCollections(ctx context.Context, dbAlias string) ([]utils.DatabaseCollections, error) {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return nil, err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return nil, err
	}

	return crud.GetCollections(ctx)
}

// GetConnectionState gets the current state of client
func (m *Module) GetConnectionState(ctx context.Context, dbAlias string) bool {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return false
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return false
	}

	return crud.GetConnectionState(ctx)
}

// DeleteTable drop specified table from database
func (m *Module) DeleteTable(ctx context.Context, dbAlias, col string) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(ctx); err != nil {
		return err
	}

	return crud.DeleteCollection(ctx, col)
}

// IsPreparedQueryPresent checks if id exist
func (m *Module) IsPreparedQueryPresent(dbAlias, id string) bool {
	m.RLock()
	defer m.RUnlock()
	_, p := m.queries[getPreparedQueryKey(dbAlias, id)]
	return p
}

// GetSchema function gets schema
func (m *Module) GetSchema(dbAlias, col string) (model.Fields, bool) {
	m.RLock()
	defer m.RUnlock()

	dbSchema, p := m.schemaDoc[dbAlias]
	if !p {
		return nil, false
	}

	colSchema, p := dbSchema[col]
	if !p {
		return nil, false
	}

	fields := make(model.Fields, len(colSchema))
	for k, v := range colSchema {
		fields[k] = v
	}

	return fields, true
}
