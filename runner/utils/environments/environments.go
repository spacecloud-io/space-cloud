package environments

import (
	"sync"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// Manager manages the various environments in the runner
type Manager struct {
	sync.RWMutex

	// For internal use
	config *Config

	// For tracking environments
	envs model.Environments
}

// New creates a new instance of the environment manager
func New(c *Config) *Manager {
	return &Manager{config: c}
}
