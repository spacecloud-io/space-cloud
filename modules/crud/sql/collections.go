package sql

import (
	"context"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetCollections returns collection / tables name of specified database
func (s *SQL) GetCollections(ctx context.Context, project, dbType string) ([]utils.DatabaseCollections, error) {
	queryString := `SELECT table_name FROM information_schema.tables WHERE table_schema = $1`

	rows, err := s.client.Queryx(queryString, []interface{}{project}...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []utils.DatabaseCollections{}
	count := 0
	for rows.Next() {
		count++
		fieldType := new(utils.DatabaseCollections)

		if err := rows.StructScan(fieldType); err != nil {
			return nil, err
		}

		result = append(result, *fieldType)
	}
	
	return result, nil
}
