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
	err = crud.Create(ctx, project, col, req)

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

	return crud.Read(ctx, project, col, req)
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
	err = crud.Update(ctx, project, col, req)

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

	// Perfrom the delete operation
	err = crud.Delete(ctx, project, col, req)

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

	// Perfrom the batch operation
	err = crud.Batch(ctx, project, req)

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

	var sql string
	switch utils.DBType(dbType) {
	case utils.MySQL:
		sql = "create database if not exists " + project
	case utils.Postgres:
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
