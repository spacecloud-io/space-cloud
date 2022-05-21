package migrate

import (
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func getDBConfig(resource *model.SCConfig, configPath string) error {
	databaseConfigs, err := utils.ReadSpecObjectsFromFile(filepath.Join(configPath, "03-db-configs.yaml"))
	if err != nil {
		return err
	}

	for _, databaseConfig := range databaseConfigs {
		value := new(config.DatabaseConfig)
		if err := mapstructure.Decode(databaseConfig.Spec, value); err != nil {
			return err
		}

		projectID := databaseConfig.Meta["project"]
		dbAlias := databaseConfig.Meta["dbAlias"]

		value.DbAlias = dbAlias
		res := model.ResourceObject{
			Meta: model.ResourceMeta{
				Module: "project",
				Type:   "config",
				Name:   dbAlias,
				Parents: map[string]string{
					"project": projectID,
				},
			},
			Spec: value,
		}

		module, ok := resource.Config["database"]
		if !ok {
			module = make(model.ConfigModule)
		}
		resourceObjects, ok := module["config"]
		if !ok {
			resourceObjects = make([]*model.ResourceObject, 0)
		}

		resourceObjects = append(resourceObjects, &res)
		module["config"] = resourceObjects
		resource.Config["database"] = module
	}
	return nil
}

func getDBRule(resource *model.SCConfig, configPath string) error {
	databaseRules, err := utils.ReadSpecObjectsFromFile(filepath.Join(configPath, "04-db-rules.yaml"))
	if err != nil {
		return err
	}

	for _, databaseRule := range databaseRules {
		value := new(config.DatabaseRule)
		if err := mapstructure.Decode(databaseRule.Spec, value); err != nil {
			return err
		}

		projectID := databaseRule.Meta["project"]
		dbAlias := databaseRule.Meta["dbAlias"]
		table := databaseRule.Meta["col"]
		value.Table = table
		value.DbAlias = dbAlias

		value.DbAlias = dbAlias
		res := model.ResourceObject{
			Meta: model.ResourceMeta{
				Module: "project",
				Type:   "rules",
				Name:   table,
				Parents: map[string]string{
					"project":  projectID,
					"database": dbAlias,
				},
			},
			Spec: value,
		}

		module, ok := resource.Config["database"]
		if !ok {
			module = make(model.ConfigModule)
		}
		resourceObjects, ok := module["rules"]
		if !ok {
			resourceObjects = make([]*model.ResourceObject, 0)
		}

		resourceObjects = append(resourceObjects, &res)
		module["rules"] = resourceObjects
		resource.Config["database"] = module
	}
	return nil
}

func getDBSchema(resource *model.SCConfig, configPath string) error {
	databaseSchemas, err := utils.ReadSpecObjectsFromFile(filepath.Join(configPath, "05-db-schemas.yaml"))
	if err != nil {
		return err
	}

	for _, databaseSchema := range databaseSchemas {
		value := new(config.DatabaseSchema)
		if err := mapstructure.Decode(databaseSchema.Spec, value); err != nil {
			return err
		}

		projectID := databaseSchema.Meta["project"]
		dbAlias := databaseSchema.Meta["dbAlias"]
		table := databaseSchema.Meta["col"]
		value.Table = table
		value.DbAlias = dbAlias

		value.DbAlias = dbAlias
		res := model.ResourceObject{
			Meta: model.ResourceMeta{
				Module: "project",
				Type:   "schema",
				Name:   table,
				Parents: map[string]string{
					"project":  projectID,
					"database": dbAlias,
				},
			},
			Spec: value,
		}

		module, ok := resource.Config["database"]
		if !ok {
			module = make(model.ConfigModule)
		}
		resourceObjects, ok := module["schema"]
		if !ok {
			resourceObjects = make([]*model.ResourceObject, 0)
		}

		resourceObjects = append(resourceObjects, &res)
		module["schema"] = resourceObjects
		resource.Config["database"] = module
	}
	return nil
}

func getDBPreparedQuery(resource *model.SCConfig, configPath string) error {
	databasePreparedQueries, err := utils.ReadSpecObjectsFromFile(filepath.Join(configPath, "06-db-prepared-query.yaml"))
	if err != nil {
		return err
	}

	for _, databasedatabasePreparedQuery := range databasePreparedQueries {
		value := new(config.DatbasePreparedQuery)
		if err := mapstructure.Decode(databasedatabasePreparedQuery.Spec, value); err != nil {
			return err
		}

		projectID := databasedatabasePreparedQuery.Meta["project"]
		dbAlias := databasedatabasePreparedQuery.Meta["db"]
		id := databasedatabasePreparedQuery.Meta["id"]
		value.ID = id
		value.DbAlias = dbAlias

		value.DbAlias = dbAlias
		res := model.ResourceObject{
			Meta: model.ResourceMeta{
				Module: "project",
				Type:   "prepared-query",
				Name:   id,
				Parents: map[string]string{
					"project":  projectID,
					"database": dbAlias,
				},
			},
			Spec: value,
		}

		module, ok := resource.Config["database"]
		if !ok {
			module = make(model.ConfigModule)
		}
		resourceObjects, ok := module["prepared-query"]
		if !ok {
			resourceObjects = make([]*model.ResourceObject, 0)
		}

		resourceObjects = append(resourceObjects, &res)
		module["prepared-query"] = resourceObjects
		resource.Config["database"] = module
	}
	return nil
}
