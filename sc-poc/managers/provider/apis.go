package provider

import (
	"github.com/spacecloud-io/space-cloud/managers/apis"
)

// GetAPIRoutes returns all the apis that are exposed by this app
func (a *App) GetAPIRoutes() apis.APIs {
	return a.apis
}

func (a *App) prepareAPIRoutes() {
	// Get the API Routes for each provider
	for _, workspace := range a.workspaces {
		for _, provider := range workspace.providers {
			apiGetter, ok := provider.(apis.App)
			if !ok {
				// Simply skip if that provider does not expose any apis
				continue
			}

			// Get the apis exposed by this provider
			apis := apiGetter.GetAPIRoutes()

			// Inject our workspace specific header for all non main workspaces
			for _, api := range apis {
				if api.Headers == nil {
					api.Headers = map[string][]string{}
				}

				// We want to inject the header only for non main workspaces
				if workspace.name != "main" {
					api.Headers["x-sc-workspace"] = []string{workspace.name}
				}

				// Add the workspace name to the api to gurantee uniqueness
				api.Name += workspace.name

				a.apis = append(a.apis, api)
			}

		}
	}
}
