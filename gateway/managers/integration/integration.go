package integration

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// Module is responsible for handling all integration related tasks
type Manager struct {
	lock sync.RWMutex

	adminMan adminManager

	config map[string]*config.IntegrationConfig
}

// New creates a new instance of the integration module
func New(adminMan adminManager) *Manager {
	return &Manager{adminMan: adminMan, config: map[string]*config.IntegrationConfig{}}
}

const (
	module    string = "integration"
	checkAuth string = "check-auth"
)
