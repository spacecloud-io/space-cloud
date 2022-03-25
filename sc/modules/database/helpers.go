package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/modules/database/connectors"
	"github.com/spacecloud-io/space-cloud/modules/database/connectors/schema"
)

// CombineDBConfigKey joing project id and db alias to form the db config key
func CombineDBConfigKey(projectID, dbAlias string) string {
	return fmt.Sprintf("%s---%s", projectID, dbAlias)
}

// SplitDBConfigKey split the db config key into project id and dbAlias
func SplitDBConfigKey(key string) (project, dbAlias string) {
	arr := strings.Split(key, "---")
	return arr[0], arr[1]
}

func generateUniqueDBKey(projectID string, c *config.DatabaseConfig) string {
	return fmt.Sprintf("%s---%s--%v", CombineDBConfigKey(projectID, c.DbAlias), c.DBName, c.DriverConf)
}

func (l *App) processDBSchemaHook(ctx context.Context, obj *configman.ResourceObject) error {
	// Convert object to a known type
	dbSchema := new(config.DatabaseSchema)
	if err := mapstructure.Decode(obj.Spec, dbSchema); err != nil {
		return err
	}

	// Check if database exists
	db, p := l.connectors[obj.Meta.Parents["db"]]
	if !p {
		return fmt.Errorf("unknown database alias '%s' provided", obj.Meta.Parents["db"])
	}

	// Set some spec values which may be absent
	m := obj.Spec.(map[string]interface{})
	m["col"] = obj.Meta.Name
	m["dbAlias"] = obj.Meta.Parents["db"]

	// Try to create the table in the database
	newSchema, err := schema.Parser(config.DatabaseSchemas{obj.Meta.Name: dbSchema})
	if err != nil {
		return err
	}
	return db.ApplyCollectionSchema(ctx, obj.Meta.Name, newSchema)
}

func processPreparedQuery(obj *configman.ResourceObject) error {
	// Set some spec values which may be absent
	m := obj.Spec.(map[string]interface{})
	m["id"] = obj.Meta.Name
	m["dbAlias"] = obj.Meta.Parents["db"]

	return nil
}

func processConfig(obj *configman.ResourceObject) error {
	// Set some spec values which may be absent
	m := obj.Spec.(map[string]interface{})
	m["dbAlias"] = obj.Meta.Name

	return nil
}

func (l *App) getConnector(project, db string) (*connectors.Module, error) {
	conn, p := l.connectors[CombineDBConfigKey(project, db)]
	if !p {
		return nil, fmt.Errorf("database '%s' does not exist in project '%s'", db, project)
	}

	return conn, nil
}
