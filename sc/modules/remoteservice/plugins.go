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

func (a *App) getRequestTemplatePlugin(endpoint *config.Endpoint) (string, string, string, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, config.PluginRequestTemplate)
	if err != nil {
		return "", "", "", err
	}

	template, ok := plugin.Params.(map[string]interface{})["template"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginRequestTemplate))
	}

	format, ok := plugin.Params.(map[string]interface{})["format"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s format value for endpoint", string(config.PluginRequestTemplate))
	}

	templateEngine, ok := plugin.Params.(map[string]interface{})["templateEngine"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginRequestTemplate))
	}

	return template.(string), format.(string), templateEngine.(string), err
}

func (a *App) getResponseTemplatePlugin(endpoint *config.Endpoint) (string, string, string, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, config.PluginResponseTemplate)
	if err != nil {
		return "", "", "", err
	}

	template, ok := plugin.Params.(map[string]interface{})["template"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginResponseTemplate))
	}

	format, ok := plugin.Params.(map[string]interface{})["format"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s format value for endpoint", string(config.PluginResponseTemplate))
	}

	templateEngine, ok := plugin.Params.(map[string]interface{})["templateEngine"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginRequestTemplate))
	}
	return template.(string), format.(string), templateEngine.(string), err
}

func (a *App) getGraphTemplatePlugin(endpoint *config.Endpoint) (string, string, string, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, config.PluginGraphTemplate)
	if err != nil {
		return "", "", "", err
	}

	template, ok := plugin.Params.(map[string]interface{})["template"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginGraphTemplate))
	}

	format, ok := plugin.Params.(map[string]interface{})["format"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s format value for endpoint", string(config.PluginGraphTemplate))
	}

	templateEngine, ok := plugin.Params.(map[string]interface{})["templateEngine"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginRequestTemplate))
	}
	return template.(string), format.(string), templateEngine.(string), err
}

func (a *App) getClaimsPlugin(endpoint *config.Endpoint) (string, string, string, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, config.PluginClaims)
	if err != nil {
		return "", "", "", err
	}

	template, ok := plugin.Params.(map[string]interface{})["template"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginClaims))
	}

	format, ok := plugin.Params.(map[string]interface{})["format"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s format value for endpoint", string(config.PluginClaims))
	}

	templateEngine, ok := plugin.Params.(map[string]interface{})["templateEngine"]
	if !ok {
		return "", "", "", fmt.Errorf("unable to get %s template value for endpoint", string(config.PluginRequestTemplate))
	}
	return template.(string), format.(string), templateEngine.(string), err
}

func (a *App) getTokenPlugin(endpoint *config.Endpoint) (string, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, config.PluginToken)
	if err != nil {
		return "", err
	}

	res, ok := plugin.Params.(string)
	if !ok {
		return "", fmt.Errorf("unable to get %s value for endpoint", string(config.PluginToken))
	}
	return res, err
}

func (a *App) getTimeoutPlugins(endpoint *config.Endpoint) (int, error) {
	plugin, err := a.getEndpointPluginsParams(endpoint.Plugins, config.PluginTimeout)
	if err != nil {
		return 0, err
	}

	res, ok := plugin.Params.(float64)
	if !ok {
		return 0, fmt.Errorf("unable to get %s value for endpoint", config.PluginTimeout)
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
