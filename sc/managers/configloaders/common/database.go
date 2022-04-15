package common

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/modules/database"
)

func prepareDatabaseApp(fileConfig *model.SCConfig) []byte {
	dbConfigs := map[string]*database.Config{}

	module, ok := fileConfig.Config["database"]
	if !ok {
		return prepareReturn(dbConfigs)
	}
	resourceObjects, ok := module["config"]
	if !ok {
		return prepareReturn(dbConfigs)
	}

	for _, resourceObject := range resourceObjects {
		dbConfig := new(config.DatabaseConfig)
		if err := mapstructure.Decode(resourceObject.Spec, dbConfig); err != nil {
			return prepareReturn(map[string]*database.Config{})
		}

		projectID := resourceObject.Meta.Parents["project"]
		key := database.CombineDBConfigKey(projectID, dbConfig.DbAlias)
		dbConfigs[key] = &database.Config{
			Connector:       dbConfig,
			Schemas:         make(config.DatabaseSchemas),
			PreparedQueries: make(config.DatabasePreparedQueries),
		}
	}

	resourceObjects, ok = module["schema"]
	if ok {
		for _, resourceObject := range resourceObjects {
			schema := new(config.DatabaseSchema)
			if err := mapstructure.Decode(resourceObject.Spec, schema); err != nil {
				return prepareReturn(map[string]*database.Config{})
			}

			projectID := resourceObject.Meta.Parents["project"]
			key := database.CombineDBConfigKey(projectID, schema.DbAlias)
			dbConfigs[key].Schemas[schema.Table] = schema
		}
	}

	resourceObjects, ok = module["prepared-query"]
	if ok {
		for _, resourceObject := range resourceObjects {
			query := new(config.DatbasePreparedQuery)
			if err := mapstructure.Decode(resourceObject.Spec, query); err != nil {
				return prepareReturn(map[string]*database.Config{})
			}

			projectID := resourceObject.Meta.Parents["project"]
			key := database.CombineDBConfigKey(projectID, query.DbAlias)
			dbConfigs[key].PreparedQueries[query.ID] = query
		}
	}
	return prepareReturn(dbConfigs)
}

func prepareReturn(dbConfigs map[string]*database.Config) []byte {
	data, _ := json.Marshal(map[string]interface{}{"dbConfigs": dbConfigs})
	return data
}
