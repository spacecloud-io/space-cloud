package deploy

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/deploy/kubernetes"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is the main object for handling all deployments
type Module struct {
	lock   sync.RWMutex
	config *config.Deploy
	driver Driver
	client http.Client
}

// Driver is the interface every deployment driver must implement
type Driver interface {
	Deploy(ctx context.Context, conf *model.Deploy) error
}

// New creates a new instance of the deploy module
func New() *Module {
	m := new(Module)
	m.client = http.Client{}
	return m
}

// SetConfig initialises the driver for the deployment module
func (m *Module) SetConfig(c *config.Deploy) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Exit if module isn't enabled
	if c == nil || !c.Enabled {
		return nil
	}

	// Store the config
	m.config = c

	// Sign in to the registry
	if err := m.signIn(); err != nil {
		return err
	}

	// Create a new instance of the appropriate driver
	var err error
	switch c.Orchestrator {
	case utils.Kubernetes:
		m.driver, err = kubernetes.New(&c.Registry)
	default:
		err = errors.New("Deploy: Invalid orchestrator provided")
	}

	return err
}
