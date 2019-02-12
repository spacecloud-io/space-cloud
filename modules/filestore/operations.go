package filestore

import (
	"context"
	"io"

	"github.com/spaceuptech/space-cloud/model"
)

// CreateFile creates a file at the provided path
func (m *Module) CreateFile(ctx context.Context, project string, req *model.CreateFileRequest, file io.Reader) error {
	m.RLock()
	defer m.RUnlock()

	return m.CreateFile(ctx, project, req, file)
}

// CreateDir creates a directory at the provided path
func (m *Module) CreateDir(ctx context.Context, project string, req *model.CreateFileRequest) error {
	m.RLock()
	defer m.RUnlock()

	return m.CreateDir(ctx, project, req)
}

// ListDir lists the contents of the directory at the proided path
func (m *Module) ListDir(ctx context.Context, project string, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	m.RLock()
	defer m.RUnlock()

	return m.ListDir(ctx, project, req)
}

// ReadFile reads a file at the provided path
func (m *Module) ReadFile(ctx context.Context, project, path string) (*model.File, error) {
	m.RLock()
	defer m.RUnlock()

	return m.ReadFile(ctx, project, path)
}

// DeleteDir deletes a directory from the provided path
func (m *Module) DeleteDir(ctx context.Context, project, path string) error {
	m.RLock()
	defer m.RUnlock()

	return m.DeleteDir(ctx, project, path)
}

// DeleteFile deletes a file from the provided path
func (m *Module) DeleteFile(ctx context.Context, project, path string) error {
	m.RLock()
	defer m.RUnlock()

	return m.DeleteFile(ctx, project, path)
}
