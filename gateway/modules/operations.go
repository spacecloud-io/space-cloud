package modules

import (
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
)

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(config *config.Config, le *letsencrypt.LetsEncrypt, ingressRouting *routing.Routing) {
	if config.Projects != nil && len(config.Projects) > 0 {
		p := config.Projects[0]

		logrus.Debugln("Setting config of crud module")
		if err := m.Crud.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting crud module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of schema module")
		if err := m.Schema.SetConfig(p.Modules.Crud, p.ID); err != nil {
			logrus.Errorf("error setting schema module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of auth module")
		if err := m.Auth.SetConfig(p.ID, p.Secret, p.AESkey, p.Modules.Crud, p.Modules.FileStore, p.Modules.Services, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting auth module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of functions module")
		m.Functions.SetConfig(p.ID, p.Modules.Services)

		logrus.Debugln("Setting config of user management module")
		m.User.SetConfig(p.Modules.Auth)

		logrus.Debugln("Setting config of file storage module")
		if err := m.File.SetConfig(p.Modules.FileStore); err != nil {
			logrus.Errorf("error setting filestore module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of eventing module")
		if err := m.Eventing.SetConfig(p.ID, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting eventing module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of realtime module")
		if err := m.Realtime.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting realtime module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of graphql module")
		m.Graphql.SetConfig(p.ID)

		logrus.Debugln("Setting config of lets encrypt module")
		if err := le.SetProjectDomains(p.ID, p.Modules.LetsEncrypt); err != nil {
			logrus.Errorf("error setting letsencypt module config - %s", err.Error())
		}

		logrus.Debugln("Setting config of ingress routing module")
		ingressRouting.SetProjectRoutes(p.ID, p.Modules.Routes)
	}
}

// SetGlobalConfig sets the auth secret and AESKey
func (m *Modules) SetGlobalConfig(projectID, secret, aesKey string) {
	m.Auth.SetSecret(secret)
	m.Auth.SetAESKey(aesKey)
}

// SetCrudConfig sets the config of crud, auth, schema and realtime modules
func (m *Modules) SetCrudConfig(projectID string, crudConfig config.Crud) error {
	logrus.Debugln("Setting config of crud module")
	if err := m.Crud.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting crud module config - %s", err.Error())
		return err
	}
	logrus.Debugln("Setting config of auth module")
	m.Auth.SetCrudConfig(projectID, crudConfig)

	logrus.Debugln("Setting config of schema module")
	if err := m.Schema.SetConfig(crudConfig, projectID); err != nil {
		logrus.Errorf("error setting schema module config - %s", err.Error())
		return err
	}
	if err := m.Realtime.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting realtime module config - %s", err.Error())
		return err
	}
	logrus.Debugln("Setting config of file storage module")
	return nil
}

// SetServicesConfig sets the config of auth and functions modules
func (m *Modules) SetServicesConfig(projectID string, services *config.ServicesModule) error {
	logrus.Debugln("Setting config of auth module")
	m.Auth.SetServicesConfig(projectID, services)

	logrus.Debugln("Setting config of remote services module")
	m.Functions.SetConfig(projectID, services)
	return nil
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(projectID string, fileStore *config.FileStore) error {
	logrus.Debugln("Setting config of auth module")
	m.Auth.SetFileStoreConfig(projectID, fileStore)

	logrus.Debugln("Setting config of file storage module")
	if err := m.File.SetConfig(fileStore); err != nil {
		logrus.Errorf("error setting filestore module config - %s", err.Error())
		return err
	}
	return nil
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	logrus.Debugln("Setting config of eventing module")
	if err := m.Eventing.SetConfig(projectID, eventingConfig); err != nil {
		logrus.Errorf("error setting eventing module config - %s", err.Error())
		return err
	}
	return nil
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(projectID string, auth config.Auth) {
	logrus.Debugln("Setting config of user management module")
	m.User.SetConfig(auth)
}
