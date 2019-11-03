package crud

import (
	"context"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// InternalCreate inserts a document (or multiple when op is "all") into the database based on dbType.
// It does not invoke any hooks. This should only be used by the eventing module.
func (m *Module) InternalCreate(ctx context.Context, dbType, project, col string, req *model.CreateRequest) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Perform the create operation
	n, err := crud.Create(ctx, project, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbType, col, n, utils.Create)
	}

	return err
}

// InternalUpdate updates the document(s) which match a query from the database based on dbType.
// It does not invoke any hooks. This should only be used by the eventing module.
func (m *Module) InternalUpdate(ctx context.Context, dbType, project, col string, req *model.UpdateRequest) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Perform the update operation
	n, err := crud.Update(ctx, project, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbType, col, n, utils.Update)
	}

	return err
}

// InternalDelete removes the document(s) which match a query from the database based on dbType.
// It does not invoke any hooks. This should only be used by the eventing module.
func (m *Module) InternalDelete(ctx context.Context, dbType, project, col string, req *model.DeleteRequest) error {
	m.RLock()
	defer m.RUnlock()

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return err
	}

	if err := crud.IsClientSafe(); err != nil {
		return err
	}

	// Perform the delete operation
	n, err := crud.Delete(ctx, project, col, req)

	// Invoke the metric hook if the operation was successful
	if err == nil {
		m.metricHook(m.project, dbType, col, n, utils.Update)
	}

	return err
}
