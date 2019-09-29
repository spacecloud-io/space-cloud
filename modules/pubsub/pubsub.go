package pubsub

import (
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/pubsub/nats"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for Pubsub
type Module struct {
	sync.RWMutex
	broker     string
	connection pubsubBroker
	enabled    bool
	storage    sync.Map
	auth       *auth.Module
}

// Init returns a new instance of the Pubsub module
func Init(auth *auth.Module) *Module {
	return &Module{enabled: false, auth: auth}
}

// SetConfig sets the config required by the Pubsub module
func (m *Module) SetConfig(pubsub *config.Pubsub) error {
	m.Lock()
	defer m.Unlock()

	if pubsub == nil || !pubsub.Enabled {
		m.enabled = false
		return nil
	}

	switch pubsub.Broker {
	case utils.Nats:
		m.broker = string(pubsub.Broker)
		if m.connection == nil {
			nc, err := nats.Connect(pubsub.Conn)
			if err != nil {
				return err
			}
			m.connection = nc
		}
	default:
		return errors.New("Pubsub Error: Broker is not supported")
	}

	m.enabled = true
	return nil
}

// IsEnabled checks if the Pubsub module is enabled
func (m *Module) IsEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return m.enabled
}
