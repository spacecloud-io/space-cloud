package syncman

// Service is used to watch changes in services
type Service interface {
	WatchServices(cb func(projects scServices)) error
}
