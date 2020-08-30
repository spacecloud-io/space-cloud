package modules

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers"
	"github.com/spaceuptech/space-cloud/gateway/modules/global"
)

// Modules is an object that sets up the modules
type Modules struct {
	lock   sync.RWMutex
	blocks map[string]*Module

	nodeID string

	// Global Modules
	GlobalMods *global.Global

	// Managers
	Managers *managers.Managers
}

// New creates a new modules instance
func New(nodeID string, managers *managers.Managers, globalMods *global.Global) (*Modules, error) {
	return &Modules{
		blocks:     map[string]*Module{},
		nodeID:     nodeID,
		GlobalMods: globalMods,
		Managers:   managers,
	}, nil
}

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(config *config.Project) error {
	module, err := m.loadModule(config.ID)
	if err != nil {
		module, err = m.newModule(config)
		if err != nil {
			return err
		}
	}
	_ = module.SetProjectConfig(config)
	return nil
}

// SetGlobalConfig sets the auth secret and AESKey
func (m *Modules) SetGlobalConfig(projectID, secretSource string, secret []*config.Secret, aesKey string) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetGlobalConfig(projectID, secretSource, secret, aesKey)
}

// SetCrudConfig sets the config of db, auth, schema and realtime modules
func (m *Modules) SetCrudConfig(projectID string, crudConfig config.Crud) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetCrudConfig(projectID, crudConfig)
}

// SetServicesConfig sets the config of auth and functions modules
func (m *Modules) SetServicesConfig(projectID string, services *config.ServicesModule) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetServicesConfig(projectID, services)
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(projectID string, fileStore *config.FileStore) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetFileStoreConfig(projectID, fileStore)
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetEventingConfig(projectID, eventingConfig)
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(projectID string, auth config.Auth) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	_ = module.SetUsermanConfig(projectID, auth)
	return nil
}

func (m *Modules) projects() *config.Config {
	m.lock.RLock()
	defer m.lock.RUnlock()

	c := &config.Config{Projects: []*config.Project{}}
	for id := range m.blocks {
		c.Projects = append(c.Projects, &config.Project{ID: id})
	}
	return c
}

func (m *Modules) Delete(projectID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if block, p := m.blocks[projectID]; p {
		// Close all the modules here
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of db module", nil)
		if err := block.db.CloseConfig(); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing db module config", err, map[string]interface{}{"project": projectID})
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of filestore module", nil)
		if err := block.file.CloseConfig(); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing filestore module config", err, map[string]interface{}{"project": projectID})
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of eventing module", nil)
		if err := block.eventing.CloseConfig(); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing eventing module config", err, map[string]interface{}{"project": projectID})
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of realtime module", nil)
		if err := block.realtime.CloseConfig(); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error closing realtime module config", err, map[string]interface{}{"project": projectID})
		}
	}

	delete(m.blocks, projectID)

	// Remove config from global modules
	_ = m.LetsEncrypt().DeleteProjectDomains(projectID)
	m.Routing().DeleteProjectRoutes(projectID)
}

func (m *Modules) loadModule(projectID string) (*Module, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if module, p := m.blocks[projectID]; p {
		return module, nil
	}

	return nil, fmt.Errorf("project (%s) not found in server state", projectID)
}

func (m *Modules) newModule(config *config.Project) (*Module, error) {
	projects := m.projects()
	m.lock.Lock()
	defer m.lock.Unlock()

	if ok := m.Managers.Admin().ValidateProjectSyncOperation(projects, config); !ok {
		helpers.Logger.LogWarn("", "Cannot create new project. Upgrade your plan", nil)
		return nil, errors.New("upgrade your plan to create new project")
	}

	module := newModule(m.nodeID, m.Managers, m.GlobalMods)
	m.blocks[config.ID] = module
	return module, nil
}
