package common

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/utils"
)

// PrepareConfig prepares a new caddy config based on a SC config object
func PrepareConfig(scConfig *config.Config) *caddy.Config {
	// First load the admin config
	c := utils.LoadAdminConfig(false)
	c.AppsRaw = make(caddy.ModuleMap, 0)

	// Load all the apps. Each app will have data for all the projects combined
	c.AppsRaw["database"] = prepareDatabaseApp(scConfig)

	return c
}
