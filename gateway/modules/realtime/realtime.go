package realtime

import (
	"encoding/base64"
	"os"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/pubsub"
)

// Module is responsible for managing the realtime module
type Module struct {
	sync.RWMutex

	// The static configuration required by the realtime module
	project string
	nodeID  string
	groups  sync.Map

	dbConfigs config.DatabaseConfigs
	dbRules   config.DatabaseRules
	dbSchemas config.DatabaseSchemas

	// The external module realtime depends on
	eventing model.EventingRealtimeInterface
	auth     model.AuthRealtimeInterface
	crud     model.CrudRealtimeInterface
	schema   schemaInterface
	metrics  *metrics.Module
	syncMan  *syncman.Manager

	// Pubsub client
	pubsubClient *pubsub.Module

	// Auth module
	aesKey []byte
}

// Init creates a new instance of the realtime module
func Init(projectID, nodeID string, eventing model.EventingRealtimeInterface, auth model.AuthRealtimeInterface, crud model.CrudRealtimeInterface, schema schemaInterface, metrics *metrics.Module, syncMan *syncman.Manager) (*Module, error) {
	// Create a new pubsub client
	pubsubClient, err := pubsub.New(projectID, os.Getenv("REDIS_CONN"))
	if err != nil {
		return nil, err
	}

	m := &Module{project: projectID, nodeID: nodeID, syncMan: syncMan,
		eventing: eventing, auth: auth, crud: crud, schema: schema, metrics: metrics, pubsubClient: pubsubClient}

	// Start the internal routines
	go m.routineHandleMessages()

	return m, nil
}

// SetConfig set the rules and secret key required by the realtime block
func (m *Module) SetConfig(dbConfigs config.DatabaseConfigs, dbRules config.DatabaseRules, dbSchemas config.DatabaseSchemas) error {
	m.Lock()
	defer m.Unlock()

	// Store the project id
	m.dbConfigs = dbConfigs
	m.dbRules = dbRules
	m.dbSchemas = dbSchemas

	// Add the rules to the eventing module
	url := m.syncMan.GetRealtimeURL(m.project)
	m.eventing.SetRealtimeTriggers(generateEventRules(dbConfigs, dbRules, dbSchemas, m.project, url))

	return nil
}

// SetDatabaseConfig sets database config of realtime
func (m *Module) SetDatabaseConfig(dbConfigs config.DatabaseConfigs) {
	m.Lock()
	defer m.Unlock()
	m.dbConfigs = dbConfigs

	// Add the rules to the eventing module
	url := m.syncMan.GetRealtimeURL(m.project)
	m.eventing.SetRealtimeTriggers(generateEventRules(m.dbConfigs, m.dbRules, m.dbSchemas, m.project, url))
}

// SetDatabaseRules sets database rules config of realtime
func (m *Module) SetDatabaseRules(databaseRules config.DatabaseRules) {
	m.Lock()
	defer m.Unlock()
	m.dbRules = databaseRules

	// Add the rules to the eventing module
	url := m.syncMan.GetRealtimeURL(m.project)
	m.eventing.SetRealtimeTriggers(generateEventRules(m.dbConfigs, m.dbRules, m.dbSchemas, m.project, url))
}

// SetDatabaseSchemas sets database schemas config of realtime
func (m *Module) SetDatabaseSchemas(databaseSchemas config.DatabaseSchemas) {
	m.Lock()
	defer m.Unlock()
	m.dbSchemas = databaseSchemas

	// Add the rules to the eventing module
	url := m.syncMan.GetRealtimeURL(m.project)
	m.eventing.SetRealtimeTriggers(generateEventRules(m.dbConfigs, m.dbRules, m.dbSchemas, m.project, url))
}

// SetProjectAESKey set aes key
func (m *Module) SetProjectAESKey(aesKey string) error {
	m.Lock()
	defer m.Unlock()

	decodedAESKey, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return err
	}
	m.aesKey = decodedAESKey
	return nil
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
