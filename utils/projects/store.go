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

		// Create a new project
		state, err = p.NewProject(config.ID)
		if err != nil {
			return err
		}
	}

	// Set the configuration for the auth module
	if err := state.Auth.SetConfig(config.ID, config.Secret, config.Modules.Crud,
		config.Modules.FileStore, config.Modules.Functions, config.Modules.Pubsub); err != nil {
		log.Println("Auth module config error:", err)
	}

	if err := state.Pubsub.SetConfig(config.Modules.Pubsub); err != nil {
		log.Println("Pubsub module config error:", err)
	}

	// Set the configuration for the crud module
	if err := state.Crud.SetConfig(config.Modules.Crud); err != nil {
		log.Println("Crud module config error:", err)
	}

	if err := state.Eventing.SetConfig(config.ID, &config.Modules.Eventing); err != nil {
		log.Println("Eventing module config error:", err)
	}

	// Set the configuration for the user management module
	state.UserManagement.SetConfig(config.Modules.Auth)

	// Set the configuration for the file storage module
	if err := state.FileStore.SetConfig(config.Modules.FileStore); err != nil {
		log.Println("File storage module config error:", err)
	}

	// Set the configuration for the Realtime module
	state.Realtime.SetConfig(config.ID, config.Modules.Crud)

	state.Graph.SetConfig(config.ID)

	return nil
}
