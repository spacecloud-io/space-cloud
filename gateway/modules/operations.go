package modules

import (
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
)

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(c *config.Config, le *letsencrypt.LetsEncrypt, ingressRouting *routing.Routing) {
	if c.Projects != nil && len(c.Projects) > 0 {
		p := c.Projects[0]

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

		logrus.Debugln("Setting config of db module")
		if err := m.db.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting db module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of schema module")
		if err := m.schema.SetConfig(p.Modules.Crud, p.ID); err != nil {
			logrus.Errorf("error setting schema module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of auth module")
		if err := m.auth.SetConfig(p.ID, p.Secrets, p.AESKey, p.Modules.Crud, p.Modules.FileStore, p.Modules.Services, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting auth module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of functions module")
		if err := m.functions.SetConfig(p.ID, p.Modules.Services); err != nil {
			logrus.Errorf("error setting remote services module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of user management module")
		m.user.SetConfig(p.Modules.Auth)

		logrus.Debugln("Setting config of file storage module")
		if err := m.file.SetConfig(p.ID, p.Modules.FileStore); err != nil {
			logrus.Errorf("error setting filestore module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of eventing module")
		if err := m.eventing.SetConfig(p.ID, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting eventing module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of realtime module")
		if err := m.realtime.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting realtime module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of graphql module")
		m.graphql.SetConfig(p.ID)

		logrus.Debugln("Setting config of lets encrypt module")
		if err := le.SetProjectDomains(p.ID, p.Modules.LetsEncrypt); err != nil {
			logrus.Errorf("error setting letsencypt module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of ingress routing module")
		if err := ingressRouting.SetProjectRoutes(p.ID, p.Modules.Routes); err != nil {
			logrus.Errorf("error setting routing module config - %s", err.Error())
		}
	}
}

// SetGlobalConfig sets the auth secret and AESKey
func (m *Modules) SetGlobalConfig(projectID string, secrets []*config.Secret, aesKey string) error {
	m.auth.SetSecrets(secrets)
	return m.auth.SetAESKey(aesKey)
}

// SetCrudConfig sets the config of db, auth, schema and realtime modules
func (m *Modules) SetCrudConfig(projectID string, crudConfig config.Crud) error {
	logrus.Debugln("Setting config of db module")
	if err := m.db.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting db module config - %s", err.Error())
		return err
	}
	logrus.Debugln("Setting config of auth module")
	m.auth.SetCrudConfig(projectID, crudConfig)

	logrus.Debugln("Setting config of schema module")
	if err := m.schema.SetConfig(crudConfig, projectID); err != nil {
		logrus.Errorf("error setting schema module config - %s", err.Error())
		return err
	}
	if err := m.realtime.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting realtime module config - %s", err.Error())
		return err
	}
	logrus.Debugln("Setting config of file storage module")
	return nil
}

// SetServicesConfig sets the config of auth and functions modules
func (m *Modules) SetServicesConfig(projectID string, services *config.ServicesModule) error {
	logrus.Debugln("Setting config of auth module")
	m.auth.SetServicesConfig(projectID, services)

	logrus.Debugln("Setting config of remote services module")
	return m.functions.SetConfig(projectID, services)
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(projectID string, fileStore *config.FileStore) error {
	logrus.Debugln("Setting config of auth module")
	m.auth.SetFileStoreConfig(projectID, fileStore)

	logrus.Debugln("Setting config of file storage module")
	if err := m.file.SetConfig(projectID, fileStore); err != nil {
		logrus.Errorf("error setting filestore module config - %s", err.Error())
		return err
	}
	return nil
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	logrus.Debugln("Setting config of eventing module")
	if err := m.eventing.SetConfig(projectID, eventingConfig); err != nil {
		logrus.Errorf("error setting eventing module config - %s", err.Error())
		return err
	}
	m.auth.SetEventingConfig(eventingConfig.SecurityRules)
	return nil
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(projectID string, auth config.Auth) {
	logrus.Debugln("Setting config of user management module")
	m.user.SetConfig(auth)
}
