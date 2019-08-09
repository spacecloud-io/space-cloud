package projects

import (
	"log"

	"github.com/spaceuptech/space-cloud/config"
)

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

	// Set the configuration for the file storage module
	if err := state.FileStore.SetConfig(config.Modules.FileStore); err != nil {
		log.Println("File storage module config error:", err)
	}

	// Set the configuration for the functions module
	if err := state.Functions.SetConfig(config.Modules.Functions); err != nil {
		log.Println("Functions module config error:", err)
	}

	// Set the configuration for the Realtime module
	if err := state.Realtime.SetConfig(config.ID, config.Modules.Realtime); err != nil {
		log.Println("Realtime module config error:", err)
	}

	// Set the configuration for the crud module
	state.Crud.SetConfig(config.Modules.Crud)

	return nil
}
