package functions

import (
	"errors"
	"sync"

	nats "github.com/nats-io/nats.go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for functions
type Module struct {
	sync.RWMutex
	nc              *nats.Conn
	enabled         bool
	services        sync.Map
	channel         chan *nats.Msg
	pendingRequests sync.Map
}

// SendPayload is the function called whenever a data point (payload) is to be sent
type SendPayload func(*model.FunctionsPayload)

// Init returns a new instance of the Functions module
func Init() *Module {
	m := new(Module)
	go m.removeStaleRequests()
	return m
}

// SetConfig set the config required by the Functions module
func (m *Module) SetConfig(functions *config.Functions) error {
	m.Lock()
	defer m.Unlock()

	if functions == nil || !functions.Enabled {
		m.enabled = false
		return nil
	}

	// Connect and create a new nats client
	if functions.Broker != utils.Nats {
		return errors.New("Functions Error: Broker is not supported")
	}

	if m.nc == nil {
		nc, err := nats.Connect(functions.Conn)
		if err != nil {
			return err
		}
		m.nc = nc

		// Create new channel and start worker routines
		m.channel = make(chan *nats.Msg, 10)
		m.initWorkers(utils.FunctionsWorkerCount)
	}

	m.enabled = true
	return nil
}

// IsEnabled checks if the Functions module is enabled
func (m *Module) IsEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return m.enabled
}

func getSubjectName(service string) string {
	return "functions:" + service
}
