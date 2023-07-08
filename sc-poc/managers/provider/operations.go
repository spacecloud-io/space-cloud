package provider

import (
	"fmt"
	"sync"
)

// The necesary global objects to hold all registered providers
var (
	providersLock sync.RWMutex

	registeredProviders providers
)

// Register adds a caddy module as a provider. The providers with the highest
// priority values gets provisioned first
func Register(name string, priority int) {
	providersLock.Lock()
	defer providersLock.Unlock()

	registeredProviders = append(registeredProviders, provider{name, priority})
	registeredProviders.sort()
}

func (a *App) GetProvider(workspace, provider string) (any, error) {
	// Get the provider from the workspace
	for _, ws := range a.workspaces {
		// Skip if workspace name doesn't match
		if ws.name != workspace {
			continue
		}

		p, ok := ws.providers[provider]
		if !ok {
			return nil, fmt.Errorf("unable to find provider '%s' in workspace '%s'", provider, workspace)
		}

		return p, nil
	}

	return nil, fmt.Errorf("unable to find workspace '%s' to get provider '%s'", workspace, provider)
}
