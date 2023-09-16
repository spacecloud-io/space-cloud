package pubsub

import (
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/provider"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/pubsub/connectors"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var connectorPool = caddy.NewUsagePool()

func init() {
	caddy.RegisterModule(Module{})
	provider.Register("pubsub", 50)
}

// Module defines struct for pubsub app
type Module struct {
	Workspace string `json:"workspace"`

	// For pubsub engine
	pubSub *gochannel.GoChannel

	// APIs
	apis apis.APIs

	// For internal usage
	logger      *zap.Logger
	asyncapiDoc *AsyncAPI
	channels    []v1alpha1.PubsubChannelSpec
}

// CaddyModule returns the Caddy module information.
func (Module) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "provider.pubsub",
		New: func() caddy.Module { return new(Module) },
	}
}

// Provision sets up the pubsub module.
func (m *Module) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger(m)

	poolKey := getPoolKey()
	val, _, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
		return connectors.New()
	})
	if err != nil {
		m.logger.Error("Unable to open pubsub connector", zap.Error(err))
		return err
	}
	connector := val.(*connectors.Connector)
	m.pubSub = connector.GetPubsubClient()

	// Get all space-cloud defined pubsub channel sources
	m.createInternalChannels()

	// Get all user defined pubsub channel sources
	sourceManT, err := ctx.App("source")
	if err != nil {
		m.logger.Error("Unable to load the source manager", zap.Error(err))
	}
	sourceMan := sourceManT.(*source.App)
	sources := sourceMan.GetSources(m.Workspace, "pubsub")
	for _, src := range sources {
		channelSrc, ok := src.(Source)
		if ok {
			topic := channelSrc.GetChannel()
			m.channels = append(m.channels, topic)
		}
	}

	// Generate publish and subscribe API for each channel
	for path, channel := range m.Channels().Channels {
		// Get the publish and subscribe API of the channel
		publisherAPI := m.getProducerAPI(path, channel)
		subscriberAPI := m.getConsumerAPI(path, channel)

		m.apis = append(m.apis, publisherAPI, subscriberAPI)
	}

	// Generate AsyncAPI Doc and expose it via an API
	m.asyncapiDoc = m.generateASyncAPIDoc()
	m.apis = append(m.apis, m.exposeAsyncAPIDoc())
	return nil
}

// Start begins the pubsub app operations
func (m *Module) Start() error {
	return nil
}

// Stop ends the pubsub app operations
func (m *Module) Stop() error {
	return nil
}

// Cleanup cleans up the app
func (m *Module) Cleanup() error {
	_, err := connectorPool.Delete(getPoolKey())
	if err != nil {
		m.logger.Error("Unable to gracefully close  pubsub connector", zap.Error(err))
		return err
	}
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*Module)(nil)
	_ caddy.App         = (*Module)(nil)
)
