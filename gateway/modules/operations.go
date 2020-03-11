package modules

import (
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(config *config.Config) error {
	if config.Projects != nil {
		p := config.Projects[0]

		if err := m.crud.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting crud module config - %s", err.Error())
			return err
		}

		if err := m.schema.SetConfig(p.Modules.Crud, p.ID); err != nil {
			logrus.Errorf("error setting schema module config - %s", err.Error())
			return err
		}

		if err := m.auth.SetConfig(p.ID, p.Secret, p.AESkey, p.Modules.Crud, p.Modules.FileStore, p.Modules.Services, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting auth module config - %s", err.Error())
			return err
		}

		m.functions.SetConfig(p.ID, p.Modules.Services)

		m.user.SetConfig(p.Modules.Auth)

		if err := m.file.SetConfig(p.Modules.FileStore); err != nil {
			logrus.Errorf("error setting filestore module config - %s", err.Error())
			return err
		}

		if err := m.eventing.SetConfig(p.ID, &p.Modules.Eventing); err != nil {
			logrus.Errorf("error setting eventing module config - %s", err.Error())
			return err
		}

		if err := m.realtime.SetConfig(p.ID, p.Modules.Crud); err != nil {
			logrus.Errorf("error setting realtime module config - %s", err.Error())
			return err
		}

		m.graphql.SetConfig(p.ID)
	}
	return nil
}

// SetGlobalConfig sets the auth secret and AESkey
func (m *Modules) SetGlobalConfig(projectID, secret, aesKey string) {
	m.auth.SetSecret(secret)
	m.auth.SetAESKey(aesKey)
}

// SetCrudConfig sets the config of crud, auth, schema and realtime modules
func (m *Modules) SetCrudConfig(projectID, secret, aesKey string, crudConfig config.Crud, fileStore *config.FileStore, services *config.ServicesModule, eventing *config.Eventing) error {
	if err := m.crud.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting crud module config - %s", err.Error())
		return err
	}
	if err := m.auth.SetConfig(projectID, secret, aesKey, crudConfig, fileStore, services, eventing); err != nil {
		logrus.Errorf("error setting auth module config - %s", err.Error())
		return err
	}
	if err := m.schema.SetConfig(crudConfig, projectID); err != nil {
		logrus.Errorf("error setting schema module config - %s", err.Error())
		return err
	}
	if err := m.realtime.SetConfig(projectID, crudConfig); err != nil {
		logrus.Errorf("error setting realtime module config - %s", err.Error())
		return err
	}
	return nil
}

// SetServicesConfig sets the config of auth and functions modules
func (m *Modules) SetServicesConfig(projectID, secret, aesKey string, crudConfig config.Crud, fileStore *config.FileStore, services *config.ServicesModule, eventing *config.Eventing) error {
	if err := m.auth.SetConfig(projectID, secret, aesKey, crudConfig, fileStore, services, eventing); err != nil {
		logrus.Errorf("error setting auth module config - %s", err.Error())
		return err
	}
	m.functions.SetConfig(projectID, services)
	return nil
}

// SetFileStoreConfig sets the config of auth and filestore modules
func (m *Modules) SetFileStoreConfig(projectID, secret, aesKey string, crudConfig config.Crud, fileStore *config.FileStore, services *config.ServicesModule, eventing *config.Eventing) error {
	if err := m.auth.SetConfig(projectID, secret, aesKey, crudConfig, fileStore, services, eventing); err != nil {
		logrus.Errorf("error setting auth module config - %s", err.Error())
		return err
	}
	if err := m.file.SetConfig(fileStore); err != nil {
		logrus.Errorf("error setting filestore module config - %s", err.Error())
		return err
	}
	return nil
}

// SetEventingConfig sets the config of eventing module
func (m *Modules) SetEventingConfig(projectID string, eventingConfig *config.Eventing) error {
	if err := m.eventing.SetConfig(projectID, eventingConfig); err != nil {
		logrus.Errorf("error setting eventing module config - %s", err.Error())
		return err
	}
	return nil
}

// SetUsermanConfig set the config of the userman module
func (m *Modules) SetUsermanConfig(projectID string, auth config.Auth) {
	m.user.SetConfig(auth)
}
