package sql

import (
	"context"

	"github.com/spaceuptech/space-cloud/utils"
)

// not to be exposed externally
// ExecuteRawQuery return a structure of sql table
func (s *SQL) DescribeTable(ctx context.Context, project, col string) ([]utils.FieldType, []utils.ForeignKeysType, error) {
	fields, err := s.getDescribeDetails(ctx, col)
	if err != nil {
		return nil, nil, err
	}
	foreignKeys, err := s.getForeignKeyDetails(ctx, project, col)
	if err != nil {
		return nil, nil, err
	}
	return fields, foreignKeys, nil
}

func (s *SQL) getDescribeDetails(ctx context.Context, col string) ([]utils.FieldType, error) {
	rows, err := s.client.Queryx("DESCRIBE " + col)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []utils.FieldType{}
	for rows.Next() {
		fieldType := new(utils.FieldType)

		if err := rows.StructScan(fieldType); err != nil {
			return nil, err
		}

		result = append(result, *fieldType)
	}
	return result, nil
}

func (s *SQL) getForeignKeyDetails(ctx context.Context, project, col string) ([]utils.ForeignKeysType, error) {
	rows, err := s.client.Queryx("select TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = '" + project + "' and TABLE_NAME = '" + col + "'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []utils.ForeignKeysType{}
	for rows.Next() {
		foreignKey := new(utils.ForeignKeysType)

		if err := rows.StructScan(foreignKey); err != nil {
			return nil, err
		}

		result = append(result, *foreignKey)
	}
	return result, nil
}
