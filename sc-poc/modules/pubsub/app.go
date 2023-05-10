package pubsub

import (
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/pubsub/connectors"
	"github.com/spacecloud-io/space-cloud/sources/pubsub_channel"
)

var connectorPool = caddy.NewUsagePool()

func init() {
	caddy.RegisterModule(App{})
	apis.RegisterApp("pubsub", 100)
}

// App defines struct for pubsub app
type App struct {
	// For pubsub engine
	pubSub *gochannel.GoChannel

	// APIs
	apis apis.APIs

	// For internal usage
	logger      *zap.Logger
	asyncapiDoc *AsyncAPI
	channels    []Channel
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "pubsub",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the pubsub module.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)

	poolKey := getPoolKey()
	val, _, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
		return connectors.New()
	})
	if err != nil {
		a.logger.Error("Unable to open pubsub connector", zap.Error(err))
		return err
	}
	connector := val.(*connectors.Connector)
	a.pubSub = connector.GetPubsubClient()

	// Get all space-cloud defined pubsub channel sources
	a.createInternalChannels()

	// Get all user defined pubsub channel sources
	sourceManT, err := ctx.App("source")
	if err != nil {
		a.logger.Error("Unable to load the source manager", zap.Error(err))
	}
	sourceMan := sourceManT.(*source.App)
	sources := sourceMan.GetSources("pubsub")
	for _, src := range sources {
		channelSrc := src.(*pubsub_channel.PubsubChannelSource)
		topic := Channel{
			Name: channelSrc.Spec.Channel,
			Payload: ChannelPayload{
				Schema: channelSrc.Spec.Payload,
			},
		}

		a.channels = append(a.channels, topic)
	}

	// Generate publish and subscribe API for each channel
	for path, channel := range a.Channels().Channels {
		// Get the publish and subscribe API of the channel
		publisherAPI := a.getPublisherAPI(path, channel)
		subscriberAPI := a.getSubscriberAPI(path, channel.Name)

		a.apis = append(a.apis, publisherAPI, subscriberAPI)
	}

	// Generate AsyncAPI Doc and expose it via an API
	a.asyncapiDoc = a.generateASyncAPIDoc()
	a.apis = append(a.apis, a.exposeAsyncAPIDoc())
	return nil
}

// Start begins the pubsub app operations
func (a *App) Start() error {
	return nil
}

// Stop ends the pubsub app operations
func (a *App) Stop() error {
	return nil
}

// Cleanup cleans up the app
func (a *App) Cleanup() error {
	_, err := connectorPool.Delete(getPoolKey())
	if err != nil {
		a.logger.Error("Unable to gracefully close  pubsub connector", zap.Error(err))
		return err
	}
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
