package common

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spf13/viper"
)

func prepareAdminApp(fileConfig *model.SCConfig) json.RawMessage {
	user := viper.GetString("admin-user")
	pass := viper.GetString("admin-pass")
	secret := viper.GetString("admin-secret")
	isDev := viper.GetBool("dev")

	config := map[string]interface{}{
		"user":     user,
		"pass":     pass,
		"secret":   secret,
		"isDev":    isDev,
		"projects": prepareAdminProjects(fileConfig),
	}

	data, _ := json.Marshal(config)
	return data
}

func prepareAdminProjects(fileConfig *model.SCConfig) map[string]*config.AdminProject {
	adminProjects := make(map[string]*config.AdminProject)
	module, ok := fileConfig.Config["adminman"]
	if !ok {
		return adminProjects
	}
	resourceObjects, ok := module["config"]
	if !ok {
		return adminProjects
	}

	for _, resourceObject := range resourceObjects {
		adminConfig := new(config.AdminProjectConfig)
		if err := mapstructure.Decode(resourceObject.Spec, adminConfig); err != nil {
			return map[string]*config.AdminProject{}
		}

		projectID := resourceObject.Meta.Parents["project"]
		adminProjects[projectID] = &config.AdminProject{
			Config:  *adminConfig,
			Secrets: make([]*config.Secret, 0),
		}
	}

	resourceObjects, ok = module["aes"]
	if ok {
		for _, resourceObject := range resourceObjects {
			aes := new(config.AdminProjectAesKey)
			if err := mapstructure.Decode(resourceObject.Spec, aes); err != nil {
				return map[string]*config.AdminProject{}
			}

			projectID := resourceObject.Meta.Parents["project"]
			adminProject := adminProjects[projectID]
			adminProject.AesKey = *aes
			adminProjects[projectID] = adminProject
		}
	}

	resourceObjects, ok = module["secrets"]
	if ok {
		for _, resourceObject := range resourceObjects {
			schema := new(config.Secret)
			if err := mapstructure.Decode(resourceObject.Spec, schema); err != nil {
				return map[string]*config.AdminProject{}
			}

			projectID := resourceObject.Meta.Parents["project"]
			adminProject := adminProjects[projectID]
			adminProject.Secrets = append(adminProject.Secrets, schema)
			adminProjects[projectID] = adminProject
		}
	}

	return adminProjects
}
