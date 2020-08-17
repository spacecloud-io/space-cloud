package integration

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetConfig sets the config of the integration manager
func (m *Manager) SetConfig(array config.Integrations) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Check if integration array is valid
	if err := m.adminMan.ValidateIntegrationSyncOperation(array); err != nil {
		m.config = map[string]*config.IntegrationConfig{}
		return err
	}

	// Reset existing config
	m.config = make(map[string]*config.IntegrationConfig, len(array))

	// Add the integration config on by one
	for _, i := range array {
		m.config[i.ID] = i
	}

	return nil
}
