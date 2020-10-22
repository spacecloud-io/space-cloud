package modules

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/managers"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore"
	"github.com/spaceuptech/space-cloud/gateway/modules/functions"
	"github.com/spaceuptech/space-cloud/gateway/modules/global"
	"github.com/spaceuptech/space-cloud/gateway/modules/realtime"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/modules/userman"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
)

// Modules is an object that sets up the modules
type Modules struct {
	auth      *auth.Module
	db        *crud.Module
	user      *userman.Module
	file      *filestore.Module
	functions *functions.Module
	realtime  *realtime.Module
	eventing  *eventing.Module
	graphql   *graphql.Module
	schema    *schema.Schema

	// Global Modules
	GlobalMods *global.Global

	// Managers
	Managers *managers.Managers
}

// New creates a new modules instance
func New(nodeID string, managers *managers.Managers, globalMods *global.Global) (*Modules, error) {

	// Extract managers
	adminMan := managers.Admin()
	syncMan := managers.Sync()

	// Extract global modules
	metrics := globalMods.Metrics()

	c := crud.Init()
	c.SetGetSecrets(syncMan.GetSecrets)
	s := schema.Init(c)
	c.SetSchema(s)

	a := auth.Init(nodeID, c, adminMan)
	a.SetMakeHTTPRequest(syncMan.MakeHTTPRequest)
	c.SetAuth(a)

	fn := functions.Init(a, syncMan, metrics.AddFunctionOperation)
	f := filestore.Init(a, metrics.AddFileOperation)
	f.SetGetSecrets(syncMan.GetSecrets)

	e := eventing.New(a, c, s, adminMan, syncMan, f, metrics.AddEventingType)
	f.SetEventingModule(e)

	c.SetHooks(&model.CrudHooks{
		Create: e.HookDBCreateIntent,
		Update: e.HookDBUpdateIntent,
		Delete: e.HookDBDeleteIntent,
		Batch:  e.HookDBBatchIntent,
		Stage:  e.HookStage,
	}, metrics.AddDBOperation)

	rt, err := realtime.Init(nodeID, e, a, c, s, metrics, syncMan)
	if err != nil {
		return nil, err
	}

	u := userman.Init(c, a)
	graphqlMan := graphql.New(a, c, fn, s)

	return &Modules{auth: a, db: c, user: u, file: f, functions: fn, realtime: rt, eventing: e, graphql: graphqlMan, schema: s, GlobalMods: globalMods, Managers: managers}, nil
}

// Delete deletes a project
func (m *Modules) Delete(projectID string) {
	// Close all the modules here
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of db module", nil)
	if err := m.db.CloseConfig(); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing db module config", err, map[string]interface{}{"project": projectID})
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of filestore module", nil)
	if err := m.file.CloseConfig(); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing filestore module config", err, map[string]interface{}{"project": projectID})
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of eventing module", nil)
	if err := m.eventing.CloseConfig(); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing eventing module config", err, map[string]interface{}{"project": projectID})
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of realtime module", nil)
	if err := m.realtime.CloseConfig(); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing realtime module config", err, map[string]interface{}{"project": projectID})
	}
}
