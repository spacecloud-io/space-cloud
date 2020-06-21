package eventing

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
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

	stopChan chan struct{}
	wg       sync.WaitGroup
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
	}

	// Start the internal processes
	m.wg.Add(1)
	go m.routineProcessIntents()
	m.wg.Add(1)
	go m.routineProcessStaged()

	return m
}

// SetConfig sets the module config
func (m *Module) SetConfig(project string, eventing *config.Eventing) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.stopChan = make(chan struct{})
	m.wg = sync.WaitGroup{}
	if eventing == nil || !eventing.Enabled {
		m.config.Enabled = false
		return nil
	}

	for eventType, schemaObj := range eventing.Schemas {
		dummyCrud := config.Crud{
			"dummyDBName": &config.CrudStub{
				Collections: map[string]*config.TableRule{
					eventType: &config.TableRule{
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
	// err := m.CloseConfig(project, eventing)
	// if err != nil {
	// 	fmt.Printf("setconfig\n")
	// }
	return nil
}

// CloseConfig closes the module config
func (m *Module) CloseConfig(project string, eventing *config.Eventing) error {
	m.eventChanMap = sync.Map{}
	m.processingEvents = sync.Map{}
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
	for k := range eventing.Rules {
		delete(eventing.Rules, k)
	}
	for k := range eventing.InternalRules {
		delete(eventing.InternalRules, k)
	}
	for k := range eventing.SecurityRules {
		delete(eventing.SecurityRules, k)
	}
	for k := range eventing.Schemas {
		delete(eventing.Schemas, k)
	}
	//m.stopChan = make(chan struct{})
	fmt.Println("main: telling goroutines to stop")
	close(m.stopChan)
	fmt.Println("main: telling goroutines to stop")
	m.wg.Wait()
	fmt.Println("main: all goroutines have told us they've finished")
	return nil
}
