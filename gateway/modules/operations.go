package modules

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(c *config.Project) error {
	p := c

	if p.Modules == nil {
		p.Modules = &config.Modules{
			FileStore:   &config.FileStore{},
			Services:    &config.ServicesModule{},
			Auth:        map[string]*config.AuthStub{},
			Crud:        map[string]*config.CrudStub{},
			Routes:      []*config.Route{},
			LetsEncrypt: config.LetsEncrypt{WhitelistedDomains: []string{}},
		}
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of db module", nil)
	if err := m.db.SetConfig(p.ID, p.Modules.Crud); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set db module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of schema module", nil)
	if err := m.schema.SetConfig(p.Modules.Crud, p.ID); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set schema module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of auth module", nil)
	if err := m.auth.SetConfig(p.ID, p.SecretSource, p.Secrets, p.AESKey, p.Modules.Crud, p.Modules.FileStore, p.Modules.Services, &p.Modules.Eventing); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set auth module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of functions module", nil)
	if err := m.functions.SetConfig(p.ID, p.Modules.Services); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set remote services module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of user management module", nil)
	m.user.SetConfig(p.Modules.Auth)

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of file storage module", nil)
	if err := m.file.SetConfig(p.ID, p.Modules.FileStore); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set filestore module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of eventing module", nil)
	if err := m.eventing.SetConfig(p.ID, &p.Modules.Eventing); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set eventing module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of realtime module", nil)
	if err := m.realtime.SetConfig(p.ID, p.Modules.Crud); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set realtime module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of graphql module", nil)
	m.graphql.SetConfig(p.ID)

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of lets encrypt module", nil)
	if err := m.GlobalMods.LetsEncrypt().SetProjectDomains(p.ID, p.Modules.LetsEncrypt); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set letsencypt module config", err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of ingress routing module", nil)
	if err := m.GlobalMods.Routing().SetProjectRoutes(p.ID, p.Modules.Routes); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set routing module config", err, nil)
	}
	m.GlobalMods.Routing().SetGlobalConfig(p.Modules.GlobalRoutes)

	return nil
}

// SetCrudConfig sets the config of db, auth, schema and realtime modules
func (m *Modules) SetCrudConfig(projectID string, crudConfig config.Crud) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of db module", nil)
	if err := m.db.SetConfig(projectID, crudConfig); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set db module config", err, nil)
	}
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of auth module", nil)
	m.auth.SetCrudConfig(projectID, crudConfig)

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of schema module", nil)
	if err := m.schema.SetConfig(crudConfig, projectID); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set schema module config", err, nil)
	}
	if err := m.realtime.SetConfig(projectID, crudConfig); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set realtime module config", err, nil)
	}
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of file storage module", nil)
	return nil
}

// SetServicesConfig sets the config of auth and functions modules
func (m *Modules) SetServicesConfig(projectID string, services *config.ServicesModule) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of auth module", nil)
	m.auth.SetServicesConfig(projectID, services)

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of remote services module", nil)
	return m.functions.SetConfig(projectID, services)
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(projectID string, fileStore *config.FileStore) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of auth module", nil)
	m.auth.SetFileStoreConfig(projectID, fileStore)

	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of file storage module", nil)
	if err := m.file.SetConfig(projectID, fileStore); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set filestore module config", err, nil)
	}
	return nil
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of eventing module", nil)
	if err := m.eventing.SetConfig(projectID, eventingConfig); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set eventing module config", err, nil)
	}
	m.auth.SetEventingConfig(eventingConfig.SecurityRules)
	return nil
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(projectID string, auth config.Auth) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting config of user management module", nil)
	m.user.SetConfig(auth)
	return nil
}
