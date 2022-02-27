package modules

import (
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore"
	"github.com/spaceuptech/space-cloud/gateway/modules/functions"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/caching"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/modules/userman"
)

// Auth returns the auth module
func (m *Modules) Auth(projectID string) (*auth.Module, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.auth, nil
}

// DB returns the auth module
func (m *Modules) DB(projectID string) (*crud.Module, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.db, nil
}

// User returns the auth module
func (m *Modules) User(projectID string) (*userman.Module, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.user, nil
}

// File returns the auth module
func (m *Modules) File(projectID string) (*filestore.Module, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.file, nil
}

// Functions returns the auth module
func (m *Modules) Functions(projectID string) (*functions.Module, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.functions, nil
}

// Realtime returns the auth module
func (m *Modules) Realtime(projectID string) (RealtimeInterface, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.realtime, nil
}

// Eventing returns the auth module
func (m *Modules) Eventing(projectID string) (*eventing.Module, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.eventing, nil
}

// GraphQL returns the auth module
func (m *Modules) GraphQL(projectID string) (GraphQLInterface, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.graphql, nil
}

// Schema returns the auth module
func (m *Modules) Schema(projectID string) (*schema.Schema, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.schema, nil
}

// GetSchemaModuleForSyncMan returns schema module for sync manager
func (m *Modules) GetSchemaModuleForSyncMan(projectID string) (model.SchemaEventingInterface, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.schema, nil
}

// GetAuthModuleForSyncMan returns auth module for sync manager
func (m *Modules) GetAuthModuleForSyncMan(projectID string) (model.AuthSyncManInterface, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.auth, nil
}

// LetsEncrypt returns the letsencrypt module
func (m *Modules) LetsEncrypt() *letsencrypt.LetsEncrypt {
	return m.GlobalMods.LetsEncrypt()
}

// Routing returns the routing module
func (m *Modules) Routing() *routing.Routing {
	return m.GlobalMods.Routing()
}

// Caching returns the caching module
func (m *Modules) Caching() *caching.Cache {
	return m.GlobalMods.Caching()
}