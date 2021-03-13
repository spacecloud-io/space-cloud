package syncman

// LocalService is an object for storing localservice information
type LocalService struct {
	services scServices
}

// NewLocalService creates a new local service
func NewLocalService(nodeID string) (*LocalService, error) {
	services := scServices{}
	return &LocalService{services: append(services, &service{id: nodeID})}, nil
}

// WatchServices maintains consistency over all services
func (s *LocalService) WatchServices(cb func(scServices)) error {
	cb(s.services)
	return nil
}
