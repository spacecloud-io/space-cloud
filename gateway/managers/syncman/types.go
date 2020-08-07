package syncman

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

// AdminSyncmanInterface is an interface consisting of functions of admin module used by eventing module
type AdminSyncmanInterface interface {
	GetInternalAccessToken() (string, error)
	IsTokenValid(token, resource, op string, attr map[string]string) (model.RequestParams, error)
	IsRegistered() bool
	GetSessionID() string
	RenewLicense(bool) error
	ValidateProjectSyncOperation(c *config.Config, project *config.Project) bool
	SetConfig(admin *config.Admin, isFirst bool) error
	GetConfig() *config.Admin
}

// ModulesInterface is an interface consisting of functions of the modules module used by syncman
type ModulesInterface interface {
	// SetProjectConfig sets the config all modules
	SetProjectConfig(config *config.Project) error
	// SetGlobalConfig sets the auth secret and AESKey
	SetGlobalConfig(projectID, secretSource string, secrets []*config.Secret, aesKey string) error
	// SetCrudConfig sets the config of crud, auth, schema and realtime modules
	SetCrudConfig(projectID string, crudConfig config.Crud) error
	// SetServicesConfig sets the config of auth and functions modules
	SetServicesConfig(projectID string, services *config.ServicesModule) error
	// SetFileStoreConfig sets the config of auth and filestore modules
	SetFileStoreConfig(projectID string, fileStore *config.FileStore) error
	// SetEventingConfig sets the config of eventing module
	SetEventingConfig(projectID string, eventingConfig *config.Eventing) error
	// SetUsermanConfig set the config of the 0]userman module
	SetUsermanConfig(projectID string, auth config.Auth) error

	// Getters
	GetSchemaModuleForSyncMan(projectID string) (model.SchemaEventingInterface, error)
	GetAuthModuleForSyncMan(projectID string) (model.AuthSyncManInterface, error)
	LetsEncrypt() *letsencrypt.LetsEncrypt
	Routing() *routing.Routing

	// Delete
	Delete(projectID string)
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
	Schema string `json:"schema"`
}
