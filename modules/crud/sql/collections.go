package sql

import (
	"context"
	"strings"

	"github.com/doug-martin/goqu/v8"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetCollections returns collection / tables name of specified database
func (s *SQL) GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error) {
	dialect := goqu.Dialect(s.dbType)
	query := dialect.From("information_schema.tables").Prepared(true).Select("table_name").Where(goqu.Ex{"table_schema": project})

	sqlString, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)
	rows, err := s.client.Queryx(sqlString, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]utils.DatabaseCollections, 0)
	for rows.Next() {
		var tableName string

		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		result = append(result, utils.DatabaseCollections{TableName: tableName})
	}

	return result, nil
}
