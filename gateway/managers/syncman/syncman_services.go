package syncman

import (
	"context"
	"fmt"
	"net/http"

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

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceRemoteService, service)
	if projectConfig.RemoteService == nil {
		projectConfig.RemoteService = config.Services{resourceID: value}
	} else {
		projectConfig.RemoteService[resourceID] = value
	}

	if err := s.modules.SetRemoteServiceConfig(ctx, project, projectConfig.RemoteService); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.SetResource(ctx, resourceID, value); err != nil {
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

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceRemoteService, service)
	delete(projectConfig.RemoteService, resourceID)

	if err := s.modules.SetRemoteServiceConfig(ctx, project, projectConfig.RemoteService); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
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
		service, ok := projectConfig.RemoteService[config.GenerateResourceID(s.clusterID, project, config.ResourceRemoteService, serviceID)]
		if !ok {
			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("service with id (%s) does not exists", serviceID), nil, nil)
		}
		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.RemoteService {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}
