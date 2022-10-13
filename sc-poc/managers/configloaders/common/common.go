package common

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/utils"
)

// PrepareConfig prepares a new caddy config based on the configuration provided
func PrepareConfig() (*caddy.Config, error) {
	// First load the admin config
	c, err := utils.LoadAdminConfig(false)
	if err != nil {
		return nil, err
	}

	// Load all the apps.
	c.AppsRaw = make(caddy.ModuleMap)
	c.AppsRaw["graphql"] = prepareGraphQLApp()
	c.AppsRaw["http"] = prepareHTTPHanndlerApp()

	return c, nil
}
