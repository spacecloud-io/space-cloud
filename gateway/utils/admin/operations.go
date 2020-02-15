package admin

import (
	"errors"
	"net/http"

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

	if m.admin.Operation.Mode == 0 && len(config) > 2 {
		return errors.New("community edition can have a maximum of 2 dbs in a single project")
	}

	return nil
}

// ValidateSyncOperation validates if an operation is permitted based on the mode
func (m *Manager) ValidateSyncOperation(projects []string, project *config.Project) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, p := range projects {
		if p == project.ID {
			return true
		}
	}

	maxProjects := 1
	if m.admin.Operation.Mode == 1 {
		maxProjects = 3
	} else if m.admin.Operation.Mode == 2 {
		maxProjects = 5
	}

	if len(projects) < maxProjects {
		return true
	}

	return false
}

// IsAdminOpAuthorised checks if the admin operation is authorised.
// TODO add scope level restrictions as well
func (m *Manager) IsAdminOpAuthorised(token, scope string) (int, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.isProd {
		return http.StatusOK, nil
	}

	auth, err := m.parseToken(token)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	user, p := auth["id"]
	if !p {
		return http.StatusUnauthorized, errors.New("Invalid Token")
	}

	if user == utils.InternalUserID {
		return http.StatusOK, nil
	}

	for _, u := range m.admin.Users {
		if u.User == user {

			// Allow full access for scope name `all`
			if _, p := u.Scopes["all"]; p {
				return http.StatusOK, nil
			}

			// Check if scope is present
			if _, p := u.Scopes[scope]; p {
				return http.StatusOK, nil
			}

			break
		}
	}

	return http.StatusForbidden, errors.New("You are not authorized to make this request")
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

func (m *Manager) reduceOpMode() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	m.admin.Operation.Mode = 0
}
