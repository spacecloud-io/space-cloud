package admin

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
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

// IsDBConfigValid checks if the database config is valid
func (m *Manager) IsDBConfigValid(config config.Crud) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if len(config) > m.quotas.MaxDatabases {
		return fmt.Errorf("plan needs to be upgraded to use more than %d databases", m.quotas.MaxDatabases)
	}

	return nil
}

// ValidateProjectSyncOperation validates if an operation is permitted based on the mode
func (m *Manager) ValidateProjectSyncOperation(projects []string, projectID string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, p := range projects {
		if p == projectID {
			return true
		}
	}

	if len(projects) < m.quotas.MaxProjects {
		return true
	}

	return false
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

func (m *Manager) IsEnterpriseMode() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.isEnterpriseMode()
}
