package admin

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetIntegrationToken returns the admin token required by an intergation
func (m *Manager) GetIntegrationToken(id string) (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.createToken(map[string]interface{}{"id": id, "role": "integration"})
}

func (m *Manager) ParseLicense(license string) (map[string]interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Fetch the license key if it isn't already present
	if m.publicKey == nil {
		if err := m.fetchPublicKeyWithoutLock(); err != nil {
			return nil, err
		}
	}

	return m.parseLicenseToken(license)
}

func (m *Manager) ValidateIntegrationSyncOperation(integrations config.Integrations) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Fetch the license key if it isn't already present
	if m.publicKey == nil {
		if err := m.fetchPublicKeyWithoutLock(); err != nil {
			return err
		}
	}

	// Iterate over each integration
	for _, i := range integrations {
		obj, err := m.parseLicenseToken(i.License)
		if err != nil {
			m.config.Integrations = removeIntegration(m.config.Integrations, i.ID)
			return err
		}

		// Return error if license does not belong to integration
		if obj["id"] != i.ID {
			m.config.Integrations = removeIntegration(m.config.Integrations, i.ID)
			return utils.LogError(fmt.Sprintf("Integration (%s) has an invlaid license", i.ID), "admin", "validate-integration", nil)
		}

		// Check if the level is larger than the licensed level
		if obj["level"].(float64) > m.quotas.IntegrationLevel {
			m.config.Integrations = removeIntegration(m.config.Integrations, i.ID)
			return utils.LogError(fmt.Sprintf("Integration (%s) cannot be used with the current plan", i.ID), "admin", "validate-integration", nil)
		}
	}

	return nil
}

func removeIntegration(arr config.Integrations, id string) config.Integrations {
	length := len(arr)
	for index, integrationConfig := range arr {
		if integrationConfig.ID == id {
			arr[index] = arr[length-1]
			return arr[:length-1]
		}
	}
	return arr
}
