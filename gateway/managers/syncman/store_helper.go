package syncman

import (
	"reflect"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func (s *PostgresStore) helperAddOrUpdate(resourceID string, resourceType config.Resource, resource interface{}, project string) (string, string, config.Resource, interface{}) {
	var evenType string
	var obj interface{}

	switch resourceType {
	case config.ResourceProject:
		projectConfig := s.globalConfig.Projects[project].ProjectConfig
		if !reflect.DeepEqual(projectConfig, resource) {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceDatabaseConfig:
		databaseConfig, ok := s.globalConfig.Projects[project].DatabaseConfigs[resourceID]
		if ok {
			if !reflect.DeepEqual(databaseConfig, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceDatabaseRule:
		databaseRule, ok := s.globalConfig.Projects[project].DatabaseRules[resourceID]
		if ok {
			if !reflect.DeepEqual(databaseRule, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)

		}
	case config.ResourceDatabaseSchema:
		databaseSchema, ok := s.globalConfig.Projects[project].DatabaseSchemas[resourceID]
		if ok {
			if !reflect.DeepEqual(databaseSchema, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceDatabasePreparedQuery:
		databasePreparedQuery, ok := s.globalConfig.Projects[project].DatabasePreparedQueries[resourceID]
		if ok {
			if !reflect.DeepEqual(databasePreparedQuery, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceFileStoreConfig:
		fileStoreConfig := s.globalConfig.Projects[project].FileStoreConfig
		if !reflect.DeepEqual(fileStoreConfig, resource) {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceFileStoreRule:
		fileStoreRule, ok := s.globalConfig.Projects[project].FileStoreRules[resourceID]
		if ok {
			if !reflect.DeepEqual(fileStoreRule, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceEventingConfig:
		eventingConfig := s.globalConfig.Projects[project].EventingConfig
		if !reflect.DeepEqual(eventingConfig, resource) {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceEventingTrigger:
		eventingTrigger, ok := s.globalConfig.Projects[project].EventingTriggers[resourceID]
		if ok {
			if !reflect.DeepEqual(eventingTrigger, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceEventingRule:
		eventingRule, ok := s.globalConfig.Projects[project].EventingRules[resourceID]
		if ok {
			if !reflect.DeepEqual(eventingRule, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceEventingSchema:
		eventingSchemas, ok := s.globalConfig.Projects[project].EventingSchemas[resourceID]
		if ok {
			if !reflect.DeepEqual(eventingSchemas, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceRemoteService:
		remoteService, ok := s.globalConfig.Projects[project].RemoteService[resourceID]
		if ok {
			if !reflect.DeepEqual(remoteService, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceIngressGlobal:
		global := s.globalConfig.Projects[project].IngressGlobal
		if !reflect.DeepEqual(global, resource) {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceIngressRoute:
		routes, ok := s.globalConfig.Projects[project].IngressRoutes[resourceID]
		if ok {
			if !reflect.DeepEqual(routes, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceAuthProvider:
		auth, ok := s.globalConfig.Projects[project].Auths[resourceID]
		if ok {
			if !reflect.DeepEqual(auth, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceProjectLetsEncrypt:
		letsEncrypt := s.globalConfig.Projects[project].LetsEncrypt
		if !reflect.DeepEqual(letsEncrypt, resource) {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)

		}
	case config.ResourceCluster:
		cluster := s.globalConfig.ClusterConfig
		if !reflect.DeepEqual(cluster, resource) {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceIntegration:
		integrationconfig, ok := s.globalConfig.Integrations[resourceID]
		if ok {
			if !reflect.DeepEqual(integrationconfig, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	case config.ResourceIntegrationHook:
		integrationHook, ok := s.globalConfig.IntegrationHooks[resourceID]
		if ok {
			if !reflect.DeepEqual(integrationHook, resource) {
				evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
			}
		} else {
			evenType, resourceID, resourceType, resource = onAddOrUpdateResource(config.ResourceAddEvent, obj)
		}
	}

	return evenType, resourceID, resourceType, resource
}
