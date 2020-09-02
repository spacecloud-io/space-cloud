package syncman

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetService adds a remote service
func (s *Manager) SetService(ctx context.Context, project, service string, value *config.Service, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = service
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if projectConfig.Modules.Services.Services == nil {
		projectConfig.Modules.Services.Services = config.Services{}
	}

	// Check timeout field and add default value if not present
	for _, val := range value.Endpoints {
		if val.Timeout == 0 {
			val.Timeout = int(10 * time.Second)
		}
	}

	projectConfig.Modules.Services.Services[service] = value

	if err := s.modules.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteService deletes a remotes service
func (s *Manager) DeleteService(ctx context.Context, project, service string, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	delete(projectConfig.Modules.Services.Services, service)

	if err := s.modules.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetServices gets a remotes service
func (s *Manager) GetServices(ctx context.Context, project, serviceID string, params model.RequestParams) (int, []interface{}, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if serviceID != "*" {
		service, ok := projectConfig.Modules.Services.Services[serviceID]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("service with id (%s) does not exists", serviceID), nil, nil)
		}
		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Services.Services {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}
