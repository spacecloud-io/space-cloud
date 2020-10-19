package realtime

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/metrics"
)

// Module is responsible for managing the realtime module
type Module struct {
	sync.RWMutex

	// The static configuration required by the realtime module
	nodeID string
	groups sync.Map

	dbConfigs config.DatabaseConfigs
	dbRules   config.DatabaseRules
	dbSchemas config.DatabaseSchemas

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
func (m *Module) SetConfig(project string, dbConfigs config.DatabaseConfigs, dbRules config.DatabaseRules, dbSchemas config.DatabaseSchemas) error {
	m.Lock()
	defer m.Unlock()

	// Store the project id
	m.project = project
	m.dbConfigs = dbConfigs
	m.dbRules = dbRules
	m.dbSchemas = dbSchemas

	url := m.syncMan.GetRealtimeURL(m.project)

	// add the rules to the eventing module
	m.eventing.SetRealtimeTriggers(generateEventRules(dbConfigs, dbRules, dbSchemas, project, url))

	return nil
}

// SetDatabaseConfig sets database config of realtime
func (m *Module) SetDatabaseConfig(dbConfigs config.DatabaseConfigs) {
	m.Lock()
	defer m.Unlock()
	m.dbConfigs = dbConfigs
	url := m.syncMan.GetRealtimeURL(m.project)
	m.eventing.SetRealtimeTriggers(generateEventRules(m.dbConfigs, m.dbRules, m.dbSchemas, m.project, url))
}

// SetDatabaseRules sets database rules config of realtime
func (m *Module) SetDatabaseRules(databaseRules config.DatabaseRules) {
	m.Lock()
	defer m.Unlock()
	m.dbRules = databaseRules
	url := m.syncMan.GetRealtimeURL(m.project)
	m.eventing.SetRealtimeTriggers(generateEventRules(m.dbConfigs, m.dbRules, m.dbSchemas, m.project, url))
}

// SetDatabaseSchemas sets database schemas config of realtime
func (m *Module) SetDatabaseSchemas(databaseSchemas config.DatabaseSchemas) {
	m.Lock()
	defer m.Unlock()
	m.dbSchemas = databaseSchemas
	url := m.syncMan.GetRealtimeURL(m.project)
	m.eventing.SetRealtimeTriggers(generateEventRules(m.dbConfigs, m.dbRules, m.dbSchemas, m.project, url))
}

// CloseConfig close the rules and secret key required by the realtime block
func (m *Module) CloseConfig() error {
	m.Lock()
	defer m.Unlock()
	// erase map
	m.groups.Range(func(key interface{}, value interface{}) bool {
		m.groups.Delete(key)
		return true
	})
	return nil
}
