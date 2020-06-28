package realtime

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
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
	eventing model.EventingRealtimeInterface
	auth     model.AuthRealtimeInterface
	crud     model.CrudRealtimeInterface
	schema   schemaInterface
	metrics  *metrics.Module
	syncMan  *syncman.Manager
}

// Init creates a new instance of the realtime module
func Init(nodeID string, eventing model.EventingRealtimeInterface, auth model.AuthRealtimeInterface, crud model.CrudRealtimeInterface, schema schemaInterface, metrics *metrics.Module, syncMan *syncman.Manager) (*Module, error) {

	m := &Module{nodeID: nodeID, syncMan: syncMan,
		eventing: eventing, auth: auth, crud: crud, schema: schema, metrics: metrics}

	return m, nil
}

// SetConfig set the rules and secret key required by the realtime block
func (m *Module) SetConfig(project string, crudConfig config.Crud) error {
	m.Lock()
	defer m.Unlock()

	// Store the project id
	m.project = project

	url := m.syncMan.GetRealtimeURL(m.project)

	// add the rules to the eventing module
	m.eventing.SetRealtimeTriggers(generateEventRules(crudConfig, project, url))

	return nil
}
