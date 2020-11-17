package admin

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/spaceuptech/helpers"

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
func (m *Manager) IsTokenValid(ctx context.Context, token, resource, op string, attr map[string]string) (model.RequestParams, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.isProd {
		return model.RequestParams{}, nil
	}

	claims, err := m.parseToken(ctx, token)
	if err != nil {
		return model.RequestParams{}, err
	}

	// Check if its an integration request and return the integration response if its an integration request
	res := m.integrationMan.HandleConfigAuth(ctx, resource, op, claims, attr)
	if res.CheckResponse() && res.Error() != nil {
		return model.RequestParams{}, res.Error()
	}

	// Otherwise just return nil for backward compatibility
	return model.RequestParams{Resource: resource, Op: op, Attributes: attr, Claims: claims}, nil
}

// CheckIfAdmin simply checks the token
func (m *Manager) CheckIfAdmin(ctx context.Context, token string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.isProd {
		return nil
	}

	claims, err := m.parseToken(ctx, token)
	if err != nil {
		return err
	}

	// Check if role is admin
	role, p := claims["role"]
	if !p {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid token provided. Claim `role` is absent.", nil, nil)
	}

	if !strings.Contains(role.(string), "admin") {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Only admins are authorised to make this request.", nil, nil)
	}

	return nil
}

// IsDBConfigValid checks if the database config is valid
func (m *Manager) IsDBConfigValid(config config.DatabaseConfigs) error {
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
func (m *Manager) ValidateProjectSyncOperation(c *config.Config, project *config.ProjectConfig) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Allow if project is an integration
	if _, p := m.integrations.Get(project.ID); p {
		return true
	}

	var totalProjects int

	for _, p := range c.Projects {
		// Return true if the project already exists in
		if p.ProjectConfig.ID == project.ID {
			return true
		}

		// Increment count if it isn't an integration
		if m.integrations == nil {
			totalProjects++
			continue
		}

		if _, p := m.integrations.Get(p.ProjectConfig.ID); !p {
			totalProjects++
		}
	}

	return totalProjects < m.quotas.MaxProjects
}

// RefreshToken is used to create a new token based on an existing one
func (m *Manager) RefreshToken(ctx context.Context, token string) (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Parse the token to get userID and userRole
	tokenClaims, err := m.parseToken(ctx, token)
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
func (m *Manager) GetSessionID() (string, error) {
	if licenseMode == licenseModeOffline {
		return m.getOfflineLicenseSessionID(), nil
	}
	if m.isEnterpriseMode() {
		licenseObj, err := m.decryptLicense(m.license.License)
		if err != nil {
			return "", helpers.Logger.LogError("get-session-id", "Unable to validate license key", err, nil)
		}
		return licenseObj.SessionID, nil
	}
	return selectRandomSessionID(m.services), nil // first time license renewal
}

func (m *Manager) GetEnterpriseClusterID() string {
	return m.license.LicenseKey
}

// GetSecret returns the admin secret
func (m *Manager) GetSecret() string {
	return m.user.Secret
}

// GetPermissions returns the permissions the user has. The permissions is for the format `projectId:resource`.
// This only applies to the config level endpoints.
func (m *Manager) GetPermissions(ctx context.Context, params model.RequestParams) (int, interface{}, error) {
	hookResponse := m.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		return hookResponse.Status(), hookResponse.Result(), nil
	}

	return http.StatusOK, []interface{}{map[string]interface{}{"project": "*", "resource": "*", "verb": "*"}}, nil
}
