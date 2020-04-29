package modules

import (
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore"
	"github.com/spaceuptech/space-cloud/gateway/modules/functions"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/modules/userman"
)

// Auth returns the auth module
func (m *Modules) Auth() *auth.Module {
	return m.auth
}

// DB returns the auth module
func (m *Modules) DB() *crud.Module {
	return m.db
}

// User returns the auth module
func (m *Modules) User() *userman.Module {
	return m.user
}

// File returns the auth module
func (m *Modules) File() *filestore.Module {
	return m.file
}

// Functions returns the auth module
func (m *Modules) Functions() *functions.Module {
	return m.functions
}

// Realtime returns the auth module
func (m *Modules) Realtime() RealtimeInterface {
	return m.realtime
}

// Eventing returns the auth module
func (m *Modules) Eventing() *eventing.Module {
	return m.eventing
}

// GraphQL returns the auth module
func (m *Modules) GraphQL() GraphQLInterface {
	return m.graphql
}

// Schema returns the auth module
func (m *Modules) Schema() *schema.Schema {
	return m.schema
}

// GetSchemaModuleForSyncMan returns schema module for sync manager
func (m *Modules) GetSchemaModuleForSyncMan() model.SchemaEventingInterface {
	return m.schema
}
