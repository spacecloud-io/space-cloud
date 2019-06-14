package realtime

import (
	"sync"
	"strconv"

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

// SendFeed is the function called whenever a data point (feed) is to be sent
type SendFeed func(*model.FeedData)

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

func AcceptableIdType(id interface{}) (string, bool) {
	switch v := id.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case float64:
		// json.Unmarshal converts all numbers to float64
		if float64(int64(v)) == v { // v is actually an int
			return strconv.FormatInt(int64(v), 10), true
		}
		return "", false
	default:
		return "", false
	}
}
