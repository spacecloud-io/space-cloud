package filestore

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// CreateDir creates a directory at the provided path
func (m *Module) CreateDir(ctx context.Context, project, token string, req *model.CreateFileRequest, meta map[string]interface{}) (int, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	_, err := m.auth.IsFileOpAuthorised(ctx, project, token, req.Path, utils.FileCreate, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	intent, err := m.eventing.CreateFileIntentHook(ctx, req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err = m.store.CreateDir(req); err != nil {
		m.eventing.HookStage(ctx, intent, err)
		return http.StatusInternalServerError, nil
	}

	m.eventing.HookStage(ctx, intent, nil)
	m.metricsHook(project, string(m.store.GetStoreType()), utils.Create)
	return http.StatusOK, err
}

// DeleteFile deletes a file at the provided path
func (m *Module) DeleteFile(ctx context.Context, project, token string, path string, meta map[string]interface{}) (int, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	_, err := m.auth.IsFileOpAuthorised(ctx, project, token, path, utils.FileDelete, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	intent, err := m.eventing.DeleteFileIntentHook(ctx, path, meta)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err = m.store.DeleteFile(path); err != nil {
		m.eventing.HookStage(ctx, intent, err)
		return http.StatusInternalServerError, nil
	}

	m.eventing.HookStage(ctx, intent, nil)
	m.metricsHook(project, string(m.store.GetStoreType()), utils.Delete)
	return http.StatusOK, err
}

// ListFiles lists all the files in the provided path
func (m *Module) ListFiles(ctx context.Context, project, token string, req *model.ListFilesRequest) (int, []*model.ListFilesResponse, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, nil, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	_, err := m.auth.IsFileOpAuthorised(ctx, project, token, req.Path, utils.FileRead, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, nil, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	res, err := m.store.ListDir(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	m.metricsHook(project, string(m.store.GetStoreType()), utils.List)
	return http.StatusOK, res, nil
}

// UploadFile uploads a file to the provided path
func (m *Module) UploadFile(ctx context.Context, project, token string, req *model.CreateFileRequest, reader io.Reader) (int, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	_, err := m.auth.IsFileOpAuthorised(ctx, project, token, req.Path, utils.FileCreate, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	intent, err := m.eventing.CreateFileIntentHook(ctx, req)
	if err != nil {
		return 500, err
	}

	if err = m.store.CreateFile(req, reader); err != nil {
		m.eventing.HookStage(ctx, intent, err)
		return http.StatusInternalServerError, nil
	}

	m.eventing.HookStage(ctx, intent, nil)
	m.metricsHook(project, string(m.store.GetStoreType()), utils.Create)
	return http.StatusOK, err
}

// DownloadFile downloads a file from the provided path
func (m *Module) DownloadFile(ctx context.Context, project, token, path string) (int, *model.File, error) {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, nil, errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	_, err := m.auth.IsFileOpAuthorised(ctx, project, token, path, utils.FileRead, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, nil, errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	// Read the file from file storage
	file, err := m.store.ReadFile(path)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	m.metricsHook(project, string(m.store.GetStoreType()), utils.Read)
	return http.StatusOK, file, nil
}

// DoesExists checks if the provided path exists
func (m *Module) DoesExists(ctx context.Context, project, token, path string) error {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return errors.New("This feature isn't enabled")
	}

	// Check if the user is authorised to make this request
	_, err := m.auth.IsFileOpAuthorised(ctx, project, token, path, utils.FileRead, map[string]interface{}{})
	if err != nil {
		return errors.New("You are not authorized to make this request")
	}

	m.RLock()
	defer m.RUnlock()

	// Read the file from file storage
	return m.store.DoesExists(path)
}

// GetState checks if selected storage is active
func (m *Module) GetState(ctx context.Context) error {
	// Exit if file storage is not enabled
	if !m.IsEnabled() {
		return errors.New("This feature isn't enabled")
	}
	m.RLock()
	defer m.RUnlock()
	// Read the state from file storage
	return m.store.GetState(ctx)
}
