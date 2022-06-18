package remoteservice

import (
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func getServiceConfigTypes() model.ConfigTypes {
	reflector := utils.GetJSONSchemaReflector()
	return model.ConfigTypes{
		"config": &model.ConfigTypeDefinition{
			Schema:          reflector.Reflect(&config.Service{}),
			RequiredParents: []string{"project"},
		},
	}
}
