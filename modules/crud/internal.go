package crud

import (
	"context"

	"github.com/spaceuptech/space-cloud/model"
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
	return crud.Create(ctx, project, col, req)
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
	return crud.Update(ctx, project, col, req)
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

	// Perfrom the delete operation
	return crud.Delete(ctx, project, col, req)
}
