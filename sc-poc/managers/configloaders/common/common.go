package common

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// PrepareConfig prepares a new caddy config based on the configuration provided
func PrepareConfig(configuration map[string][]*unstructured.Unstructured) (*caddy.Config, error) {
	// First load the admin config
	c, err := utils.LoadAdminConfig(false)
	if err != nil {
		return nil, err
	}

	// Load all the apps.
	c.AppsRaw = make(caddy.ModuleMap)
	c.AppsRaw["auth"] = prepareAuthApp(configuration)
	c.AppsRaw["graphql"] = prepareGraphQLApp(configuration)
	c.AppsRaw["rest"] = prepareRestApp(configuration)
	c.AppsRaw["http"] = prepareHTTPHanndlerApp()

	return c, nil
}
