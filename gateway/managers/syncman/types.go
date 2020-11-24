package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

// AdminSyncmanInterface is an interface consisting of functions of admin module used by eventing module
type AdminSyncmanInterface interface {
	GetInternalAccessToken() (string, error)
	IsTokenValid(ctx context.Context, token, resource, op string, attr map[string]string) (model.RequestParams, error)
	IsRegistered() bool
	GetSessionID() (string, error)
	SetServices(eventType string, services model.ScServices)
	RenewLicense(bool) error
	ValidateProjectSyncOperation(c *config.Config, project *config.ProjectConfig) bool
	SetConfig(admin *config.License) error
	GetConfig() *config.License
	SetIntegrationConfig(integrations config.Integrations)

	// For integrations
	GetIntegrationToken(id string) (string, error)
	ParseLicense(license string) (map[string]interface{}, error)
	ValidateIntegrationSyncOperation(integrations config.Integrations) error
}

type integrationInterface interface {
	SetConfig(integrations config.Integrations, integrationHooks config.IntegrationHooks) error
	SetIntegrations(integrations config.Integrations) error
	SetIntegrationHooks(integrationHooks config.IntegrationHooks)
	InvokeHook(context.Context, model.RequestParams) config.IntegrationAuthResponse
}

// ModulesInterface is an interface consisting of functions of the modules module used by syncman
type ModulesInterface interface {
	// SetInitialProjectConfig sets the config all modules
	SetInitialProjectConfig(ctx context.Context, config config.Projects) error

	// SetProjectConfig sets specific project config
	SetProjectConfig(ctx context.Context, config *config.ProjectConfig) error

	// SetDatabaseConfig sets the config of crud, auth, schema and realtime modules
	SetDatabaseConfig(ctx context.Context, projectID string, databaseConfigs config.DatabaseConfigs, schemaConfigs config.DatabaseSchemas, ruleConfigs config.DatabaseRules, prepConfigs config.DatabasePreparedQueries) error
	SetDatabaseSchemaConfig(ctx context.Context, projectID string, schemaConfigs config.DatabaseSchemas) error
	SetDatabaseRulesConfig(ctx context.Context, projectID string, ruleConfigs config.DatabaseRules) error
	SetDatabasePreparedQueryConfig(ctx context.Context, projectID string, prepConfigs config.DatabasePreparedQueries) error

	// SetFileStoreConfig sets the config of auth and filestore modules
	SetFileStoreConfig(ctx context.Context, projectID string, fileStore *config.FileStoreConfig) error
	SetFileStoreSecurityRuleConfig(ctx context.Context, projectID string, fileRule config.FileStoreRules) error

	// SetServicesConfig sets the config of auth and functions modules
	SetRemoteServiceConfig(ctx context.Context, projectID string, services config.Services) error

	SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error

	SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error
	SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error

	// SetEventingConfig sets the config of eventing module
	SetEventingConfig(ctx context.Context, projectID string, eventingConfig *config.EventingConfig, secureObj config.EventingRules, eventingSchemas config.EventingSchemas, eventingTriggers config.EventingTriggers) error
	SetEventingSchemaConfig(ctx context.Context, projectID string, schemaObj config.EventingSchemas) error
	SetEventingTriggerConfig(ctx context.Context, projectID string, triggerObj config.EventingTriggers) error
	SetEventingRuleConfig(ctx context.Context, projectID string, secureObj config.EventingRules) error

	// SetUsermanConfig set the config of the userman module
	SetUsermanConfig(ctx context.Context, projectID string, auth config.Auths) error

	// Getters
	GetSchemaModuleForSyncMan(projectID string) (model.SchemaEventingInterface, error)
	GetAuthModuleForSyncMan(projectID string) (model.AuthSyncManInterface, error)
	LetsEncrypt() *letsencrypt.LetsEncrypt
	Routing() *routing.Routing

	// Delete
	Delete(projectID string)
}

// GlobalModulesInterface is an interface consisting of functions of the global modules
type GlobalModulesInterface interface {
	// SetMetricsConfig set the config of the metrics module
	SetMetricsConfig(isMetricsEnabled bool)
}
