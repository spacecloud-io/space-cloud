package functions

import (
	"sync"

	"github.com/nats-io/nats.go"

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
func Init() (*Module, error) {
	m := new(Module)
	go m.removeStaleRequests()

	// Create a nats client
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	m.nc = nc

	// Create new channel and start worker routines
	m.channel = make(chan *nats.Msg, 10)
	m.initWorkers(utils.FunctionsWorkerCount)

	// Enable the module
	m.enabled = true

	return m, nil
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
