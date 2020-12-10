package model

import "time"

const (
	// AllContainers gets all container running in a cluster
	AllContainers = "all"
	// ServiceContainers gets all the service containers running in a cluster
	ServiceContainers = "service"
	// DbContainers gets all the database containers running in a cluster
	DbContainers = "db"
	// RegistryContainers gets all the registry containers running in a cluster
	RegistryContainers = "registry"
	// ScContainers gets runner & gateway container running in a cluster
	ScContainers = "space-cloud"

	// ApplyWithNoDelay used as param to indicate to used 0 delay while applying spec objects
	ApplyWithNoDelay = time.Duration(0)

	// HelmSpaceCloudChartDownloadURL url of space cloud helm chart
	HelmSpaceCloudChartDownloadURL = "https://storage.googleapis.com/space-cloud/helm/space-cloud.tgz"
	// HelmPostgresChartDownloadURL url of postgres helm chart
	HelmPostgresChartDownloadURL = "https://storage.googleapis.com/space-cloud/helm/postgres.tgz"
	// HelmMysqlChartDownloadURL url of mysql helm chart
	HelmMysqlChartDownloadURL = "https://storage.googleapis.com/space-cloud/helm/mysql.tgz"
	// HelmSQLServerCloudChartDownloadURL url of sql server helm chart
	HelmSQLServerCloudChartDownloadURL = "https://storage.googleapis.com/space-cloud/helm/sqlserver.tgz"
	// HelmMongoChartDownloadURL url of mongo helm chart
	HelmMongoChartDownloadURL = "https://storage.googleapis.com/space-cloud/helm/mongo.tgz"
	// HelmSpaceCloudNamespace space cloud namespace for helm
	HelmSpaceCloudNamespace = "space-cloud"
)

// Version version of space cli
const Version string = "0.20.1"
