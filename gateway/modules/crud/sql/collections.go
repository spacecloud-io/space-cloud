package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v8"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetCollections returns collection / tables name of specified database
func (s *SQL) GetCollections(ctx context.Context) ([]utils.DatabaseCollections, error) {
	dialect := goqu.Dialect(s.dbType)
	query := dialect.From("information_schema.tables").Prepared(true).Select("table_name").Where(goqu.Ex{"table_schema": s.name})

	sqlString, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}
	if s.dbType == "sqlserver" {
		new := strings.Replace(sqlString, "?", "@p1", -1)
		sqlString = new
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)
	rows, err := s.getClient().QueryxContext(ctx, sqlString, args...)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to query database to get tables in database (%s)", s.name), err, nil)
	}
	defer func() { _ = rows.Close() }()

	result := make([]utils.DatabaseCollections, 0)
	for rows.Next() {
		var tableName string

		if err := rows.Scan(&tableName); err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to process database result", err, nil)
		}

		result = append(result, utils.DatabaseCollections{TableName: tableName})
	}

	return result, nil
}
