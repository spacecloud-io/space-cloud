package provider

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
)

func (a *App) generateProviders(ctx caddy.Context, workspace string) (map[string]any, error) {
	// Make map of apps with capacity equals to number of registered providers
	m := make(map[string]any, len(registeredProviders))

	// Create body
	providerBody, _ := json.Marshal(map[string]string{"workspace": workspace})
	for _, p := range registeredProviders {
		provider, err := ctx.LoadModuleByID(fmt.Sprintf("provider.%s", p.name), providerBody)
		if err != nil {
			a.logger.Error("Unable to create module", zap.String("provider", p.name), zap.String("workspace", workspace), zap.Error(err))
			return nil, fmt.Errorf("unable to create provider '%s' in worskspace '%s'", p.name, workspace)
		}
		m[p.name] = provider
	}

	return m, nil
}
