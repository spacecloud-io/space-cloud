package syncman

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/getlantern/deepcopy"
	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func (s *Manager) delete(projectID string) {
	delete(s.projectConfig.Projects, projectID)
}

type scServices []*service

func (a scServices) Len() int           { return len(a) }
func (a scServices) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a scServices) Less(i, j int) bool { return a[i].id < a[j].id }

func calcTokens(n int, tokens int, i int) (start int, end int) {
	tokensPerMember := int(math.Ceil(float64(tokens) / float64(n)))
	start = tokensPerMember * i
	end = start + tokensPerMember - 1
	if end > tokens {
		end = tokens - 1
	}
	return
}

func calcIndex(token, totalTokens, n int) int {
	bucketSize := totalTokens / n
	return token / bucketSize
}

func (s *Manager) getGatewayIndex() int {
	index := 0

	for i, v := range s.services {
		if v.id == s.nodeID {
			index = i
			break
		}
	}
	return index
}

// getConfigWithoutLock returns the config present in the state
func (s *Manager) getConfigWithoutLock(ctx context.Context, projectID string) (*config.Project, error) {
	project, ok := s.projectConfig.Projects[projectID]
	if ok {
		p := new(config.Project)
		_ = deepcopy.Copy(p, project)
		return p, nil
	}

	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unknown project (%s) provided", projectID), nil, nil)
}

func (s *Manager) checkIfDbAliasExists(dbConfigs config.DatabaseConfigs, dbAlias string) (*config.DatabaseConfig, bool) {
	for _, databaseConfig := range dbConfigs {
		if dbAlias == databaseConfig.DbAlias {
			return databaseConfig, true
		}
	}
	return nil, false
}

func splitResourceID(ctx context.Context, resourceID string) (clusterID string, projectID string, resource config.Resource, err error) {
	arr := strings.Split(resourceID, "--")
	// ResourceId format --> clusterId--ProjectId--resourceType--someId-...
	if len(arr) < 4 {
		return "", "", "", helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid resource id (%s) provided", resourceID), nil, nil)
	}
	return arr[0], arr[1], config.Resource(arr[2]), nil
}

func (s *Manager) validateResource(ctx context.Context, eventType string, resourceID string, resourceType config.Resource, resource interface{}) (bool, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	globalConfig := s.projectConfig

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
	case config.ResourceSecurityFunction:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.SecurityFunction)
			if err := mapstructure.Decode(resource, value); err != nil {
				return false, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.SecurityFunction{}", reflect.TypeOf(resource)), nil, nil)
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
	}

	if resourceType == config.ResourceProject {
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.ProjectConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to find project", err, map[string]interface{}{"project": projectID})
		return nil
	}

	switch resourceType {
	case config.ResourceAuthProvider:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.AuthStub)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
	case config.ResourceSecurityFunction:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.SecurityFunction)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.SecurityFunction{}", reflect.TypeOf(resource)), nil, nil)
			}

			if project.SecurityFunctions == nil {
				project.SecurityFunctions = config.SecurityFunctions{resourceID: value}
			} else {
				project.SecurityFunctions[resourceID] = value
			}
		case config.ResourceDeleteEvent:
			delete(project.SecurityFunctions, resourceID)
		}

		return nil

	case config.ResourceDatabaseConfig:
		switch eventType {
		case config.ResourceAddEvent, config.ResourceUpdateEvent:
			value := new(config.DatabaseConfig)
			if err := mapstructure.Decode(resource, value); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type provided for resource (%s) expecting (%v) got (%v)", resourceType, "config.Auth{}", reflect.TypeOf(resource)), nil, nil)
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
