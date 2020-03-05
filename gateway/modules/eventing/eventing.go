package eventing

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

const (
	eventingLogs     string = "event_logs"
	invocationLogs   string = "invocation_logs"
	invocationSchema string = `type invocation_logs {
		_id: ID! @primary
		event_id: ID! @foreign(table: "event_logs", field: "_id")
		invocation_time: DateTime!
		request_payload: String
		response_status_code: Integer
		response_body: String
		remark: String	
	  }`
	eventSchema string = `type event_logs {
		_id: ID! @primary
		batchid: String
		type: String
		token: Integer
		timestamp: Integer
		event_timestamp: Integer
		payload: String
		status: String
		retries: Integer
		url: String
		remark: String
		invocations: [invocation_logs]! @link(table: "invocation_logs", from: "_id", to: "event_id")
	  }`
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

	adminMan  *admin.Manager
	syncMan   *syncman.Manager
	fileStore model.FilestoreEventingInterface

	schemas map[string]model.Fields

	// stores mapping of batchID w.r.t channel for sending synchronous event response
	eventChanMap sync.Map // key here is batchID
}

// synchronous event response
type eventResponse struct {
	time     time.Time
	response chan interface{}
}

// New creates a new instance of the eventing module
func New(auth model.AuthEventingInterface, crud model.CrudEventingInterface, schemaModule model.SchemaEventingInterface, adminMan *admin.Manager, syncMan *syncman.Manager, file model.FilestoreEventingInterface) *Module {

	m := &Module{
		auth:      auth,
		crud:      crud,
		schema:    schemaModule,
		adminMan:  adminMan,
		syncMan:   syncMan,
		schemas:   map[string]model.Fields{},
		fileStore: file,
		config:    &config.Eventing{Enabled: false, InternalRules: map[string]config.EventingRule{}},
	}

	// Start the internal processes
	go m.routineProcessIntents()
	go m.routineProcessStaged()

	return m
}

// SetConfig sets the module config
func (m *Module) SetConfig(project string, eventing *config.Eventing) error {
	m.lock.Lock()
	defer m.lock.Unlock()

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

	if eventing == nil || !eventing.Enabled {
		m.config.Enabled = false
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.schema.SchemaModifyAll(ctx, eventing.DBType, project, map[string]*config.TableRule{eventingLogs: &config.TableRule{Schema: eventSchema}, invocationLogs: &config.TableRule{Schema: invocationSchema}})
	if err != nil {
		logrus.Errorf("Could not create tables required for eventing module - %s", err.Error())
		return err
	}

	if eventing.DBType == "" {
		return errors.New("invalid eventing config provided")
	}

	m.project = project
	m.config = eventing

	if m.config.Rules == nil {
		m.config.Rules = map[string]config.EventingRule{}
	}

	// Reset the internal rules
	m.config.InternalRules = map[string]config.EventingRule{}

	return nil
}
