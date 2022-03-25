package database

import (
	"context"
	"fmt"

	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/model"
)

// Hook implements the configman hook functionality
func (l *App) Hook(ctx context.Context, obj *configman.ResourceObject) error {
	// Check if the resource belongs to this app
	if obj.Meta.Module != "database" {
		return fmt.Errorf("hook invoked for invalid resource type '%s/%s'", obj.Meta.Module, obj.Meta.Type)
	}

	// Process hook based on the resource type
	switch obj.Meta.Type {
	case "config":
		return processConfig(obj)
	case "schema":
		return l.processDBSchemaHook(ctx, obj)
	case "prepared-query":
		return processPreparedQuery(obj)
	default:
		return fmt.Errorf("hook invoked for invalid resource type '%s/%s'", obj.Meta.Module, obj.Meta.Type)
	}
}

// GetParsedSchemas returns the parsed schemas for databases across all projects
func (l *App) GetParsedSchemas() map[string]model.DBSchemas {
	allSchemas := make(map[string]model.DBSchemas)

	// Iterate over all connectors
	for k, connector := range l.connectors {
		project, dbAlias := SplitDBConfigKey(k)

		if _, p := allSchemas[project]; !p {
			allSchemas[project] = model.DBSchemas{}
		}

		allSchemas[project][dbAlias] = connector.GetParsedSchemas()
	}

	return allSchemas
}

// Read performs a db read operation
func (l *App) Read(ctx context.Context, project, db, col string, req *model.ReadRequest, params model.RequestParams) (interface{}, *model.SQLMetaData, error) {
	conn, err := l.getConnector(project, db)
	if err != nil {
		return nil, nil, err
	}

	return conn.Read(ctx, col, req, params)
}
