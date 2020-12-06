package integration

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// Module is responsible for handling all integration related tasks
type Manager struct {
	lock sync.RWMutex

	adminMan adminManager

	integrationConfig     config.Integrations
	integrationHookConfig config.IntegrationHooks
}

// New creates a new instance of the integration module
func New(adminMan adminManager) *Manager {
	return &Manager{adminMan: adminMan, integrationConfig: make(config.Integrations), integrationHookConfig: make(config.IntegrationHooks)}
}
