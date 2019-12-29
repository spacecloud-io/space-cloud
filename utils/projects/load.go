package projects

// // LoadConfigFromDB reads the config from the specified datbase directly
// func (p *Projects) LoadConfigFromDB(account, dbType, conn string) error {
// 	state, err := p.NewProject(utils.SpaceCloudProject)
// 	if err != nil {
// 		return err
// 	}
//
// 	crudConfig := map[string]*config.CrudStub{
// 		dbType: {
// 			Conn:        conn,
// 			Collections: map[string]*config.TableRule{},
// 		},
// 	}
//
// 	state.Crud.SetConfig(crudConfig)
//
// 	if err := state.Realtime.SetConfig(utils.SpaceCloudProject); err != nil {
// 		return err
// 	}
//
// 	feedData, err := state.Realtime.DoRealtimeSubscribe(context.Background(), utils.SpaceCloudProject, state.Crud, &model.RealtimeRequest{
// 		DBType:  dbType,
// 		Group:   utils.SpaceCloudConfigTable,
// 		Project: utils.SpaceCloudProject,
// 		ID:      utils.SpaceCloudProject,
// 		Where:   map[string]interface{}{"account": account},
// 	}, func(data *model.FeedData) {
//
// 		switch data.Type {
// 		case utils.RealtimeDelete:
// 			project := data.Payload["project"].(string)
// 			p.DeleteProject(project)
//
// 		case utils.RealtimeInsert, utils.RealtimeUpdate:
// 			project := data.Payload["project"].(string)
// 			config := data.Payload["config"].(string)
//
// 			err := p.setConfig(data.Type, project, config)
// 			if err != nil {
// 				log.Println("Projects: Error - Could not load config", err)
// 			}
// 		}
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	for _, data := range feedData {
// 		project := data.Payload["project"].(string)
// 		config := data.Payload["config"].(string)
//
// 		err := p.setConfig(data.Type, project, config)
// 		if err != nil {
// 			log.Println("Projects: Error - Could not load config", err)
// 		}
// 	}
// 	return nil
// }
//
// func (p *Projects) setConfig(action, project string, data string) error {
//
// 	// Delete the project if the action was delete
// 	if action == utils.RealtimeDelete {
// 		p.DeleteProject(project)
// 		return nil
// 	}
//
// 	// Parse the config string to a type config.Project
// 	config := new(config.Project)
// 	err := json.Unmarshal([]byte(data), config)
// 	if err != nil {
// 		return err
// 	}
//
// 	return p.StoreProject(config)
// }
