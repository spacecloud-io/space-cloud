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
	ValidateProjectSyncOperation(c *config.Config, project *config.ProjectConfig) bool
	// SetConfig(admin *config.Admin) error
	// GetConfig() *config.Admin
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
	SetDatabaseRulesConfig(ctx context.Context, ruleConfigs config.DatabaseRules) error
	SetDatabasePreparedQueryConfig(ctx context.Context, prepConfigs config.DatabasePreparedQueries) error

	// SetFileStoreConfig sets the config of auth and filestore modules
	SetFileStoreConfig(ctx context.Context, projectID string, fileStore *config.FileStoreConfig) error
	SetFileStoreSecurityRuleConfig(ctx context.Context, projectID string, fileRule config.FileStoreRules)

	// SetServicesConfig sets the config of auth and functions modules
	SetRemoteServiceConfig(ctx context.Context, projectID string, services config.Services) error

	SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error

	SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error
	SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error

	// SetEventingConfig sets the config of eventing module
	SetEventingConfig(ctx context.Context, projectID string, eventingConfig *config.EventingConfig, secureObj config.EventingRules, eventingSchemas config.EventingSchemas, eventingTriggers config.EventingTriggers) error
	SetEventingSchemaConfig(ctx context.Context, schemaObj config.EventingSchemas) error
	SetEventingTriggerConfig(ctx context.Context, triggerObj config.EventingTriggers) error
	SetEventingRuleConfig(ctx context.Context, secureObj config.EventingRules) error

	// SetUsermanConfig set the config of the userman module
	SetUsermanConfig(ctx context.Context, projectID string, auth config.Auths) error
	SetSecurityFunctionConfig(ctx context.Context, _ string, securityFunctions config.SecurityFunctions) error
	// Getters
	GetSchemaModuleForSyncMan() model.SchemaEventingInterface
	GetAuthModuleForSyncMan() model.AuthSyncManInterface
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
