package common

import (
	"github.com/caddyserver/caddy/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/spacecloud-io/space-cloud/utils"
)

// PrepareConfig prepares a new caddy config based on the configuration provided
// TODO: Remove the previous configuration object
func PrepareConfig(configuration, newConfig map[string][]*unstructured.Unstructured) (*caddy.Config, error) {
	// First load the admin config
	c, err := utils.LoadAdminConfig()
	if err != nil {
		return nil, err
	}

	// Load all the managers
	c.AppsRaw = make(caddy.ModuleMap)
	c.AppsRaw["http"] = prepareHTTPHandlerApp()
	c.AppsRaw["source"] = prepareSourceManagerApp(newConfig)

	// Load our providers
	c.AppsRaw["graphql"] = prepareEmptyApp()
	c.AppsRaw["rpc"] = prepareEmptyApp()
	c.AppsRaw["auth"] = prepareEmptyApp()

	return c, nil
}

func prepareEmptyApp() []byte {
	return []byte("{}")
}
