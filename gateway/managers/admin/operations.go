package admin

import (
	"fmt"

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
func (m *Manager) IsTokenValid(token, resource, op string, attr map[string]string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.isProd {
		return nil
	}

	_, err := m.parseToken(token)
	return err
}

// CheckToken simply checks the token
func (m *Manager) CheckToken(token string) error {
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

	// Only count the length of enabled databases
	var length int
	for _, c := range config {
		if c.Enabled {
			length++
		}
	}

	if length > m.quotas.MaxDatabases {
		return fmt.Errorf("plan needs to be upgraded to use more than %d databases", m.quotas.MaxDatabases)
	}

	return nil
}

// ValidateProjectSyncOperation validates if an operation is permitted based on the mode
func (m *Manager) ValidateProjectSyncOperation(c *config.Config, project *config.Project) bool {
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

func (m *Manager) IsRegistered() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.isRegistered()
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
func (m *Manager) GetSessionID() string {
	return m.sessionID
}

func (m *Manager) GetEnterpriseClusterID() string {
	return m.config.ClusterID
}
