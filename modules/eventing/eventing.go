package eventing

import (
	"errors"
	"sync"

	nats "github.com/nats-io/nats.go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for managing the eventing system
type Module struct {
	lock sync.RWMutex

	// Configurable variables
	project string
	config  *config.Eventing

	// Atomic maps to handle pending requests
	pendingCreates sync.Map
	pending        sync.Map
	pendingDelete  sync.Map

	// Variables defined during initialisation
	nc        *nats.Conn
	maxTokens int
	crud      *crud.Module
}

// New creates a new instance of the eventing module
func New(crud *crud.Module) *Module {
	return &Module{
		maxTokens: utils.MaxEventTokens,
		crud:      crud,
		config:    &config.Eventing{Enabled: false},
	}
}

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

		if _, err := m.nc.ChanSubscribe("core-eventing", channel); err != nil {
			return err
		}
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
