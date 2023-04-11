package pubsub

import (
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/modules/pubsub/connectors"
	"go.uber.org/zap"
)

var connectorPool = caddy.NewUsagePool()

func init() {
	caddy.RegisterModule(App{})
}

// App defines struct for pubsub app
type App struct {
	// For pubsub engine
	pubSub *gochannel.GoChannel

	// For internal usage
	logger *zap.Logger
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
	a.pubSub = connector.PubSubClient

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
	return nil
}
