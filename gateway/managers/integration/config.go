package integration

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetConfig sets the config of the integration manager
func (m *Manager) SetConfig(integrations config.Integrations, integrationHooks config.IntegrationHooks) error {
	if err := m.SetIntegrations(integrations); err != nil {
		return err
	}

	m.SetIntegrationHooks(integrationHooks)
	return nil
}

func (m *Manager) SetIntegrations(integrations config.Integrations) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Check if integration are valid
	if err := m.adminMan.ValidateIntegrationSyncOperation(integrations); err != nil {
		m.integrationConfig = map[string]*config.IntegrationConfig{}
		return err
	}

	m.integrationConfig = integrations
	return nil
}

func (m *Manager) SetIntegrationHooks(integrationHooks config.IntegrationHooks) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.integrationHookConfig = integrationHooks
}
