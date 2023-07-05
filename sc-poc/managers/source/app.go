package source

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"go.uber.org/zap"
)

// App describes the source manager app
type App struct {
	Config map[string][]json.RawMessage `json:"config"`

	// Internal stuff
	logger    *zap.Logger
	sourceMap map[string]Sources
	plugins   []v1alpha1.HTTPPlugin
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "source",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the source manager.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)

	// Create a map of sources
	a.sourceMap = make(map[string]Sources, len(a.Config))
	a.plugins = []v1alpha1.HTTPPlugin{
		{
			Name:   "",
			Driver: "deny_user",
		},
		{
			Name:   "",
			Driver: "authenticate-user",
		},
	}
	for key, list := range a.Config {
		gvr := GetResourceGVR(key)

		// Make one module for each source
		for _, c := range list {
			// LoadModuleByID will automatically call provision and validate for us. We can safely assume that the source
			// module is ready to be used if no error is returned
			t, err := ctx.LoadModuleByID(key, c)
			if err != nil {
				a.logger.Warn("Unable to load module for source", zap.String("group", gvr.Group), zap.String("version", gvr.Version), zap.String("resource", gvr.Resource))
				continue
			}

			source, ok := t.(Source)
			if !ok {
				a.logger.Error("Loaded source is not of a valid type", zap.String("group", gvr.Group), zap.String("version", gvr.Version), zap.String("resource", gvr.Resource))
				continue
			}

			if plugin, ok := source.(Plugin); ok {
				a.plugins = append(a.plugins, plugin.GetPluginDetails())
			}

			// Add the provider for all supported providers
			for _, provider := range source.GetProviders() {
				a.sourceMap[provider] = append(a.sourceMap[provider], source)
			}
		}
	}

	// Sort the sources for each provider
	for _, s := range a.sourceMap {
		s.Sort()
	}

	return nil
}

// Start begins the source manager operations
func (a *App) Start() error {
	return nil
}

// Stop ends the source manager operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
