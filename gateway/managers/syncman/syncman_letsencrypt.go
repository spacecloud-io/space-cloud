package syncman

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetProjectLetsEncryptDomains sets a projects whitelisted domains
func (s *Manager) SetProjectLetsEncryptDomains(ctx context.Context, project string, c config.LetsEncrypt) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	// Update the projects domains
	projectConfig.Modules.LetsEncrypt = c
	if err := s.modules.LetsEncrypt().SetProjectDomains(project, c); err != nil {
		logrus.Errorf("error setting letsencrypt project domains - %s", err.Error())
		return err
	}
	// Persist the config
	return s.setProject(ctx, projectConfig)
}

// GetLetsEncryptConfig returns the letsencrypt config for the particular project
func (s *Manager) GetLetsEncryptConfig(project string, params model.RequestParams) (config.LetsEncrypt, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return config.LetsEncrypt{}, err
	}

	return projectConfig.Modules.LetsEncrypt, nil
}
