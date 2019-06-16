package projects

import (
	"context"
	"encoding/json"
	"log"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// LoadConfigFromDB reads the config from the specified datbase directly
func (p *Projects) LoadConfigFromDB(account, dbType, conn string) error {
	state := p.NewProject(utils.SpaceCloudProject)
	crudConfig := map[string]*config.CrudStub{
		dbType: &config.CrudStub{
			Conn:        conn,
			Collections: map[string]*config.TableRule{},
		},
	}

	if err := state.Crud.SetConfig(crudConfig); err != nil {
		return err
	}

	if err := state.Realtime.SetConfig(&config.Realtime{Enabled: true, Broker: utils.Nats, Conn: "nats://localhost:4222"}); err != nil {
		return err
	}

	feedData, err := state.Realtime.DoRealtimeSubscribe(context.TODO(), utils.SpaceCloudProject, state.Crud, &model.RealtimeRequest{
		DBType:  dbType,
		Group:   utils.SpaceCloudConfigTable,
		Project: utils.SpaceCloudProject,
		ID:      utils.SpaceCloudProject,
		Where:   map[string]interface{}{"account": account},
	}, func(data *model.FeedData) {

		switch data.Type {
		case utils.RealtimeDelete:
			project := data.Payload["project"].(string)
			p.DeleteProject(project)

		case utils.RealtimeWrite, utils.RealtimeUpdate:
			project := data.Payload["project"].(string)
			config := data.Payload["config"].(string)

			err := p.setConfig(data.Type, project, config)
			if err != nil {
				log.Println("Projects: Error - Could not load config", err)
			}
		}
	})
	if err != nil {
		return err
	}

	for _, data := range feedData {
		project := data.Payload["project"].(string)
		config := data.Payload["config"].(string)

		err := p.setConfig(data.Type, project, config)
		if err != nil {
			log.Println("Projects: Error - Could not load config", err)
		}
	}
	return nil
}

func (p *Projects) setConfig(action, project string, data string) error {

	// Delete the project if the action was delete
	if action == utils.RealtimeDelete {
		p.DeleteProject(project)
		return nil
	}

	// Get the project. Create if not exists
	state, err := p.LoadProject(project)
	if err != nil {
		state = p.NewProject(project)
	}

	// Parse the config string to a type config.Project
	config := new(config.Project)
	err = json.Unmarshal([]byte(data), config)
	if err != nil {
		return err
	}

	// Set the configuration for the auth module
	state.Auth.SetConfig(config.ID, config.Secret, config.Modules.Crud, config.Modules.FileStore, config.Modules.Functions)

	// Set the configuration for the user management module
	state.UserManagement.SetConfig(config.Modules.Auth)

	// Set the configuration for the file storage module
	if err := state.FileStore.SetConfig(config.Modules.FileStore); err != nil {
		return err
	}

	// Set the configuration for the functions module
	if err := state.Functions.SetConfig(config.Modules.Functions); err != nil {
		return err
	}

	// Set the configuration for the Realtime module
	if err := state.Realtime.SetConfig(config.Modules.Realtime); err != nil {
		return err
	}

	// Set the configuration for the crud module
	return state.Crud.SetConfig(config.Modules.Crud)
}
