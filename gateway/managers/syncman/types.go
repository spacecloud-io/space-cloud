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
	ValidateProjectSyncOperation(c *config.Config, project *config.Project) bool
	SetConfig(admin *config.Admin) error
	GetConfig() *config.Admin
}

// ModulesInterface is an interface consisting of functions of the modules module used by syncman
type ModulesInterface interface {
	// SetProjectConfig sets the config all modules
	SetProjectConfig(config *config.Project) error
	// SetCrudConfig sets the config of crud, auth, schema and realtime modules
	SetCrudConfig(projectID string, crudConfig config.Crud) error
	// SetServicesConfig sets the config of auth and functions modules
	SetServicesConfig(projectID string, services *config.ServicesModule) error
	// SetFileStoreConfig sets the config of auth and filestore modules
	SetFileStoreConfig(projectID string, fileStore *config.FileStore) error
	// SetEventingConfig sets the config of eventing module
	SetEventingConfig(projectID string, eventingConfig *config.Eventing) error
	// SetUsermanConfig set the config of the userman module
	SetUsermanConfig(projectID string, auth config.Auth) error

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

type preparedQueryResponse struct {
	ID        string       `json:"id"`
	DBAlias   string       `json:"db"`
	SQL       string       `json:"sql"`
	Arguments []string     `json:"args"`
	Rule      *config.Rule `json:"rule"`
}

type dbRulesResponse struct {
	IsRealTimeEnabled bool                    `json:"isRealtimeEnabled"`
	Rules             map[string]*config.Rule `json:"rules"`
}

type dbSchemaResponse struct {
	DbAlias   string       `json:"dbAlias"`
	Col       string       `json:"col"`
	Schema    string       `json:"schema"`
	SchemaObj model.Fields `json:"schemaObj"`
}
