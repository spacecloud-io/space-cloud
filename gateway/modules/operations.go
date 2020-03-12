package modules

import (
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(config *config.Config) error {
	if config.Projects != nil {
		p := config.Projects[0]

		if err := m.Crud.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting crud module config - %s", err.Error())
			return err
		}

		if err := m.Schema.SetConfig(p.Modules.Crud, p.ID); err != nil {
			logrus.Errorf("error setting schema module config - %s", err.Error())
			return err
		}

		if err := m.Auth.SetConfig(p.ID, p.Secret, p.AESkey, p.Modules.Crud, p.Modules.FileStore, p.Modules.Services, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting auth module config - %s", err.Error())
			return err
		}

		m.Functions.SetConfig(p.ID, p.Modules.Services)

		m.User.SetConfig(p.Modules.Auth)

		if err := m.File.SetConfig(p.Modules.FileStore); err != nil {
			logrus.Errorf("error setting filestore module config - %s", err.Error())
			return err
		}

		if err := m.Eventing.SetConfig(p.ID, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting eventing module config - %s", err.Error())
			return err
		}

		if err := m.Realtime.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting realtime module config - %s", err.Error())
			return err
		}

		m.Graphql.SetConfig(p.ID)
	}
	return nil
}

// SetGlobalConfig sets the auth secret and AESkey
func (m *Modules) SetGlobalConfig(projectID, secret, aesKey string) {
	m.Auth.SetSecret(secret)
	m.Auth.SetAESKey(aesKey)
}

// SetCrudConfig sets the config of crud, auth, schema and realtime modules
func (m *Modules) SetCrudConfig(projectID, secret, aesKey string, crudConfig config.Crud, fileStore *config.FileStore, services *config.ServicesModule, eventing *config.Eventing) error {
	if err := m.Crud.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting crud module config - %s", err.Error())
		return err
	}
	if err := m.Auth.SetConfig(projectID, secret, aesKey, crudConfig, fileStore, services, eventing); err != nil {
		logrus.Errorf("error setting auth module config - %s", err.Error())
		return err
	}
	if err := m.Schema.SetConfig(crudConfig, projectID); err != nil {
		logrus.Errorf("error setting schema module config - %s", err.Error())
		return err
	}
	if err := m.Realtime.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting realtime module config - %s", err.Error())
		return err
	}
	return nil
}

// SetServicesConfig sets the config of auth and functions modules
func (m *Modules) SetServicesConfig(projectID, secret, aesKey string, crudConfig config.Crud, fileStore *config.FileStore, services *config.ServicesModule, eventing *config.Eventing) error {
	if err := m.Auth.SetConfig(projectID, secret, aesKey, crudConfig, fileStore, services, eventing); err != nil {
		logrus.Errorf("error setting auth module config - %s", err.Error())
		return err
	}
	m.Functions.SetConfig(projectID, services)
	return nil
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(projectID, secret, aesKey string, crudConfig config.Crud, fileStore *config.FileStore, services *config.ServicesModule, eventing *config.Eventing) error {
	if err := m.Auth.SetConfig(projectID, secret, aesKey, crudConfig, fileStore, services, eventing); err != nil {
		logrus.Errorf("error setting auth module config - %s", err.Error())
		return err
	}
	if err := m.File.SetConfig(fileStore); err != nil {
		logrus.Errorf("error setting filestore module config - %s", err.Error())
		return err
	}
	return nil
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	if err := m.Eventing.SetConfig(projectID, eventingConfig); err != nil {
		logrus.Errorf("error setting eventing module config - %s", err.Error())
		return err
	}
	return nil
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(projectID string, auth config.Auth) {
	m.User.SetConfig(auth)
}
