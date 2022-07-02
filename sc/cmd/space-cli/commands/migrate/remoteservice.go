package migrate

import (
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func getRemoteServices(resource *model.SCConfig, configPath string) error {
	remoteServices, err := utils.ReadSpecObjectsFromFile(filepath.Join(configPath, "13-remote-services.yaml"))
	if err != nil {
		return err
	}

	for _, remoteService := range remoteServices {
		projectID := remoteService.Meta["project"]
		id := remoteService.Meta["id"]

		value := new(config.ServiceMigrate)
		if err := mapstructure.Decode(remoteService.Spec, value); err != nil {
			return err
		}

		newValue := new(config.Service)
		newValue.ID = id
		newValue.URL = value.URL
		newValue.Endpoints = make(map[string]*config.Endpoint)

		for endpointName, endpoint := range value.Endpoints {
			tempEndpoint := new(config.Endpoint)
			tempEndpoint.Kind = endpoint.Kind
			tempEndpoint.Method = endpoint.Method
			tempEndpoint.Path = endpoint.Path
			tempEndpoint.Plugins = []*config.EndpointPlugin{}

			reqTemplatePlugin := &config.EndpointPlugin{
				Type: config.PluginRequestTemplate,
				Params: map[string]interface{}{
					"template":       endpoint.ReqTmpl,
					"format":         endpoint.ReqPayloadFormat,
					"templateEngine": endpoint.Tmpl,
				},
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, reqTemplatePlugin)

			resTemplatePlugin := &config.EndpointPlugin{
				Type: config.PluginResponseTemplate,
				Params: map[string]interface{}{
					"template":       endpoint.ResTmpl,
					"format":         endpoint.ReqPayloadFormat,
					"templateEngine": endpoint.Tmpl,
				},
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, resTemplatePlugin)

			graphTemplatePlugin := &config.EndpointPlugin{
				Type: config.PluginGraphTemplate,
				Params: map[string]interface{}{
					"template":       endpoint.GraphTmpl,
					"format":         endpoint.ReqPayloadFormat,
					"templateEngine": endpoint.Tmpl,
				},
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, graphTemplatePlugin)

			claimsTemplatePlugin := &config.EndpointPlugin{
				Type: config.PluginClaims,
				Params: map[string]interface{}{
					"template":       endpoint.Claims,
					"format":         endpoint.ReqPayloadFormat,
					"templateEngine": endpoint.Tmpl,
				},
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, claimsTemplatePlugin)

			tokenPlugin := &config.EndpointPlugin{
				Type:   config.PluginToken,
				Params: endpoint.Token,
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, tokenPlugin)

			timeoutPlugin := &config.EndpointPlugin{
				Type:   config.PluginTimeout,
				Params: endpoint.Timeout,
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, timeoutPlugin)

			headerPlugin := &config.EndpointPlugin{
				Type:   config.PluginHeaders,
				Params: endpoint.Headers,
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, headerPlugin)

			cacheOptionsPlugin := &config.EndpointPlugin{
				Type:   config.PluginCacheOptions,
				Params: endpoint.CacheOptions,
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, cacheOptionsPlugin)

			rulePlugin := &config.EndpointPlugin{
				Type:   config.PluginRule,
				Params: endpoint.Rule,
			}
			tempEndpoint.Plugins = append(tempEndpoint.Plugins, rulePlugin)

			newValue.Endpoints[endpointName] = tempEndpoint
		}

		res := model.ResourceObject{
			Meta: model.ResourceMeta{
				Module: "remote-service",
				Type:   "config",
				Name:   id,
				Parents: map[string]string{
					"project": projectID,
					"id":      id,
				},
			},
			Spec: newValue,
		}

		module, ok := resource.Config["remote-service"]
		if !ok {
			module = make(model.ConfigModule)
		}
		resourceObjects, ok := module["config"]
		if !ok {
			resourceObjects = make([]*model.ResourceObject, 0)
		}

		resourceObjects = append(resourceObjects, &res)
		module["config"] = resourceObjects
		resource.Config["remote-service"] = module
	}
	return nil
}
