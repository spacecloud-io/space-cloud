package filestore

import (
	"context"
	"io"
	"net/http"
	"errors"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// CreateDir creates a directory at the provided path
func (m *Module) CreateDir(ctx context.Context, project, token string, req *model.CreateFileRequest) (int, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	err := m.auth.IsFileOpAuthorised(project, token, req.Path, utils.FileCreate, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	err = m.store.CreateDir(ctx, project, req)
	if err != nil {
		return 500, err
	} else {
		return 200, nil
	}
}

// DeleteFile deletes a file at the provided path
func (m *Module) DeleteFile(ctx context.Context, project, token, path string) (int, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	err := m.auth.IsFileOpAuthorised(project, token, path, utils.FileDelete, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	err = m.store.DeleteDir(ctx, project, path)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 200, nil
}

// ListFiles lists all the files in the provided path
func (m *Module) ListFiles(ctx context.Context, project, token string, req *model.ListFilesRequest) (int, []*model.ListFilesResponse, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, nil, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	err := m.auth.IsFileOpAuthorised(project, token, req.Path, utils.FileRead, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, nil, errors.New("You are not authorized to make this request")
	}
	
	m.RLock()
	defer m.RUnlock()

	res, err := m.store.ListDir(ctx, project, req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, res, nil
}

// UploadFile uploads a file to the provided path
func (m *Module) UploadFile(ctx context.Context, project, token string, req *model.CreateFileRequest, reader io.Reader) (int, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	err := m.auth.IsFileOpAuthorised(project, token, req.Path, utils.FileCreate, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	err = m.store.CreateFile(ctx, project, req, reader)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// DownloadFile downloads a file from the provided path
func (m *Module) DownloadFile(ctx context.Context, project, token, path string) (int, *model.File, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, nil, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	err := m.auth.IsFileOpAuthorised(project, token, path, utils.FileRead, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, nil, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	// Read the file from file storage
	file, err := m.store.ReadFile(ctx, project, path)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, file, nil
}
