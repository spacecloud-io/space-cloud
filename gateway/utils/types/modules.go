package types

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
)

// ModulesInterface is an interface consisting of functions of the modules module used by syncman
type ModulesInterface interface {
	// SetProjectConfig sets the config all modules
	SetProjectConfig(config *config.Config, le *letsencrypt.LetsEncrypt, ingressRouting *routing.Routing)
	// SetGlobalConfig sets the auth secret and AESKey
	SetGlobalConfig(projectID string, secrets []*config.Secret, aesKey string) error
	// SetCrudConfig sets the config of crud, auth, schema and realtime modules
	SetCrudConfig(projectID string, crudConfig config.Crud) error
	// SetServicesConfig sets the config of auth and functions modules
	SetServicesConfig(projectID string, services *config.ServicesModule) error
	// SetFileStoreConfig sets the config of auth and filestore modules
	SetFileStoreConfig(projectID string, fileStore *config.FileStore) error
	// SetEventingConfig sets the config of eventing module
	SetEventingConfig(projectID string, eventingConfig *config.Eventing) error
	// SetUsermanConfig set the config of the userman module
	SetUsermanConfig(projectID string, auth config.Auth)

	// Getters

	GetSchemaModuleForSyncMan() model.SchemaEventingInterface
}
