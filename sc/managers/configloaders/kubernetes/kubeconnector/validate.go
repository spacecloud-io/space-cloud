package kubeconnector

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spaceuptech/helpers"
)

// validateResource check type of incomming config
func (c *Connector) validateResource(ctx context.Context, eventType string, resourceID string, resourceType config.Resource, resource interface{}) (bool, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	globalConfig := c.ProjectsConfig

	_, projectID, rt, err := splitResourceID(ctx, resourceID)
	if err != nil {
		return false, err
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
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.ClusterConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(globalConfig.ClusterConfig, value) {
				return true, nil
			}
		}
		return false, nil

	case config.ResourceIntegration:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.IntegrationConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.IntegrationConfig{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(globalConfig.Integrations[resourceID], value) {
				return true, nil
			}
		}
		return false, nil

	case config.ResourceIntegrationHook:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.IntegrationHook)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.IntegrationHook{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(globalConfig.IntegrationHooks[resourceID], value) {
				return true, nil
			}
		}
		return false, nil

	}

	if resourceType == config.ResourceProject {
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.ProjectConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(globalConfig.Projects[projectID], value) {
				return true, nil
			}
		}
		return false, nil
	}

	project, ok := globalConfig.Projects[projectID]
	if !ok {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to find project", err, map[string]interface{}{"project": projectID})
		return false, nil
	}

	switch resourceType {
	case config.ResourceAuthProvider:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.AuthStub)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.Auths[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceDatabaseConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatabaseConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.DatabaseConfigs[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceDatabaseSchema:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatabaseSchema)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.DatabaseSchemas[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceDatabaseRule:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatabaseRule)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.DatabaseRules[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceDatabasePreparedQuery:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatbasePreparedQuery)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.DatabasePreparedQueries[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceEventingConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.EventingConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.EventingConfig, value) {
				return true, nil
			}

		}
		return false, nil
	case config.ResourceEventingSchema:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.EventingSchema)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.EventingSchemas[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceEventingRule:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.Rule)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.EventingRules[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceEventingTrigger:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.EventingTrigger)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.EventingTriggers[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceFileStoreConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.FileStoreConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.FileStoreConfig, value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceFileStoreRule:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.FileRule)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.FileStoreRules[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceProjectLetsEncrypt:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.LetsEncrypt)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.LetsEncrypt, value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceIngressRoute:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.Route)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.IngressRoutes[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceIngressGlobal:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.GlobalRoutesConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.IngressGlobal, value) {
				return true, nil
			}
		}
		return false, nil
	case config.ResourceRemoteService:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.Service)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
			}

			if reflect.DeepEqual(project.RemoteService[resourceID], value) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unknown resource type (%s) provided", resourceType)
	}
}
