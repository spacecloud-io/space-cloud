package config

import (
	"fmt"
	"strings"
)

// ResourceFetchingOrder order is according to space cli apply functionality
var ResourceFetchingOrder = []Resource{
	ResourceProject,
	ResourceDatabaseConfig,
	ResourceDatabaseRule,
	ResourceDatabaseSchema,
	ResourceDatabasePreparedQuery,
	ResourceFileStoreConfig,
	ResourceFileStoreRule,
	ResourceEventingConfig,
	ResourceEventingTrigger,
	ResourceEventingRule,
	ResourceEventingSchema,
	ResourceRemoteService,
	ResourceIngressGlobal,
	ResourceIngressRoute,
	ResourceAuthProvider,
	ResourceSecurityFunction,
	ResourceProjectLetsEncrypt,
	ResourceCluster,
	ResourceIntegration,
	ResourceIntegrationHook,
}

// Resource is a resource type
type Resource string

const (
	// ResourceSecurityFunction is a resource
	ResourceSecurityFunction Resource = "security-function"

	// ResourceAuthProvider is a resource
	ResourceAuthProvider Resource = "auth-provider"

	// ResourceProject is a resource
	ResourceProject Resource = "project"

	// ResourceDatabaseConfig is a resource
	ResourceDatabaseConfig Resource = "db-config"
	// ResourceDatabaseSchema is a resource
	ResourceDatabaseSchema Resource = "db-schema"
	// ResourceDatabaseRule is a resource
	ResourceDatabaseRule Resource = "db-rule"
	// ResourceDatabasePreparedQuery is a resource
	ResourceDatabasePreparedQuery Resource = "db-prepared-query"

	// ResourceEventingConfig is a resource
	ResourceEventingConfig Resource = "eventing-config"
	// ResourceEventingSchema is a resource
	ResourceEventingSchema Resource = "eventing-schema"
	// ResourceEventingTrigger is a resource
	ResourceEventingTrigger Resource = "eventing-trigger"
	// ResourceEventingRule is a resource
	ResourceEventingRule Resource = "eventing-rule"

	// ResourceFileStoreConfig is a resource
	ResourceFileStoreConfig Resource = "filestore-config"
	// ResourceFileStoreRule is a resource
	ResourceFileStoreRule Resource = "filestore-rule"

	// ResourceProjectLetsEncrypt is a resource
	ResourceProjectLetsEncrypt Resource = "letsencrypt"

	// ResourceIngressRoute is a resource
	ResourceIngressRoute Resource = "ingress-route"
	// ResourceIngressGlobal is a resource
	ResourceIngressGlobal Resource = "ingress-global"

	// ResourceRemoteService is a resource
	ResourceRemoteService Resource = "remote-service"

	// ResourceIntegration is a resource
	ResourceIntegration Resource = "integration"
	// ResourceIntegrationHook is a resource
	ResourceIntegrationHook Resource = "integration-hook"
	// ResourceCluster is a resource
	ResourceCluster Resource = "cluster"

	// ResourceDeployService is a resource
	// ResourceDeployService Resource = "service"
	// ResourceDeployServiceRoute is a resource
	// ResourceDeployServiceRoute Resource = "service-route"
	// ResourceSecret is a resource
	// ResourceSecret Resource = "secret"

	// API resource
	// ResourceDatabaseCreate      Resource = "db-create"
	// ResourceDatabaseRead        Resource = "db-read"
	// ResourceDatabaseUpdate      Resource = "db-update"
	// ResourceDatabaseDelete      Resource = "db-delete"
	// ResourceDatabaseAggregate   Resource = "db-aggregate"
	// ResourceDatabasePreparedSQL Resource = "db-prepared-sql"
	// ResourceEventTrigger        Resource = "eventing-trigger"
	// ResourceServiceCall         Resource = "service-call"
	// ResourceInternalAPIAccess    Resource = "internal-api-access"

	// Misc resource
	// ResourceLoadEnv          Resource = "load-env"
	// ResourceConfigPermission Resource = "config-permission"
	// ResourceAdminLogin       Resource = "admin-login"
	// ResourceCredentials      Resource = "credentials"
	// ResourceInternalToken    Resource = "internal-token"
	// ResourceAdminToken       Resource = "admin-token"

	// ResourceDeleteEvent syncman uses this const to represent a delete event
	ResourceDeleteEvent = "delete"
	// ResourceUpdateEvent syncman uses this const to represent a delete event
	ResourceUpdateEvent = "update"
	// ResourceAddEvent syncman uses this const to represent a delete event
	ResourceAddEvent = "add"
)

// GenerateResourceID This a generic function for all resource to generate a unique id in a space cloud cluster
// params are in the form (clusterId,projectId,idOfCurrentResource) for most of the resources
// database resources are the exception they take the form (clusterId,projectId,dbAlias,idOfCurrentResource)
func GenerateResourceID(clusterID, projectID string, resourceType Resource, arr ...string) string {
	return fmt.Sprintf("%s--%s--%s--%s", clusterID, projectID, resourceType, strings.Join(arr, "-"))
}
