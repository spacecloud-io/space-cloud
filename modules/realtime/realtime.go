package realtime

import (
	"sync"

	"github.com/nats-io/nats.go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/eventing"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// Module is responsible for managing the realtime module
type Module struct {
	sync.RWMutex

	// The static configuration required by the realtime module
	nodeID string
	groups sync.Map
	ec     *nats.EncodedConn

	// Dynamic configuration that can change over time
	project string

	// The external module realtime depends on
	eventing  *eventing.Module
	auth      *auth.Module
	crud      *crud.Module
	functions *functions.Module
	adminMan  *admin.Manager
	syncMan   *syncman.SyncManager
}

// Init creates a new instance of the realtime module
func Init(nodeID string, eventing *eventing.Module, auth *auth.Module, crud *crud.Module,
	functions *functions.Module, adminMan *admin.Manager, syncMan *syncman.SyncManager) (*Module, error) {

	m := &Module{nodeID: nodeID, adminMan: adminMan, syncMan: syncMan,
		eventing: eventing, auth: auth, crud: crud, functions: functions}

	// Register the realtime service handler
	if err := m.registerEventHandlerService(); err != nil {
		return nil, err
	}

	// Create a nats connection
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	m.ec = ec

	if _, err := ec.Subscribe(pubSubTopic, m.handleRealtimeRequests); err != nil {
		return nil, err
	}

	return m, nil
}

// SendFeed is the function called whenever a data point (feed) is to be sent
type SendFeed func(*model.FeedData)

const (
	serviceName string = "realtime-handler"
	funcName    string = "handle"
	pubSubTopic string = "realtime-message"
)

type handlerAck struct {
	Ack bool
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, crudConfig config.Crud) {
	m.Lock()
	defer m.Unlock()

	// Store the project id
	m.project = project

	// add the rules to the eventing module
	m.eventing.AddInternalRules(generateEventRules(crudConfig))
	return
}
