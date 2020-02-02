package realtime

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Module is responsible for managing the realtime module
type Module struct {
	sync.RWMutex

	// The static configuration required by the realtime module
	nodeID string
	groups sync.Map

	// Dynamic configuration
	project string

	// The external module realtime depends on
	eventing *eventing.Module
	auth     *auth.Module
	crud     *crud.Module
	schema   *schema.Schema
	metrics  *metrics.Module
	syncMan  *syncman.Manager
}

// Init creates a new instance of the realtime module
func Init(nodeID string, eventing *eventing.Module, auth *auth.Module, crud *crud.Module, schema *schema.Schema, metrics *metrics.Module, syncMan *syncman.Manager) (*Module, error) {

	m := &Module{nodeID: nodeID, syncMan: syncMan,
		eventing: eventing, auth: auth, crud: crud, schema: schema, metrics: metrics}

	return m, nil
}

// SendFeed is the function called whenever a data point (feed) is to be sent
type SendFeed func(*model.FeedData)

const (
	serviceName string = "sc-realtime"
	funcName    string = "handle"
)

type handlerAck struct {
	Ack bool
}

// SetConfig set the rules and secret key required by the realtime block
func (m *Module) SetConfig(project string, crudConfig config.Crud) error {
	m.Lock()
	defer m.Unlock()

	// Store the project id
	m.project = project

	url := m.syncMan.GetRealtimeUrl(m.project)

	// add the rules to the eventing module
	m.eventing.AddInternalRules(generateEventRules(crudConfig, project, url))

	return nil
}
