package syncman

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

type mockAdminSyncmanInterface struct {
	mock.Mock
}

func (m *mockAdminSyncmanInterface) IsRegistered() bool {
	panic("implement me")
}

func (m *mockAdminSyncmanInterface) GetSessionID() string {
	panic("implement me")
}

func (m *mockAdminSyncmanInterface) RenewLicense(b bool) error {
	panic("implement me")
}

func (m *mockAdminSyncmanInterface) ValidateProjectSyncOperation(projects *config.Config, projectID *config.Project) bool {
	return m.Called(projects, projectID).Bool(0)
}

func (m *mockAdminSyncmanInterface) SetConfig(admin *config.Admin, isFirst bool) error {
	panic("implement me")
}

func (m *mockAdminSyncmanInterface) IsTokenValid(token, resource, op string, attr map[string]string) (model.RequestParams, error) {
	c := m.Called(token, resource, op, attr)
	return c.Get(0).(model.RequestParams), c.Error(1)
}

func (m *mockAdminSyncmanInterface) GetInternalAccessToken() (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

func (m *mockAdminSyncmanInterface) GetConfig() *config.Admin {
	return m.Called().Get(0).(*config.Admin)
}

type mockModulesInterface struct {
	mock.Mock
}

func (m *mockModulesInterface) GetAuthModuleForSyncMan(projectID string) (model.AuthSyncManInterface, error) {
	panic("implement me")
}

func (m *mockModulesInterface) LetsEncrypt() *letsencrypt.LetsEncrypt {
	panic("implement me")
}

func (m *mockModulesInterface) Routing() *routing.Routing {
	panic("implement me")
}

func (m *mockModulesInterface) Delete(projectID string) {
	m.Called(projectID)
}

func (m *mockModulesInterface) SetProjectConfig(config *config.Project) error {
	return m.Called(config).Error(0)
}

func (m *mockModulesInterface) SetGlobalConfig(projectID, secretSource string, secrets []*config.Secret, aesKey string) error {
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

func (m *mockModulesInterface) SetUsermanConfig(projectID string, auth config.Auth) error {
	return m.Called(projectID, auth).Error(0)
}

func (m *mockModulesInterface) GetSchemaModuleForSyncMan(projectID string) (model.SchemaEventingInterface, error) {
	c := m.Called(projectID)
	return c.Get(0).(model.SchemaEventingInterface), c.Error(1)
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

func (m *mockStoreInterface) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	c := m.Called(ctx)
	return c.Get(0).(*config.Admin), c.Error(1)
}

func (m *mockStoreInterface) WatchAdminConfig(cb func(clusters []*config.Admin)) error {
	c := m.Called(cb)
	return c.Error(0)
}

type mockSchemaEventingInterface struct {
	mock.Mock
}

func (m *mockSchemaEventingInterface) CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool) {
	c := m.Called(dbAlias, col, obj, isFind)
	return map[string]interface{}{}, c.Bool(1)
}

func (m *mockSchemaEventingInterface) Parser(crud config.Crud) (model.Type, error) {
	c := m.Called(crud)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaValidator(col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {
	c := m.Called(col, collectionFields, doc)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, tables map[string]*config.TableRule) error {
	c := m.Called(ctx, dbAlias, logicalDBName, tables)
	return c.Error(0)
}

func (m *mockSchemaEventingInterface) SchemaInspection(ctx context.Context, dbAlias, project, col string) (string, error) {
	c := m.Called(ctx, dbAlias, project, col)
	return c.String(0), c.Error(1)
}
