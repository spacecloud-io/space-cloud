package migrate

import (
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func getAdminProjectConfig(resource *model.SCConfig, configPath string) error {
	projectConfigs, err := utils.ReadSpecObjectsFromFile(filepath.Join(configPath, "01-projects.yaml"))
	if err != nil {
		return err
	}

	module, ok := resource.Config["adminman"]
	if !ok {
		module = make(model.ConfigModule)
	}

	for _, projectConfig := range projectConfigs {
		value := new(config.ProjectConfig)
		if err := mapstructure.Decode(projectConfig.Spec, value); err != nil {
			return err
		}

		projectID := projectConfig.Meta["project"]
		adminProjectConfig := config.AdminProjectConfig{
			ID:                 value.ID,
			Name:               value.Name,
			DockerRegistry:     value.DockerRegistry,
			ContextTimeGraphQL: value.ContextTimeGraphQL,
		}

		res := &model.ResourceObject{
			Meta: model.ResourceMeta{
				Module:  "adminman",
				Type:    "project",
				Name:    projectID,
				Parents: map[string]string{},
			},
			Spec: adminProjectConfig,
		}

		projectResourceObjects, ok := module["project"]
		if !ok {
			projectResourceObjects = make([]*model.ResourceObject, 0)
		}
		projectResourceObjects = append(projectResourceObjects, res)
		module["project"] = projectResourceObjects

		adminProjectAesKey := config.AdminProjectAesKey{
			Key: value.AESKey,
		}

		res = &model.ResourceObject{
			Meta: model.ResourceMeta{
				Module: "adminman",
				Type:   "aes-key",
				Name:   projectID,
				Parents: map[string]string{
					"project": projectID,
				},
			},
			Spec: adminProjectAesKey,
		}

		aesResourceObjects, ok := module["aes-key"]

		if !ok {
			aesResourceObjects = make([]*model.ResourceObject, 0)
		}
		aesResourceObjects = append(aesResourceObjects, res)
		module["aes-key"] = aesResourceObjects

		resourceObjects, ok := module["jwt-secret"]
		if !ok {
			resourceObjects = make([]*model.ResourceObject, 0)
		}
		for _, secret := range value.Secrets {
			res = &model.ResourceObject{
				Meta: model.ResourceMeta{
					Module: "adminman",
					Type:   "jwt-secret",
					Name:   projectID,
					Parents: map[string]string{
						"project": projectID,
					},
				},
				Spec: secret,
			}
			resourceObjects = append(resourceObjects, res)
		}
		module["jwt-secret"] = resourceObjects
		resource.Config["adminman"] = module
	}

	return nil
}
