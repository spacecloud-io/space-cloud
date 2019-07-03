package projects

import "github.com/spaceuptech/space-cloud/config"

// StoreProject stores the provided project config
func (p *Projects) StoreProject(config *config.Project) error {
	// Get the project. Create if not exists
	state, err := p.LoadProject(config.ID)
	if err != nil {
		state = p.NewProject(config.ID)
	}

	// Set the configuration for the auth module
	state.Auth.SetConfig(config.ID, config.Secret, config.Modules.Crud, config.Modules.FileStore, config.Modules.Functions)

	// Set the configuration for the user management module
	state.UserManagement.SetConfig(config.Modules.Auth)

	if err := state.Static.SetConfig(config.Modules.Static); err != nil {
		return err
	}

	// Set the configuration for the file storage module
	if err := state.FileStore.SetConfig(config.Modules.FileStore); err != nil {
		return err
	}

	// Set the configuration for the functions module
	if err := state.Functions.SetConfig(config.Modules.Functions); err != nil {
		return err
	}

	// Set the configuration for the Realtime module
	if err := state.Realtime.SetConfig(config.ID, config.Modules.Realtime); err != nil {
		return err
	}

	// Set the configuration for the crud module
	return state.Crud.SetConfig(config.Modules.Crud)
}
