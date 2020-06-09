package crud

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// InternalCreate inserts a documents (or multiple when op is "all") into the database based on dbAlias.
// It does not invoke any hooks. This should only be used by the eventing module.
func (m *Module) InternalCreate(ctx context.Context, dbAlias, project, col string, req *model.CreateRequest, isIgnoreMetrics bool) error {
	m.RLock()
	defer m.RUnlock()

	// First step is to validate the create operation
	if err := m.schema.ValidateCreateOperation(dbAlias, col, req); err != nil {
		return err
	}

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	var n int64
	// Perform the create operation
	if req.IsBatch {
		n, err = m.createBatch(project, dbAlias, col, req.Document)
	} else {
		n, err = crud.Create(ctx, col, req)
	}

	// Invoke the metric hook if the operation was successful
	if err == nil && !isIgnoreMetrics {
		m.metricHook(m.project, dbAlias, col, n, utils.Create)
	}

	return err
}

// InternalUpdate updates the documents(s) which match a query from the database based on dbType.
// It does not invoke any hooks. This should only be used by the eventing module.
func (m *Module) InternalUpdate(ctx context.Context, dbAlias, project, col string, req *model.UpdateRequest) error {
	m.RLock()
	defer m.RUnlock()

	// First step is to validate the update operation
	if err := m.schema.ValidateUpdateOperation(dbAlias, col, req.Operation, req.Update, req.Find); err != nil {
		return err
	}

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Adjust where clause
	if err := m.schema.AdjustWhereClause(dbAlias, crud.GetDBType(), col, req.Find); err != nil {
		return err
	}

	// Perform the update operation
	n, err := crud.Update(ctx, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbAlias, col, n, utils.Update)
	}

	return err
}

// InternalDelete removes the documents(s) which match a query from the database based on dbType.
// It does not invoke any hooks. This should only be used by the eventing module.
func (m *Module) InternalDelete(ctx context.Context, dbAlias, project, col string, req *model.DeleteRequest) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbAlias)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Adjust where clause
	if err := m.schema.AdjustWhereClause(dbAlias, crud.GetDBType(), col, req.Find); err != nil {
		return err
	}

	// Perform the delete operation
	n, err := crud.Delete(ctx, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbAlias, col, n, utils.Update)
	}

	return err
}
