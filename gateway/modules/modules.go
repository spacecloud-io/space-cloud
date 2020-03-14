package modules

import (
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore"
	"github.com/spaceuptech/space-cloud/gateway/modules/functions"
	"github.com/spaceuptech/space-cloud/gateway/modules/realtime"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/modules/userman"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Modules is an object that sets up the modules
type Modules struct {
	Auth      *auth.Module
	Crud      *crud.Module
	User      *userman.Module
	File      *filestore.Module
	Functions *functions.Module
	Realtime  *realtime.Module
	Eventing  *eventing.Module
	Graphql   *graphql.Module
	Schema    *schema.Schema
}

// New creates a new modules instance
func New(nodeID string, removeProjectScope bool, syncMan *syncman.Manager, adminMan *admin.Manager, metrics *metrics.Module) (*Modules, error) {

	c := crud.Init(removeProjectScope)
	s := schema.Init(c, removeProjectScope)
	c.SetSchema(s)

	a := auth.Init(nodeID, c, removeProjectScope)
	a.SetMakeHTTPRequest(syncMan.MakeHTTPRequest)

	fn := functions.Init(a, syncMan)
	f := filestore.Init(a)

	e := eventing.New(a, c, s, adminMan, syncMan, f)
	f.SetEventingModule(e)

	c.SetHooks(&model.CrudHooks{
		Create: e.HookDBCreateIntent,
		Update: e.HookDBUpdateIntent,
		Delete: e.HookDBDeleteIntent,
		Batch:  e.HookDBBatchIntent,
		Stage:  e.HookStage,
	}, metrics.AddDBOperation)

	rt, err := realtime.Init(nodeID, e, a, c, metrics, syncMan)
	if err != nil {
		return nil, err
	}

	u := userman.Init(c, a)
	graphqlMan := graphql.New(a, c, fn, s)

	return &Modules{Auth: a, Crud: c, User: u, File: f, Functions: fn, Realtime: rt, Eventing: e, Graphql: graphqlMan, Schema: s}, nil
}
