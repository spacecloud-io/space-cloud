package database

import (
	"context"

	"github.com/spacecloud-io/space-cloud/model"
)

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

// Batch performs a db transaction
func (l *App) Batch(ctx context.Context, project, db, col string, req *model.BatchRequest, params model.RequestParams) error {
	conn, err := l.getConnector(project, db)
	if err != nil {
		return err
	}

	return conn.Batch(ctx, req, params)
}
