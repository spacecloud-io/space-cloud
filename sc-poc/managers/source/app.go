package source

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

// App describes the source manager app
type App struct {
	Config map[string][]json.RawMessage `json:"config"`

	// Internal stuff
	logger    *zap.Logger
	sourceMap map[string]map[string]Sources // [workspace] -> [provider] -> source
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

	// Initialise internal datastructures
	a.sourceMap = make(map[string]map[string]Sources, len(a.Config))

	// Create a map of sources
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

			// Extract the workspace this source belongs to
			workspace := GetWorkspaceNameFromSource(source)

			// Add the workspace to the source map if it doesn't already exist
			if _, p := a.sourceMap[workspace]; !p {
				a.sourceMap[workspace] = make(map[string]Sources, 1)
			}

			// Add the provider for all supported providers
			for _, provider := range source.GetProviders() {
				a.sourceMap[workspace][provider] = append(a.sourceMap[workspace][provider], source)
			}

			// Register the source against all requested providers
			if plugin, ok := source.(Plugin); ok {
				a.plugins = append(a.plugins, plugin.GetPluginDetails())
			}
		}
	}

	// Delete the `main` and `root` workspaces since they are internal

	// Sort the sources for each provider in each workspace
	for _, workspace := range a.sourceMap {
		for _, sources := range workspace {
			sources.Sort()
		}
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
