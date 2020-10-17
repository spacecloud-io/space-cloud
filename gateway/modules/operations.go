package modules

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetInitialProjectConfig sets the config all modules
func (m *Modules) SetInitialProjectConfig(ctx context.Context, projects config.Projects) error {
	for projectID, project := range projects {
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of db module", nil)
		if err := m.db.SetConfig(projectID, project.DatabaseConfigs); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set db module config", err, nil)
		}
		if err := m.db.SetSchemaConfig(ctx, project.DatabaseSchemas); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set schema db module config", err, nil)
		}
		if err := m.db.SetPreparedQueryConfig(ctx, project.DatabasePreparedQueries); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set db prepared query module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of schema module", nil)
		if err := m.schema.SetConfig(project.DatabaseSchemas, projectID); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set schema module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of auth module", nil)
		if err := m.auth.SetConfig(ctx, project.FileStoreConfig.StoreType, project.ProjectConfig, project.DatabaseRules, project.DatabasePreparedQueries, project.FileStoreRules, project.RemoteService, project.EventingRules); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set auth module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of functions module", nil)
		if err := m.functions.SetConfig(projectID, project.RemoteService); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set remote services module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of user management module", nil)
		m.user.SetConfig(project.Auths)

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of file storage module", nil)
		if err := m.file.SetConfig(projectID, project.FileStoreConfig); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set filestore module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of eventing module", nil)
		if err := m.eventing.SetConfig(projectID, project.EventingConfig); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set eventing module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of realtime module", nil)
		if err := m.realtime.SetConfig(projectID, project.DatabaseConfigs, project.DatabaseRules, project.DatabaseSchemas); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set realtime module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of graphql module", nil)
		m.graphql.SetConfig(projectID)

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of lets encrypt module", nil)
		if err := m.GlobalMods.LetsEncrypt().SetProjectDomains(projectID, project.LetsEncrypt); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set letsencypt module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of ingress routing module", nil)
		if err := m.GlobalMods.Routing().SetProjectRoutes(projectID, project.IngressRoutes); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set routing module config", err, nil)
		}
		m.GlobalMods.Routing().SetGlobalConfig(project.IngressGlobal)
	}
	return nil
}

// SetProjectConfig set project config
func (m *Modules) SetProjectConfig(ctx context.Context, p *config.ProjectConfig) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting project config", nil)
	if err := m.auth.SetProjectConfig(p); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set project config", err, nil)
	}
	m.graphql.SetConfig(p.ID)
	return nil
}

// SetDatabaseConfig sets the config of db, auth, schema and realtime modules
func (m *Modules) SetDatabaseConfig(ctx context.Context, projectID string, databaseConfigs config.DatabaseConfigs) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of db module", nil)
	if err := m.db.SetConfig(projectID, databaseConfigs); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set db module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of realtime module", nil)
	m.realtime.SetDatabaseConfig(databaseConfigs)
	return nil
}

// SetDatabaseSchemaConfig sets database schema config
func (m *Modules) SetDatabaseSchemaConfig(ctx context.Context, projectID string, schemaConfigs config.DatabaseSchemas) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of db schema in db module", nil)
	if err := m.db.SetSchemaConfig(ctx, schemaConfigs); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set db schema in db module", err, nil)
	}
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of schema module", nil)
	if err := m.schema.SetConfig(schemaConfigs, projectID); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set schema module config", err, nil)
	}
	m.realtime.SetDatabaseSchemas(schemaConfigs)
	return nil
}

// SetDatabaseRulesConfig set database rules of db module
func (m *Modules) SetDatabaseRulesConfig(ctx context.Context, ruleConfigs config.DatabaseRules) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of db rule in db module", nil)
	m.auth.SetDatabaseRules(ruleConfigs)
	m.realtime.SetDatabaseRules(ruleConfigs)
	return nil
}

// SetDatabasePreparedQueryConfig set prepared config of database moudle
func (m *Modules) SetDatabasePreparedQueryConfig(ctx context.Context, prepConfigs config.DatabasePreparedQueries) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of db prepared query in db module", nil)
	if err := m.db.SetPreparedQueryConfig(ctx, prepConfigs); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set db prepared query in db module", err, nil)
	}
	m.auth.SetDatabasePreparedQueryRules(prepConfigs)
	return nil
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(ctx context.Context, projectID string, fileStore *config.FileStoreConfig) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of file storage module", nil)
	if err := m.file.SetConfig(projectID, fileStore); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set filestore module config", err, nil)
	}
	m.auth.SetFileStoreType(fileStore.StoreType)
	return nil
}

// SetFileStoreSecurityRuleConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreSecurityRuleConfig(ctx context.Context, projectID string, fileStoreRules config.FileStoreRules) {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of file store rules in auth module", nil)
	m.auth.SetFileStoreRules(fileStoreRules)
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(ctx context.Context, projectID string, eventingConfig *config.EventingConfig) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of eventing module", nil)
	if err := m.eventing.SetConfig(projectID, eventingConfig); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set eventing module config", err, nil)
	}
	return nil
}

// SetEventingSchemaConfig sets the config of eventing module
func (m *Modules) SetEventingSchemaConfig(ctx context.Context, eventingSchemas config.EventingSchemas) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting schema config of eventing module", nil)
	return m.eventing.SetSchemaConfig(eventingSchemas)
}

// SetEventingTriggerConfig sets the config of eventing module
func (m *Modules) SetEventingTriggerConfig(ctx context.Context, eventingTriggers config.EventingTriggers) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting trigger config of eventing module", nil)
	return m.eventing.SetTriggerConfig(eventingTriggers)
}

// SetEventingRuleConfig sets the config of eventing module
func (m *Modules) SetEventingRuleConfig(ctx context.Context, secureObj config.EventingRules) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting rules config of eventing module", nil)
	if err := m.eventing.SetSecurityRuleConfig(secureObj); err != nil {
		return err
	}
	m.auth.SetEventingRules(secureObj)
	return nil
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(ctx context.Context, projectID string, auth config.Auths) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of user management module", nil)
	m.user.SetConfig(auth)
	return nil
}

// SetLetsencryptConfig set the config of letsencrypt module
func (m *Modules) SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting letsencrypt config of project", nil)
	return m.GlobalMods.LetsEncrypt().SetProjectDomains(projectID, c)
}

// SetIngressRouteConfig set the config of routing module
func (m *Modules) SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of routing module", nil)
	if err := m.GlobalMods.Routing().SetProjectRoutes(projectID, routes); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set routing module config", err, nil)
	}
	return nil
}

// SetIngressGlobalRouteConfig set config of routing module
func (m *Modules) SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of global routing", nil)
	m.GlobalMods.Routing().SetGlobalConfig(c)
	return nil
}

// SetRemoteServiceConfig set config of functions module
func (m *Modules) SetRemoteServiceConfig(ctx context.Context, projectID string, services config.Services) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of auth module", nil)
	m.auth.SetRemoteServiceConfig(services)

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of remote service module", nil)
	if err := m.functions.SetConfig(projectID, services); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set config of remote service module", err, nil)
	}
	return nil
}
