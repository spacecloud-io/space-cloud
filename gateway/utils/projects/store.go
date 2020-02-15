package projects

import (
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// StoreIgnoreErrors stores a project config while silently ignoring the errors
func (p *Projects) StoreIgnoreErrors(project *config.Project) error {
	// Get the project. Create if not exists
	s, err := p.LoadProject(project.ID)
	if err != nil {

		// Create a new project
		s, err = p.NewProject(project.ID)
		if err != nil {
			return err
		}
	}

	// Always set the config of the crud module first
	// Set the configuration for the crud module
	if err := s.Crud.SetConfig(project.ID, project.Modules.Crud); err != nil {
		logrus.Errorln("Error in crud module config:", err)
	}

	if err := s.Schema.SetConfig(project.Modules.Crud, project.ID); err != nil {
		logrus.Errorln("Error in schema module config:", err)
	}

	// Set the configuration for the auth module
	if err := s.Auth.SetConfig(project.ID, project.Secret, project.Modules.Crud, project.Modules.FileStore, project.Modules.Services, &project.Modules.Eventing); err != nil {
		logrus.Errorln("Error in auth module config:", err)
	}

	// Set the configuration for the functions module
	s.Functions.SetConfig(project.ID, project.Modules.Services)

	// Set the configuration for the user management module
	s.UserManagement.SetConfig(project.Modules.Auth)

	// Set the configuration for the file storage module
	if err := s.FileStore.SetConfig(project.Modules.FileStore); err != nil {
		logrus.Errorln("Error in files module config:", err)
	}

	if err := s.Eventing.SetConfig(project.ID, &project.Modules.Eventing); err != nil {
		logrus.Errorln("Error in eventing module config:", err)
	}

	// Set the configuration for the realtime module
	if err := s.Realtime.SetConfig(project.ID, project.Modules.Crud); err != nil {
		logrus.Errorln("Error in realtime module config:", err)
	}

	// Set the configuration for the graphql module
	s.Graph.SetConfig(project.ID)

	// Set the configuration for the letsencrypt module
	if err := p.letsencrypt.SetProjectDomains(project.ID, project.Modules.LetsEncrypt); err != nil {
		logrus.Errorln("Error in letsencrypt module config:", err)
		return err
	}

	// Set the configuration for the routing module
	p.routing.SetProjectRoutes(project.ID, project.Modules.Routes)

	return nil
}

// StoreProject stores the provided project config
func (p *Projects) StoreProject(project *config.Project) error {
	// Get the project. Create if not exists
	s, err := p.LoadProject(project.ID)
	if err != nil {

		// Create a new project
		s, err = p.NewProject(project.ID)
		if err != nil {
			return err
		}
	}

	// Always set the config of the crud module first
	// Set the configuration for the crud module
	if err := s.Crud.SetConfig(project.ID, project.Modules.Crud); err != nil {
		logrus.Errorln("Error in crud module config:", err)
		return err
	}

	if err := s.Schema.SetConfig(project.Modules.Crud, project.ID); err != nil {
		logrus.Errorln("Error in schema module config:", err)
		return err
	}

	// Set the configuration for the auth module
	if err := s.Auth.SetConfig(project.ID, project.Secret, project.Modules.Crud, project.Modules.FileStore, project.Modules.Services, &project.Modules.Eventing); err != nil {
		logrus.Errorln("Error in auth module config:", err)
		return err
	}

	// Set the configuration for the functions module
	s.Functions.SetConfig(project.ID, project.Modules.Services)

	// Set the configuration for the user management module
	s.UserManagement.SetConfig(project.Modules.Auth)

	// Set the configuration for the file storage module
	if err := s.FileStore.SetConfig(project.Modules.FileStore); err != nil {
		logrus.Errorln("Error in files module config:", err)
		return err
	}

	if err := s.Eventing.SetConfig(project.ID, &project.Modules.Eventing); err != nil {
		logrus.Errorln("Error in eventing module config:", err)
		return err
	}

	// Set the configuration for the realtime module
	if err := s.Realtime.SetConfig(project.ID, project.Modules.Crud); err != nil {
		logrus.Errorln("Error in realtime module config:", err)
		return err
	}

	// Set the configuration for the graphql module
	s.Graph.SetConfig(project.ID)

	// Set the configuration for the letsencrypt module
	if err := p.letsencrypt.SetProjectDomains(project.ID, project.Modules.LetsEncrypt); err != nil {
		logrus.Errorln("Error in letsencrypt module config:", err)
		return err
	}

	// Set the configuration for the routing module
	p.routing.SetProjectRoutes(project.ID, project.Modules.Routes)

	return nil
}

// SetGlobalConfig stores the global configuration of a project
func (p *Projects) SetGlobalConfig(projectID, secret string) error {
	s, err := p.LoadProject(projectID)
	if err != nil {
		return err
	}

	s.Auth.SetSecret(secret)
	return nil
}

// SetCrudConfig sets the crud config
func (p *Projects) SetCrudConfig(projectID string, c config.Crud) error {
	s, err := p.LoadProject(projectID)
	if err != nil {
		return err
	}

	if err := s.Crud.SetConfig(projectID, c); err != nil {
		logrus.Errorln("Error in crud module config:", err)
		return err
	}

	// Set the configuration for the auth module
	if err := s.Auth.SetCrudConfig(projectID, c); err != nil {
		logrus.Errorln("Error in crud module config:", err)
		return err
	}

	if err := s.Schema.SetConfig(c, projectID); err != nil {
		logrus.Errorln("Error in schema module config:", err)
		return err
	}

	// Set the configuration for the realtime module
	if err := s.Realtime.SetConfig(projectID, c); err != nil {
		logrus.Errorln("Error in realtime module config:", err)
		return err
	}

	return nil
}

// SetServicesConfig sets the config for the remote service module
func (p *Projects) SetServicesConfig(projectID string, c *config.ServicesModule) error {
	s, err := p.LoadProject(projectID)
	if err != nil {
		return err
	}

	s.Auth.SetServicesConfig(projectID, c)
	s.Functions.SetConfig(projectID, c)
	return nil
}

// SetFileStoreConfig sets the config for the file storage module
func (p *Projects) SetFileStoreConfig(projectID string, c *config.FileStore) error {
	s, err := p.LoadProject(projectID)
	if err != nil {
		return err
	}

	s.Auth.SetFileStoreConfig(projectID, c)

	// Set the configuration for the file storage module
	if err := s.FileStore.SetConfig(c); err != nil {
		logrus.Errorln("Error in files module config:", err)
		return err
	}

	return nil
}

// SetEventingConfig sets the config for the eventing module
func (p *Projects) SetEventingConfig(projectID string, eventing *config.Eventing) error {
	s, err := p.LoadProject(projectID)
	if err != nil {
		return err
	}

	if err := s.Eventing.SetConfig(projectID, eventing); err != nil {
		logrus.Errorln("Error in eventing module config:", err)
		return err
	}

	return nil
}

// SetUserManConfig sets the config for the user management module
func (p *Projects) SetUserManConfig(projectID string, userMan config.Auth) error {
	s, err := p.LoadProject(projectID)
	if err != nil {
		return err
	}

	// Set the configuration for the user management module
	s.UserManagement.SetConfig(userMan)
	return nil
}

// SetLetsencryptDomains sets the lets encrypt domains in the given project
func (p *Projects) SetLetsencryptDomains(projectID string, c config.LetsEncrypt) error {
	return p.letsencrypt.SetProjectDomains(projectID, c)
}

// RemoveLetsencryptDomains removes the domains for the given project
func (p *Projects) RemoveLetsencryptDomains(projectID string) error {
	return p.letsencrypt.DeleteProjectDomains(projectID)
}

// SetProjectRoutes sets the routes of a project
func (p *Projects) SetProjectRoutes(projectID string, routes config.Routes) {
	p.routing.SetProjectRoutes(projectID, routes)
}

// RemoveProjectRoutes removes the routes of a project
func (p *Projects) RemoveProjectRoutes(projectID string) {
	p.routing.DeleteProjectRoutes(projectID)
}
