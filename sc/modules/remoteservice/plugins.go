package remoteservice

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
)

func (a *App) getEndpointPluginsParams(plugins []*config.EndpointPlugin, name config.EndpointPluginType) (*config.EndpointPlugin, error) {
	for _, plugin := range plugins {
		if plugin.Type == name {
			return plugin, nil
		}
	}
	return nil, fmt.Errorf("plugin (%s) not present in endpoint configuration", name)
}

func (a *App) getStringOutputFromPlugins(endpoint *config.Endpoint, pluginType config.EndpointPluginType) (string, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, pluginType)
	if err != nil {
		return "", err
	}

	res, ok := plugin.Params.(string)
	if !ok {
		return "", fmt.Errorf("unable to get %s value for endpoint", string(pluginType))
	}
	return res, err
}

func (a *App) getIntOutputFromPlugins(endpoint *config.Endpoint, pluginType config.EndpointPluginType) (int, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, pluginType)
	if err != nil {
		return 0, err
	}

	res, ok := plugin.Params.(float64)
	if !ok {
		return 0, fmt.Errorf("unable to get %s value for endpoint", pluginType)
	}
	return int(res), nil
}

func (a *App) getHeadersFromPlugins(endpoint *config.Endpoint) (*config.Headers, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, config.PluginHeaders)
	if err != nil {
		return nil, err
	}

	headers := new(config.Headers)
	if err := mapstructure.Decode(plugin.Params, headers); err != nil {
		return nil, err
	}
	return headers, err
}
