package syncman

import (
	"context"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetProjectLetsEncryptDomains sets a projects whitelisted domains
func (s *Manager) SetProjectLetsEncryptDomains(ctx context.Context, project string, c config.LetsEncrypt, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Update the projects domains
	projectConfig.Modules.LetsEncrypt = c
	if err := s.modules.LetsEncrypt().SetProjectDomains(project, c); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting letsencrypt project domains", err, nil)
	}

	if err := s.setProject(ctx, projectConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetLetsEncryptConfig returns the letsencrypt config for the particular project
func (s *Manager) GetLetsEncryptConfig(ctx context.Context, project string, params model.RequestParams) (int, interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, config.LetsEncrypt{}, err
	}

	return http.StatusOK, projectConfig.Modules.LetsEncrypt, nil
}
