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
	pubSubClient *gochannel.GoChannel
}

// New intialise the pubsub connector
func New() (*Connector, error) {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	return &Connector{pubSubClient: pubSub}, nil
}

// GetPubsubClient returns pubsub client
func (c *Connector) GetPubsubClient() *gochannel.GoChannel {
	return c.pubSubClient
}

// Destruct destroys the pubsub module
func (c *Connector) Destruct() error {
	return c.pubSubClient.Close()
}
