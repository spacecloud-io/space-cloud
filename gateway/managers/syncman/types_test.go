package syncman

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

type mockAdminSyncmanInterface struct {
	mock.Mock
}

func (m *mockAdminSyncmanInterface) IsTokenValid(ctx context.Context, token, resource, op string, attr map[string]string) (model.RequestParams, error) {
	c := m.Called(token, resource, op, attr)
	return c.Get(0).(model.RequestParams), c.Error(1)
}

func (m *mockAdminSyncmanInterface) GetInternalAccessToken() (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

func (m *mockAdminSyncmanInterface) ValidateProjectSyncOperation(c *config.Config, project *config.ProjectConfig) bool {
	a := m.Called(c, project)
	return a.Bool(0)
}

func (m *mockAdminSyncmanInterface) SetConfig(admin *config.Admin) error {
	return m.Called(admin).Error(0)
}

func (m *mockAdminSyncmanInterface) GetConfig() *config.Admin {
	return m.Called().Get(0).(*config.Admin)
}

type mockModulesInterface struct {
	mock.Mock
}

func (m *mockModulesInterface) SetSecurityFunctionConfig(ctx context.Context, projectID string, securityFunctions config.SecurityFunctions) error {
	a := m.Called(ctx, projectID, securityFunctions)
	return a.Error(0)
}

func (m *mockModulesInterface) SetInitialProjectConfig(ctx context.Context, config config.Projects) error {
	a := m.Called(ctx, config)
	return a.Error(0)
}

func (m *mockModulesInterface) SetDatabaseConfig(ctx context.Context, projectID string, databaseConfigs config.DatabaseConfigs, schemaConfigs config.DatabaseSchemas, ruleConfigs config.DatabaseRules, prepConfigs config.DatabasePreparedQueries) error {
	return m.Called(ctx, projectID, databaseConfigs).Error(0)
}

func (m *mockModulesInterface) SetDatabaseSchemaConfig(ctx context.Context, projectID string, schemaConfigs config.DatabaseSchemas) error {
	return m.Called(ctx, projectID, schemaConfigs).Error(0)
}

func (m *mockModulesInterface) SetDatabaseRulesConfig(ctx context.Context, ruleConfigs config.DatabaseRules) error {
	return m.Called(ctx, ruleConfigs).Error(0)
}

func (m *mockModulesInterface) SetDatabasePreparedQueryConfig(ctx context.Context, prepConfigs config.DatabasePreparedQueries) error {
	return m.Called(ctx, prepConfigs).Error(0)
}

func (m *mockModulesInterface) SetFileStoreConfig(ctx context.Context, projectID string, fileStore *config.FileStoreConfig) error {
	c := m.Called(ctx, projectID, fileStore)
	return c.Error(0)
}

func (m *mockModulesInterface) SetFileStoreSecurityRuleConfig(ctx context.Context, projectID string, fileRule config.FileStoreRules) {
	m.Called(ctx, projectID, fileRule)
}

func (m *mockModulesInterface) SetRemoteServiceConfig(ctx context.Context, projectID string, services config.Services) error {
	return m.Called(ctx, projectID, services).Error(0)
}

func (m *mockModulesInterface) SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error {
	return m.Called(ctx, projectID, c).Error(0)
}

func (m *mockModulesInterface) SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error {
	return m.Called(ctx, projectID, routes).Error(0)
}

func (m *mockModulesInterface) SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error {
	return m.Called(ctx, projectID, c).Error(0)
}

func (m *mockModulesInterface) SetEventingConfig(ctx context.Context, projectID string, eventingConfig *config.EventingConfig, secureObj config.EventingRules, eventingSchemas config.EventingSchemas, eventingTriggers config.EventingTriggers) error {
	c := m.Called(ctx, projectID, eventingConfig)
	return c.Error(0)
}

func (m *mockModulesInterface) SetEventingSchemaConfig(ctx context.Context, schemaObj config.EventingSchemas) error {
	return m.Called(ctx, schemaObj).Error(0)
}

func (m *mockModulesInterface) SetEventingTriggerConfig(ctx context.Context, triggerObj config.EventingTriggers) error {
	return m.Called(ctx, triggerObj).Error(0)
}

func (m *mockModulesInterface) SetEventingRuleConfig(ctx context.Context, secureObj config.EventingRules) error {
	return m.Called(ctx, secureObj).Error(0)
}

func (m *mockModulesInterface) SetUsermanConfig(ctx context.Context, projectID string, auth config.Auths) error {
	return m.Called(ctx, projectID, auth).Error(0)
}

func (m *mockModulesInterface) GetAuthModuleForSyncMan() model.AuthSyncManInterface {
	return m.Called().Get(0).(model.AuthSyncManInterface)
}

func (m *mockModulesInterface) LetsEncrypt() *letsencrypt.LetsEncrypt {
	return m.Called().Get(0).(*letsencrypt.LetsEncrypt)
}

func (m *mockModulesInterface) Routing() *routing.Routing {
	return m.Called().Get(0).(*routing.Routing)
}

func (m *mockModulesInterface) Delete(projectID string) {
	m.Called(projectID)
}

func (m *mockModulesInterface) SetProjectConfig(ctx context.Context, config *config.ProjectConfig) error {
	return m.Called(ctx, config).Error(0)
}

func (m *mockModulesInterface) SetGlobalConfig(projectID, secretSource string, secrets []*config.Secret, aesKey string) error {
	c := m.Called(projectID, secretSource, secrets, aesKey)
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

func (m *mockModulesInterface) GetSchemaModuleForSyncMan() model.SchemaEventingInterface {
	c := m.Called()
	return c.Get(0).(*mockSchemaEventingInterface)
}

type mockStoreInterface struct {
	mock.Mock
}

func (m *mockStoreInterface) GetGlobalConfig() (*config.Config, error) {
	c := m.Called()
	return c.Get(0).(*config.Config), c.Error(1)
}

func (m *mockStoreInterface) WatchResources(cb func(eventType string, resourceId string, resourceType config.Resource, resource interface{})) error {
	panic("implement me")
}

func (m *mockStoreInterface) SetResource(ctx context.Context, resourceID string, resource interface{}) error {
	return m.Called(ctx, resourceID, resource).Error(0)
}

func (m *mockStoreInterface) DeleteResource(ctx context.Context, resourceID string) error {
	return m.Called(ctx, resourceID).Error(0)
}

func (m *mockStoreInterface) GetProjectsConfig() (config.Projects, error) {
	c := m.Called()
	return c.Get(0).(config.Projects), c.Error(1)
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

func (m *mockStoreInterface) SetAdminConfig(ctx context.Context, adminConfig *config.Admin) error {
	c := m.Called(ctx, adminConfig)
	return c.Error(0)
}

func (m *mockStoreInterface) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	c := m.Called(ctx)
	return c.Get(0).(*config.Admin), c.Error(1)
}

func (m *mockStoreInterface) WatchClusterConfig(cb func(clusters []*config.Admin)) error {
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

func (m *mockSchemaEventingInterface) Parser(dbSchemas config.DatabaseSchemas) (model.Type, error) {
	c := m.Called(dbSchemas)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaValidator(ctx context.Context, dbAlias, col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {
	c := m.Called(ctx, dbAlias, col, collectionFields, doc)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, dbSchemas config.DatabaseSchemas) error {
	c := m.Called(ctx, dbAlias, logicalDBName, dbSchemas)
	return c.Error(0)
}

func (m *mockSchemaEventingInterface) SchemaInspection(ctx context.Context, dbAlias, project, col string, realSchema model.Collection) (string, error) {
	c := m.Called(ctx, dbAlias, project, col, realSchema)
	return c.String(0), c.Error(1)
}

func (m *mockSchemaEventingInterface) GetSchema(dbAlias, col string) (model.Fields, bool) {
	c := m.Called(dbAlias, col)
	return c.Get(0).(model.Fields), c.Bool(1)
}
func (m *mockSchemaEventingInterface) GetSchemaForDB(ctx context.Context, dbAlias, col, format string) ([]interface{}, error) {
	c := m.Called(ctx, dbAlias, col, format)
	return c.Get(0).([]interface{}), c.Error(1)
}
