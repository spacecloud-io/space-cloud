package common

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/utils"
)

// PrepareConfig prepares a new caddy config based on a SC config object
func PrepareConfig(scConfig *config.Config) (*caddy.Config, error) {
	// First load the admin config
	c, err := utils.LoadAdminConfig(false)
	if err != nil {
		return nil, err
	}
	c.AppsRaw = make(caddy.ModuleMap)

	// Load all the apps. Each app will have data for all the projects combined
	c.AppsRaw["database"] = prepareDatabaseApp(scConfig)
	c.AppsRaw["graphql"] = prepareGraphQLApp()
	c.AppsRaw["http"] = prepareHTTPHanndlerApp()
	c.AppsRaw["configman"] = prepareStoreApp()

	return c, nil
}
