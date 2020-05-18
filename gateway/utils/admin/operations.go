package admin

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetInternalAccessToken returns the token that can be used internally by Space Cloud
func (m *Manager) GetInternalAccessToken() (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.createToken(map[string]interface{}{"id": utils.InternalUserID})
}

// IsTokenValid checks if the token is valid
func (m *Manager) IsTokenValid(token string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.isProd {
		return nil
	}

	_, err := m.parseToken(token)
	return err
}

// ValidateSyncOperation validates if an operation is permitted based on the mode
func (m *Manager) ValidateSyncOperation(c *config.Config, project *config.Project) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, p := range c.Projects {
		if p.ID == project.ID {
			return true
		}
	}

	return len(c.Projects) < m.quotas.MaxProjects
}

// RefreshToken is used to create a new token based on an existing one
func (m *Manager) RefreshToken(token string) (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Parse the token to get userID and userRole
	tokenClaims, err := m.parseToken(token)
	if err != nil {
		return "", err
	}
	// Create a new token
	newToken, err := m.createToken(tokenClaims)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

// GetQuotas gets number of projects & databases that can be created
func (m *Manager) GetQuotas() *model.UsageQuotas {
	return &m.quotas
}

// GetCredentials gets user name & pass
func (m *Manager) GetCredentials() map[string]interface{} {
	return map[string]interface{}{"user": m.user.User, "pass": m.user.Pass}
}

// GetClusterID returns the cluster id
func (m *Manager) GetClusterID() string {
	return m.clusterID
}
