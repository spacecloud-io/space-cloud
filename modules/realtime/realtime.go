package realtime

import (
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for managing the realtime module
type Module struct {
	sync.RWMutex
	feed    chan *model.FeedData
	enabled bool
	groups  sync.Map
}

// Init creates a new instance of the realtime module
func Init() *Module {
	return &Module{enabled: false}
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(conf *config.Realtime) error {
	m.Lock()
	defer m.Unlock()

	if conf == nil || !conf.Enabled {
		m.enabled = false
		if m.feed != nil {
			close(m.feed)
		}
		return nil
	}

	m.enabled = true
	m.feed = make(chan *model.FeedData, 5)
	m.initWorkers(utils.RealtimeWorkerCount)
	// TODO: initialise kafka client
	return nil
}

// Send broadcasts a realtime datapoint to the concerned clients
func (m *Module) Send(data *model.FeedData) {
	m.RLock()
	defer m.RUnlock()

	if m.enabled {
		m.feed <- data
	}
}
