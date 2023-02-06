package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spaceuptech/helpers"
	v1 "k8s.io/api/core/v1"
)

func splitResourceID(ctx context.Context, resourceID string) (clusterID string, projectID string, resource config.Resource, err error) {
	arr := strings.Split(resourceID, "--")
	// ResourceId format --> clusterId--ProjectId--resourceType--someId-...
	if len(arr) < 4 {
		return "", "", "", helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid resource id (%s) provided", resourceID), nil, nil)
	}
	return arr[0], arr[1], config.Resource(arr[2]), nil
}

// onAddOrUpdateResource reads data for config maps
func onAddOrUpdateResource(eventType string, obj interface{}) (string, string, config.Resource, interface{}) {
	configMap := obj.(*v1.ConfigMap)
	resourceID, ok := configMap.Data["id"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (id) field was not found in config map data", eventType), nil, nil)
		return "", "", "", nil
	}

	resourceType, ok := configMap.Labels["kind"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (kind) label was not found in config map", eventType), nil, nil)
		return "", "", "", nil
	}

	dataJSONString, ok := configMap.Data["data"]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("%s event occured on resource config map, but (resource) field was not found in config map data", eventType), nil, nil)
		return "", "", "", nil
	}

	v := make(map[string]interface{})
	if err := json.Unmarshal([]byte(dataJSONString), &v); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal resource config map data while watching kube store project", nil, map[string]interface{}{"resourceId": resourceID, "eventType": eventType})
		return "", "", "", nil
	}
	return eventType, resourceID, config.Resource(resourceType), v
}

// NOTE: any change made in this function should also be reflected into validateResource() method of sync man
func updateResource(ctx context.Context, eventType string, globalConfig *config.Config, resourceID string, resourceType config.Resource, resource interface{}) error {
	if globalConfig == nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot provide empty value for config", nil, nil)
	}

	_, projectID, rt, err := splitResourceID(ctx, resourceID)
	if err != nil {
		return err
	}

	// if resource type not provided extract in from resource id
	if resourceType == "" {
		resourceType = rt
	}

	// check cluster level resources first
	switch resourceType {
	case config.ResourceCluster:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.ClusterConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.ClusterConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			globalConfig.ClusterConfig = value
		case config.ResourceDeleteEvent:
		}
		return nil

	case config.ResourceIntegration:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.IntegrationConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.IntegrationConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			if globalConfig.Integrations == nil {
				globalConfig.Integrations = config.Integrations{resourceID: value}
			} else {
				globalConfig.Integrations[resourceID] = value
			}

		case config.ResourceDeleteEvent:
			delete(globalConfig.Integrations, resourceID)
		}
		return nil

	case config.ResourceIntegrationHook:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.IntegrationHook)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.IntegrationHook{}", reflect.TypeOf(resource)), nil, nil)
			}

			if globalConfig.Integrations == nil {
				globalConfig.IntegrationHooks = config.IntegrationHooks{resourceID: value}
			} else {
				globalConfig.IntegrationHooks[resourceID] = value
			}

		case config.ResourceDeleteEvent:
			delete(globalConfig.IntegrationHooks, resourceID)
		}
		return nil

	case config.ResourceCacheConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.CacheConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			globalConfig.CacheConfig = value
		case config.ResourceDeleteEvent:
			globalConfig.CacheConfig.Enabled = false
		}

		return nil
	}

	// check project level resources
	if resourceType == config.ResourceProject {
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.ProjectConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.ProjectConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			projectConfig, ok := globalConfig.Projects[projectID]
			if !ok {
				globalConfig.Projects[projectID] = config.GenerateEmptyProject(value)
			} else {
				projectConfig.ProjectConfig = value
			}
		case config.ResourceDeleteEvent:
			delete(globalConfig.Projects, projectID)
		}
		return nil
	}

	project, ok := globalConfig.Projects[projectID]
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to find project", err, map[string]interface{}{"project": projectID})
	}

	switch resourceType {
	case config.ResourceAuthProvider:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.AuthStub)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.AuthStub{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.Auths == nil {
				project.Auths = config.Auths{resourceID: value}
			} else {
				project.Auths[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.Auths, resourceID)
		}

		return nil

	case config.ResourceDatabaseConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatabaseConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.DatabaseConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.DatabaseConfigs == nil {
				project.DatabaseConfigs = config.DatabaseConfigs{resourceID: value}
			} else {
				project.DatabaseConfigs[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.DatabaseConfigs, resourceID)
		}

		return nil
	case config.ResourceDatabaseSchema:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatabaseSchema)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.DatabaseSchema{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.DatabaseSchemas == nil {
				project.DatabaseSchemas = config.DatabaseSchemas{resourceID: value}
			} else {
				project.DatabaseSchemas[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.DatabaseSchemas, resourceID)
		}

		return nil

	case config.ResourceDatabaseRule:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatabaseRule)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.DatabaseRule{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.DatabaseRules == nil {
				project.DatabaseRules = config.DatabaseRules{resourceID: value}
			} else {
				project.DatabaseRules[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.DatabaseRules, resourceID)
		}

		return nil

	case config.ResourceDatabasePreparedQuery:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatbasePreparedQuery)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.DatbasePreparedQuery{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.DatabasePreparedQueries == nil {
				project.DatabasePreparedQueries = config.DatabasePreparedQueries{resourceID: value}
			} else {
				project.DatabasePreparedQueries[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.DatabasePreparedQueries, resourceID)
		}

		return nil

	case config.ResourceEventingConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.EventingConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.EventingConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			project.EventingConfig = value

		case config.ResourceDeleteEvent:
			project.EventingConfig.Enabled = false
		}

		return nil

	case config.ResourceEventingSchema:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.EventingSchema)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.EventingSchema{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.EventingSchemas == nil {
				project.EventingSchemas = config.EventingSchemas{resourceID: value}
			} else {
				project.EventingSchemas[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.EventingSchemas, resourceID)
		}

		return nil

	case config.ResourceEventingRule:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.Rule)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Rule{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.EventingRules == nil {
				project.EventingRules = config.EventingRules{resourceID: value}
			} else {
				project.EventingRules[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.EventingRules, resourceID)
		}

		return nil

	case config.ResourceEventingTrigger:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.EventingTrigger)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.EventingTrigger{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.EventingTriggers == nil {
				project.EventingTriggers = config.EventingTriggers{resourceID: value}
			} else {
				project.EventingTriggers[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.EventingTriggers, resourceID)
		}

		return nil

	case config.ResourceFileStoreConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.FileStoreConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.FileStoreConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			project.FileStoreConfig = value
		case config.ResourceDeleteEvent:
		}

		return nil

	case config.ResourceFileStoreRule:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.FileRule)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.FileRule{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.FileStoreRules == nil {
				project.FileStoreRules = config.FileStoreRules{resourceID: value}
			} else {
				project.FileStoreRules[resourceID] = value
			}
		case config.ResourceDeleteEvent:
		}

		return nil

	case config.ResourceProjectLetsEncrypt:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.LetsEncrypt)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.LetsEncrypt{}", reflect.TypeOf(resource)), nil, nil)
			}

			project.LetsEncrypt = value
		case config.ResourceDeleteEvent:
		}

		return nil

	case config.ResourceIngressRoute:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.Route)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Route{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.IngressRoutes == nil {
				project.IngressRoutes = config.IngressRoutes{resourceID: value}
			} else {
				project.IngressRoutes[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.IngressRoutes, resourceID)
		}

		return nil

	case config.ResourceIngressGlobal:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.GlobalRoutesConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.GlobalRoutesConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			project.IngressGlobal = value
		case config.ResourceDeleteEvent:
		}

		return nil

	case config.ResourceRemoteService:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.Service)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Service{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.RemoteService == nil {
				project.RemoteService = config.Services{resourceID: value}
			} else {
				project.RemoteService[resourceID] = value
			}

		case config.ResourceDeleteEvent:
			delete(project.RemoteService, resourceID)
		}

		return nil

	default:
		return fmt.Errorf("unknown resource type (%s) provided", resourceType)
	}
}
