package syncman

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type mockAdminSyncmanInterface struct {
	mock.Mock
}

func (m *mockAdminSyncmanInterface) IsTokenValid(token string) error {
	c := m.Called(token)
	return c.Error(0)
}

func (m *mockAdminSyncmanInterface) GetInternalAccessToken() (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

func (m *mockAdminSyncmanInterface) ValidateSyncOperation(c *config.Config, project *config.Project) bool {
	a := m.Called(c, project)
	return a.Bool(0)
}

func (m *mockAdminSyncmanInterface) SetConfig(admin *config.Admin) {
	m.Called(admin)
}

type mockModulesInterface struct {
	mock.Mock
}

func (m *mockModulesInterface) SetProjectConfig(config *config.Config, le *letsencrypt.LetsEncrypt, ingressRouting *routing.Routing) {
	m.Called(config, le, ingressRouting)
}

func (m *mockModulesInterface) SetGlobalConfig(projectID string, secrets []*config.Secret, aesKey string) error {
	c := m.Called(projectID, secrets, aesKey)
	return c.Error(0)
}

func (m *mockModulesInterface) SetCrudConfig(projectID string, crudConfig config.Crud) error {
	c := m.Called(projectID, crudConfig)
	return c.Error(0)
}

func (m *mockModulesInterface) SetServicesConfig(projectID string, services *config.ServicesModule) error {
	c := m.Called(projectID, services)
	return c.Error(0)
}

func (m *mockModulesInterface) SetFileStoreConfig(projectID string, fileStore *config.FileStore) error {
	c := m.Called(projectID, fileStore)
	return c.Error(0)
}

func (m *mockModulesInterface) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	c := m.Called(projectID, eventingConfig)
	return c.Error(0)
}

func (m *mockModulesInterface) SetUsermanConfig(projectID string, auth config.Auth) {
	m.Called(projectID, auth)
}

func (m *mockModulesInterface) GetSchemaModuleForSyncMan() model.SchemaEventingInterface {
	c := m.Called()
	return c.Get(0).(model.SchemaEventingInterface)
}

type mockStoreInterface struct {
	mock.Mock
}

func (m *mockStoreInterface) WatchProjects(cb func(projects []*config.Project)) error {
	c := m.Called(cb)
	return c.Error(0)
}

func (m *mockStoreInterface) WatchServices(cb func(projects scServices)) error {
	c := m.Called(cb)
	return c.Error(0)
}

func (m *mockStoreInterface) Register() {
	m.Called()
}

func (m *mockStoreInterface) SetProject(ctx context.Context, project *config.Project) error {
	c := m.Called(ctx, project)
	return c.Error(0)
}

func (m *mockStoreInterface) DeleteProject(ctx context.Context, projectID string) error {
	c := m.Called(ctx, projectID)
	return c.Error(0)
}

func (m *mockStoreInterface) SetAdminConfig(ctx context.Context, adminConfig *config.Admin) error {
	c := m.Called(ctx, adminConfig)
	return c.Error(0)
}

func (m *mockStoreInterface) WatchAdminConfig(cb func(clusters []*config.Admin)) error {
	c := m.Called(cb)
	return c.Error(0)
}
