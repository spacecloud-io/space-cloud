package realtime

import (
	"errors"
	"sync"
	"time"

	"github.com/nats-io/go-nats"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for managing the realtime module
type Module struct {
	sync.RWMutex
	feed            chan *nats.Msg
	enabled         bool
	project         string
	groups          sync.Map
	nc              *nats.Conn
	pendingRequests sync.Map
	crud            *crud.Module
}

// Init creates a new instance of the realtime module
func Init(crud *crud.Module) *Module {
	m := &Module{enabled: false, crud: crud}
	go m.removeStaleRequests()
	return m
}

// SendFeed is the function called whenever a data point (feed) is to be sent
type SendFeed func(*model.FeedData)

const (
	typeIntent string = "intent"
	typeAck    string = "ack"
)

// Message is the message sent over nats
type Message struct {
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Ack  bool            `json:"ack"`
	Data *model.FeedData `json:"feed"`
}

type pendingRequest struct {
	data *model.FeedData
	time time.Time
}

// SetConfig set the rules and secret key required by the crud block
func (m *Module) SetConfig(project string, conf *config.Realtime) error {
	m.Lock()
	defer m.Unlock()

	m.project = project

	if conf == nil || !conf.Enabled {
		m.enabled = false
		return nil
	}

	// Close the nats client if exists
	if m.nc != nil {
		m.nc.Close()
	}

	// Close the channel if exists
	if m.feed != nil {
		close(m.feed)
	}

	// Conect and create a new nats client
	if conf.Broker != utils.Nats {
		return errors.New("Realtime Error: Broker is not supported")
	}

	nc, err := nats.Connect(conf.Conn)
	if err != nil {
		return err
	}

	m.feed = make(chan *nats.Msg, utils.RealtimeWorkerCount)
	m.initWorkers(utils.RealtimeWorkerCount)
	m.nc = nc
	m.enabled = true
	return nil
}

func getSubjectName(project, col string) string {
	return "realtime:" + project + ":" + col
}
