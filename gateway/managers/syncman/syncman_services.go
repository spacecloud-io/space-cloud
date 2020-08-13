package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetService adds a remote service
func (s *Manager) SetService(ctx context.Context, project, service string, value *config.Service, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	value.ID = service
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if projectConfig.Modules.Services.Services == nil {
		projectConfig.Modules.Services.Services = config.Services{}
	}
	projectConfig.Modules.Services.Services[service] = value

	if err := s.modules.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		logrus.Errorf("error setting services config - %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteService deletes a remotes service
func (s *Manager) DeleteService(ctx context.Context, project, service string, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	delete(projectConfig.Modules.Services.Services, service)

	if err := s.modules.SetServicesConfig(project, projectConfig.Modules.Services); err != nil {
		logrus.Errorf("error setting services config - %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetServices gets a remotes service
func (s *Manager) GetServices(ctx context.Context, project, serviceID string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if serviceID != "*" {
		service, ok := projectConfig.Modules.Services.Services[serviceID]
		if !ok {
			return http.StatusBadRequest, nil, fmt.Errorf("serviceID (%s) not present in config", serviceID)
		}
		return http.StatusOK, []interface{}{service}, nil
	}

	services := []interface{}{}
	for _, value := range projectConfig.Modules.Services.Services {
		services = append(services, value)
	}
	return http.StatusOK, services, nil
}
