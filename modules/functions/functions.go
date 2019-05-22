package functions

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/nats-io/go-nats"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for Functions
type Module struct {
	sync.RWMutex
	nc              *nats.Conn
	enabled         bool
	services        sync.Map
	channel         chan *nats.Msg
	pendingRequests sync.Map
}

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

	// Close the nats client if exists
	if m.nc != nil {
		m.nc.Close()
	}

	// Close the channel of exists
	if m.channel != nil {
		close(m.channel)
	}

	// Conect and create a new nats client
	nc, err := nats.Connect(functions.Nats)
	if err != nil {
		return err
	}

	// Create new channel and start worker routines
	m.channel = make(chan *nats.Msg, 10)
	m.initWorkers(utils.FunctionsWorkerCount)

	m.nc = nc
	m.enabled = true
	return nil
}

// Request calls a function on the provided service
func (m *Module) Request(service string, timeout int, val interface{}) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()

	// Marshal the object into json
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	// Send request over nats
	subject := getSubjectName(service)
	msg, err := m.nc.Request(subject, data, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}

func (m *Module) isEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return m.enabled
}

func getSubjectName(service string) string {
	return "functions:" + service
}
