package model

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
)
