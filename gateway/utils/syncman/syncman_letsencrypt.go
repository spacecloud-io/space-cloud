package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
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

	return s.setProject(ctx, projectConfig)
}
