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

	clusterID string
	nodeID    string

	// Global Modules
	GlobalMods *global.Global

	// Managers
	Managers *managers.Managers
}

// New creates a new modules instance
func New(_, clusterID, nodeID string, managers *managers.Managers, globalMods *global.Global) (*Modules, error) {
	return &Modules{
		blocks:     map[string]*Module{},
		clusterID:  clusterID,
		nodeID:     nodeID,
		GlobalMods: globalMods,
		Managers:   managers,
	}, nil
}

// SetInitialProjectConfig sets the config all modules
func (m *Modules) SetInitialProjectConfig(ctx context.Context, projects config.Projects) error {
	for projectID, project := range projects {
		module, err := m.loadModule(projectID)
		if err != nil {
			module, err = m.newModule(project.ProjectConfig)
			if err != nil {
				return err
			}
		}

		if err := module.SetInitialProjectConfig(ctx, config.Projects{projectID: project}); err != nil {
			return err
		}
	}
	return nil
}

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(ctx context.Context, config *config.ProjectConfig) error {
	module, err := m.loadModule(config.ID)
	if err != nil {
		module, err = m.newModule(config)
		if err != nil {
			return err
		}
	}
	return module.SetProjectConfig(ctx, config)
}

// SetDatabaseConfig sets the config of db, auth, schema and realtime modules
func (m *Modules) SetDatabaseConfig(ctx context.Context, projectID string, databaseConfigs config.DatabaseConfigs, schemaConfigs config.DatabaseSchemas, ruleConfigs config.DatabaseRules, prepConfigs config.DatabasePreparedQueries) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetDatabaseConfig(ctx, projectID, databaseConfigs, schemaConfigs, ruleConfigs, prepConfigs)
}

// SetDatabaseSchemaConfig sets database schema config
func (m *Modules) SetDatabaseSchemaConfig(ctx context.Context, projectID string, schemaConfigs config.DatabaseSchemas) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetDatabaseSchemaConfig(ctx, projectID, schemaConfigs)
}

// SetDatabaseRulesConfig set database rules of db module
func (m *Modules) SetDatabaseRulesConfig(ctx context.Context, projectID string, ruleConfigs config.DatabaseRules) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetDatabaseRulesConfig(ctx, ruleConfigs)
}

// SetDatabasePreparedQueryConfig set prepared config of database moudle
func (m *Modules) SetDatabasePreparedQueryConfig(ctx context.Context, projectID string, prepConfigs config.DatabasePreparedQueries) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetDatabasePreparedQueryConfig(ctx, prepConfigs)
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(ctx context.Context, projectID string, fileStore *config.FileStoreConfig) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetFileStoreConfig(ctx, projectID, fileStore)
}

// SetFileStoreSecurityRuleConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreSecurityRuleConfig(ctx context.Context, projectID string, fileStoreRules config.FileStoreRules) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	module.SetFileStoreSecurityRuleConfig(ctx, projectID, fileStoreRules)
	return nil
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(ctx context.Context, projectID string, eventingConfig *config.EventingConfig, secureObj config.EventingRules, eventingSchemas config.EventingSchemas, eventingTriggers config.EventingTriggers) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetEventingConfig(ctx, projectID, eventingConfig, secureObj, eventingSchemas, eventingTriggers)
}

// SetEventingSchemaConfig sets the config of eventing module
func (m *Modules) SetEventingSchemaConfig(ctx context.Context, projectID string, eventingSchemas config.EventingSchemas) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetEventingSchemaConfig(ctx, eventingSchemas)
}

// SetEventingTriggerConfig sets the config of eventing module
func (m *Modules) SetEventingTriggerConfig(ctx context.Context, projectID string, eventingTriggers config.EventingTriggers) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetEventingTriggerConfig(ctx, eventingTriggers)
}

// SetEventingRuleConfig sets the config of eventing module
func (m *Modules) SetEventingRuleConfig(ctx context.Context, projectID string, secureObj config.EventingRules) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetEventingRuleConfig(ctx, secureObj)
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(ctx context.Context, projectID string, auth config.Auths) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetUsermanConfig(ctx, projectID, auth)
}

// SetLetsencryptConfig set the config of letsencrypt module
func (m *Modules) SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetLetsencryptConfig(ctx, projectID, c)
}

// SetIngressRouteConfig set the config of routing module
func (m *Modules) SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetIngressRouteConfig(ctx, projectID, routes)
}

// SetIngressGlobalRouteConfig set config of routing module
func (m *Modules) SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetIngressGlobalRouteConfig(ctx, projectID, c)
}

// SetRemoteServiceConfig set config of functions module
func (m *Modules) SetRemoteServiceConfig(ctx context.Context, projectID string, services config.Services) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetRemoteServiceConfig(ctx, projectID, services)
}

func (m *Modules) projects() *config.Config {
	m.lock.RLock()
	defer m.lock.RUnlock()

	projects := make(config.Projects)
	for id := range m.blocks {
		projects[id] = &config.Project{ProjectConfig: &config.ProjectConfig{ID: id}}
	}
	return &config.Config{Projects: projects}
}

func (m *Modules) Delete(projectID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if block, p := m.blocks[projectID]; p {
		// Close all the modules here
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of auth module", nil)
		block.auth.CloseConfig()

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

func (m *Modules) newModule(config *config.ProjectConfig) (*Module, error) {
	projects := m.projects()
	m.lock.Lock()
	defer m.lock.Unlock()

	if ok := m.Managers.Admin().ValidateProjectSyncOperation(projects, config); !ok {
		helpers.Logger.LogWarn("", "Cannot create new project. Upgrade your plan", nil)
		return nil, errors.New("upgrade your plan to create new project")
	}

	module, err := newModule(config.ID, m.clusterID, m.nodeID, m.Managers, m.GlobalMods)
	if err != nil {
		return nil, err
	}

	m.blocks[config.ID] = module
	return module, nil
}
