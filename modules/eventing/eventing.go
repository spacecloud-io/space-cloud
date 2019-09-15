package eventing

import (
	"errors"
	"sync"

	nats "github.com/nats-io/nats.go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/utils/syncman"
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
	nc        *nats.Conn
	crud      *crud.Module
	functions *functions.Module
	syncMan   *syncman.SyncManager
}

// New creates a new instance of the eventing module
func New(crud *crud.Module, functions *functions.Module, syncMan *syncman.SyncManager) *Module {

	m := &Module{
		crud:      crud,
		functions: functions,
		syncMan:   syncMan,
		config:    &config.Eventing{Enabled: false},
	}

	// Start the internal processes
	go m.routineProcessIntents()
	go m.routineProcessStaged()

	return m
}

const internalEventingSubject string = "core-eventing"

// SetConfig sets the module config
func (m *Module) SetConfig(project string, eventing *config.Eventing) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.nc == nil {
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			return err
		}
		m.nc = nc
		channel := make(chan *nats.Msg, 10)

		if _, err := m.nc.ChanSubscribe(internalEventingSubject, channel); err != nil {
			return err
		}

		m.initEventWorkers(channel, 10)

		// Create new channel and start worker routines
		//m.channel = make(chan *nats.Msg, 10)
		//m.initWorkers(utils.FunctionsWorkerCount)
	}

	if eventing == nil {
		m.config.Enabled = false
		return nil
	}

	if eventing.DBType == "" || eventing.Col == "" {
		return errors.New("Invalid eventing config provided")
	}

	m.project = project
	m.config = eventing
	return nil
}
