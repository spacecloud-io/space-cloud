package connectors

import (
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// Connector implements connector for pubsub module
type Connector struct {
	lock sync.RWMutex

	// For pubsub engine
	PubSubClient *gochannel.GoChannel
}

// New intialise the pubsub connector
func New() (*Connector, error) {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	return &Connector{PubSubClient: pubSub}, nil
}

// Destruct destroys the pubsub module
func (c *Connector) Destruct() error {
	return c.PubSubClient.Close()
}
