package common

import (
	"github.com/caddyserver/caddy/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/spacecloud-io/space-cloud/utils"
)

type (
	// ConfigType describes the configuration of space-cloud
	ConfigType map[string][]*unstructured.Unstructured
)

// PrepareConfig prepares a new caddy config based on the configuration provided
// TODO: Remove the previous configuration object
func PrepareConfig(configuration ConfigType) (*caddy.Config, error) {
	// First load the admin config
	c, err := utils.LoadAdminConfig()
	if err != nil {
		return nil, err
	}

	// Load all the managers
	c.AppsRaw = make(caddy.ModuleMap)
	c.AppsRaw["http"] = prepareHTTPHandlerApp(configuration)
	c.AppsRaw["source"] = prepareSourceManagerApp(configuration)

	// Load our providers
	c.AppsRaw["graphql"] = prepareEmptyApp()
	c.AppsRaw["rpc"] = prepareEmptyApp()
	c.AppsRaw["auth"] = prepareEmptyApp()
	c.AppsRaw["pubsub"] = prepareEmptyApp()

	return c, nil
}

func prepareEmptyApp() []byte {
	return []byte("{}")
}
