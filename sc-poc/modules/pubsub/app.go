package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/caddyserver/caddy/v2"
)

// App defines struct for pubsub app
type App struct {
	// For pubsub engine
	pubSub *gochannel.GoChannel
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
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)
	a.pubSub = pubSub
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
