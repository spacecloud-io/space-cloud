package driver

import (
	"context"
	"io"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateProject creates project
func (m *Module) CreateProject(ctx context.Context, project *model.Project) error {
	return m.driver.CreateProject(ctx, project)
}

// DeleteProject deletes project
func (m *Module) DeleteProject(ctx context.Context, projectID string) error {
	return m.driver.DeleteProject(ctx, projectID)
}

// ApplyService applies service
func (m *Module) ApplyService(ctx context.Context, service *model.Service) error {
	err := m.driver.ApplyService(ctx, service)
	if err == nil {
		m.metricHook(service.ProjectID)
	}
	return err
}

// GetLogs get logs of specified service
func (m *Module) GetLogs(ctx context.Context, isFollow bool, projectID, taskID, replica string) (io.ReadCloser, error) {
	return m.driver.GetLogs(ctx, isFollow, projectID, taskID, replica)
}

// GetServices gets services
func (m *Module) GetServices(ctx context.Context, projectID string) ([]*model.Service, error) {
	return m.driver.GetServices(ctx, projectID)
}

// GetServiceStatus gets services status
func (m *Module) GetServiceStatus(ctx context.Context, projectID string) ([]*model.ServiceStatus, error) {
	return m.driver.GetServiceStatus(ctx, projectID)
}

// DeleteService delete's service
func (m *Module) DeleteService(ctx context.Context, projectID, serviceID, version string) error {
	return m.driver.DeleteService(ctx, projectID, serviceID, version)
}

// ScaleUp scales a service up from 0 replicas
func (m *Module) ScaleUp(ctx context.Context, projectID, serviceID, version string) error {
	return m.driver.ScaleUp(ctx, projectID, serviceID, version)
}

// WaitForService waits for service
func (m *Module) WaitForService(ctx context.Context, service *model.Service) error {
	return m.driver.WaitForService(ctx, service)
}

// Type gets driver type
func (m *Module) Type() model.DriverType {
	return m.driver.Type()
}

// ApplyServiceRoutes applies service routes
func (m *Module) ApplyServiceRoutes(ctx context.Context, projectID, serviceID string, routes model.Routes) error {
	return m.driver.ApplyServiceRoutes(ctx, projectID, serviceID, routes)
}

// GetServiceRoutes get's service routes
func (m *Module) GetServiceRoutes(ctx context.Context, projectID string) (map[string]model.Routes, error) {
	return m.driver.GetServiceRoutes(ctx, projectID)
}

// CreateSecret create's secret
func (m *Module) CreateSecret(ctx context.Context, projectID string, secretObj *model.Secret) error {
	return m.driver.CreateSecret(ctx, projectID, secretObj)
}

// ListSecrets list's secrets
func (m *Module) ListSecrets(ctx context.Context, projectID string) ([]*model.Secret, error) {
	return m.driver.ListSecrets(ctx, projectID)
}

// DeleteSecret delete's secret
func (m *Module) DeleteSecret(ctx context.Context, projectID, secretName string) error {
	return m.driver.DeleteSecret(ctx, projectID, secretName)
}

// SetKey set's key for secret
func (m *Module) SetKey(ctx context.Context, projectID, secretName, secretKey string, secretObj *model.SecretValue) error {
	return m.driver.SetKey(ctx, projectID, secretName, secretKey, secretObj)
}

// DeleteKey delete's key of secret
func (m *Module) DeleteKey(ctx context.Context, projectID, secretName, secretKey string) error {
	return m.driver.DeleteKey(ctx, projectID, secretName, secretKey)
}

// SetFileSecretRootPath set's file secret root path
func (m *Module) SetFileSecretRootPath(ctx context.Context, projectID string, secretName, rootPath string) error {
	return m.driver.SetFileSecretRootPath(ctx, projectID, secretName, rootPath)
}
