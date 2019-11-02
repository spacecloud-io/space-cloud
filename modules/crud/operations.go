package crud

import (
	"context"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Create inserts a document (or multiple when op is "all") into the database based on dbType
func (m *Module) Create(ctx context.Context, dbType, project, col string, req *model.CreateRequest) error {
	m.RLock()
	defer m.RUnlock()
	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Invoke the create intent hook
	intent, err := m.hooks.Create(ctx, dbType, col, req)
	if err != nil {
		return err
	}

	// Perform the create operation
	n, err := crud.Create(ctx, project, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbType, col, n, utils.Create)
	}

	// Invoke the stage hook
	m.hooks.Stage(ctx, intent, err)
	return err
}

// Read returns the document(s) which match a query from the database based on dbType
func (m *Module) Read(ctx context.Context, dbType, project, col string, req *model.ReadRequest) (interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return nil, err
	}

	if err := crud.IsClientSafe(); err != nil {
		return nil, err
	}

	n, result, err := crud.Read(ctx, project, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbType, col, n, utils.Read)
	}

	return result, err
}

// Update updates the document(s) which match a query from the database based on dbType
func (m *Module) Update(ctx context.Context, dbType, project, col string, req *model.UpdateRequest) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Invoke the update intent hook
	intent, err := m.hooks.Update(ctx, dbType, col, req)
	if err != nil {
		return err
	}

	// Perform the update operation
	n, err := crud.Update(ctx, project, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbType, col, n, utils.Update)
	}

	// Invoke the stage hook
	m.hooks.Stage(ctx, intent, err)
	return err
}

// Delete removes the document(s) which match a query from the database based on dbType
func (m *Module) Delete(ctx context.Context, dbType, project, col string, req *model.DeleteRequest) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Invoke the delete intent hook
	intent, err := m.hooks.Delete(ctx, dbType, col, req)
	if err != nil {
		return err
	}

	// Perform the delete operation
	n, err := crud.Delete(ctx, project, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbType, col, n, utils.Delete)
	}

	// Invoke the stage hook
	m.hooks.Stage(ctx, intent, err)
	return err
}

// Aggregate performs an aggregation defined via the pipeline
func (m *Module) Aggregate(ctx context.Context, dbType, project, col string, req *model.AggregateRequest) (interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return nil, err
	}

	if err := crud.IsClientSafe(); err != nil {
		return nil, err
	}

	return crud.Aggregate(ctx, project, col, req)
}

// Batch performs a batch operation on the database
func (m *Module) Batch(ctx context.Context, dbType, project string, req *model.BatchRequest) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Invoke the batch intent hook
	intent, err := m.hooks.Batch(ctx, dbType, req)
	if err != nil {
		return err
	}

	// Perform the batch operation
	counts, err := crud.Batch(ctx, project, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		for i, r := range req.Requests {
			m.metricHook(m.project, dbType, r.Col, counts[i], utils.OperationType(r.Operation))
		}
	}

	// Invoke the stage hook
	m.hooks.Stage(ctx, intent, err)
	return err
}

// DescribeTable performs a db operation for describing a table
func (m *Module) DescribeTable(ctx context.Context, dbType, project, col string) ([]utils.FieldType, []utils.ForeignKeysType, error) {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return nil, nil, err
	}

	if err := crud.IsClientSafe(); err != nil {
		return nil, nil, err
	}

	return crud.DescribeTable(ctx, project, dbType, col)
}

// RawBatch performs a db operaion for schema creation
func (m *Module) RawBatch(ctx context.Context, dbType string, batchedQueries []string) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	return crud.RawBatch(ctx, batchedQueries)
}

// GetCollections returns collection / tables name of specified database
func (m *Module) GetCollections(ctx context.Context, project, dbType string) ([]utils.DatabaseCollections, error) {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return nil, err
	}

	if err := crud.IsClientSafe(); err != nil {
		return nil, err
	}

	return crud.GetCollections(ctx, project)
}

// GetCollections returns collection / tables name of specified database
func (m *Module) CreateProjectIfNotExists(ctx context.Context, project, dbType string) error {
	m.RLock()
	defer m.RUnlock()

	// Skip if project scope is disabled
	if m.removeProjectScope {
		return nil
	}

	var sql string
	switch utils.DBType(dbType) {
	case utils.MySQL:
		sql = "create database if not exists " + project
	case utils.Postgres:
		sql = "create schema if not exists " + project
	case utils.SqlServer:
		sql = "create schema if not exists " + project
	default:
		return nil
	}

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	return crud.RawExec(ctx, sql)
}

// GetConnectionState gets the current state of client
func (m *Module) GetConnectionState(ctx context.Context, dbType string) bool {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return false
	}

	if err := crud.IsClientSafe(); err != nil {
		return false
	}

	return crud.GetConnectionState(ctx, dbType)
}

// DeleteTable drop specified table from database
func (m *Module) DeleteTable(ctx context.Context, project, dbType, col string) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	return crud.DeleteCollection(ctx, project, col)
}
