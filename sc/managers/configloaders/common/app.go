package common

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

// PrepareConfig prepares a new caddy config based on a SC config object
func PrepareConfig(fileConfig *model.SCConfig) (*caddy.Config, error) {
	// First load the admin config
	c, err := utils.LoadAdminConfig(false)
	if err != nil {
		return nil, err
	}
	c.AppsRaw = make(caddy.ModuleMap)

	// Load all the apps. Each app will have data for all the projects combined
	c.AppsRaw["http"] = prepareHTTPHanndlerApp()
	c.AppsRaw["configman"] = prepareConfigManApp()
	c.AppsRaw["config_store"] = prepareStoreApp()
	c.AppsRaw["database"] = prepareDatabaseApp(fileConfig)
	c.AppsRaw["graphql"] = prepareGraphQLApp()
	c.AppsRaw["admin"] = prepareAdminApp()

	return c, nil
}
