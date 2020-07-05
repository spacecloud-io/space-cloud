package eventing

import (
	"errors"
	"sync"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Module is responsible for managing the eventing system
type Module struct {
	lock sync.RWMutex

	// Configurable variables
	project string
	config  *config.Eventing

	// Atomic maps to handle events being processed
	processingEvents sync.Map

	// Variables defined during initialisation
	auth   model.AuthEventingInterface
	crud   model.CrudEventingInterface
	schema model.SchemaEventingInterface

	adminMan  model.AdminEventingInterface
	syncMan   model.SyncmanEventingInterface
	fileStore model.FilestoreEventingInterface

	schemas    map[string]model.Fields
	metricHook model.MetricEventingHook
	// stores mapping of batchID w.r.t channel for sending synchronous event response
	eventChanMap sync.Map // key here is batchID
	ticker       *time.Ticker
	done         chan bool
}

// synchronous event response
type eventResponse struct {
	time     time.Time
	response chan interface{}
}

// New creates a new instance of the eventing module
func New(auth model.AuthEventingInterface, crud model.CrudEventingInterface, schemaModule model.SchemaEventingInterface, adminMan *admin.Manager, syncMan *syncman.Manager, file model.FilestoreEventingInterface, hook model.MetricEventingHook) *Module {

	m := &Module{
		auth:       auth,
		crud:       crud,
		schema:     schemaModule,
		adminMan:   adminMan,
		syncMan:    syncMan,
		schemas:    map[string]model.Fields{},
		fileStore:  file,
		metricHook: hook,
		config:     &config.Eventing{Enabled: false, InternalRules: map[string]config.EventingRule{}},
		ticker:     &time.Ticker{},
		done:       make(chan bool),
	}

	// Start the internal processes
	m.ticker = time.NewTicker(5 * time.Second)
	go m.routineProcessIntents()
	go m.routineProcessStaged()

	return m
}

// SetConfig sets the module config
func (m *Module) SetConfig(project string, eventing *config.Eventing) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if eventing == nil || !eventing.Enabled {
		m.config.Enabled = false
		return nil
	}

	for eventType, schemaObj := range eventing.Schemas {
		dummyCrud := config.Crud{
			"dummyDBName": &config.CrudStub{
				Collections: map[string]*config.TableRule{
					eventType: {
						Schema: schemaObj.Schema,
					},
				},
			},
		}

		schemaType, err := m.schema.Parser(dummyCrud)
		if err != nil {
			return err
		}
		if len(schemaType["dummyDBName"][eventType]) != 0 {
			m.schemas[eventType] = schemaType["dummyDBName"][eventType]
		}
	}

	if eventing.DBAlias == "" {
		return errors.New("invalid eventing config provided")
	}

	m.project = project
	m.config.Enabled = eventing.Enabled
	m.config.DBAlias = eventing.DBAlias

	m.config.Rules = eventing.Rules
	if m.config.Rules == nil {
		m.config.Rules = map[string]config.EventingRule{}
	}

	m.config.SecurityRules = eventing.SecurityRules
	if m.config.SecurityRules == nil {
		m.config.SecurityRules = map[string]*config.Rule{}
	}

	// `m.config.InternalRules` cannot be set by the eventing module. Its used by other modules only.
	if m.config.InternalRules == nil {
		m.config.InternalRules = map[string]config.EventingRule{}
	}

	return nil
}

// CloseConfig closes the module config
func (m *Module) CloseConfig() error {
	//erase map
	m.eventChanMap.Range(func(key interface{}, value interface{}) bool {
		m.eventChanMap.Delete(key)
		return true
	})
	//erase map
	m.processingEvents.Range(func(key interface{}, value interface{}) bool {
		m.processingEvents.Delete(key)
		return true
	})

	for k := range m.schemas {
		delete(m.schemas, k)
	}
	for k := range m.config.Rules {
		delete(m.config.Rules, k)
	}
	for k := range m.config.InternalRules {
		delete(m.config.InternalRules, k)
	}
	for k := range m.config.SecurityRules {
		delete(m.config.SecurityRules, k)
	}
	for k := range m.config.Schemas {
		delete(m.config.Schemas, k)
	}
	m.ticker.Stop()
	close(m.done)
	return nil
}
