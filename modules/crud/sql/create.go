package sql

import (
	"context"
	"strings"

	goqu "github.com/doug-martin/goqu/v8"

	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postgres
	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for postgres

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Create inserts a document (or multiple when op is "all") into the database
func (s *SQL) Create(ctx context.Context, project, col string, req *model.CreateRequest) error {
	sqlQuery, args, err := s.generateCreateQuery(ctx, project, col, req)
	if err != nil {
		return err
	}
	_, err = doExecContext(ctx, sqlQuery, args, s.client)
	return err
}

//generateCreateQuery makes query for create operation
func (s *SQL) generateCreateQuery(ctx context.Context, project, col string, req *model.CreateRequest) (string, []interface{}, error) {
	// Generate a prepared query builder
	dialect := goqu.Dialect(s.dbType)
	query := dialect.From(project + "." + col).Prepared(true)

	var insert []interface{}
	if req.Operation == "one" {
		insert = []interface{}{req.Document.(map[string]interface{})}
	} else {
		var ok bool
		insert, ok = req.Document.([]interface{})
		if !ok {
			return "", nil, utils.ErrInvalidParams
		}
	}

	// Iterate over records to be inserted
	records := []interface{}{}
	for _, temp := range insert {
		// Genrate a record out of object
		record, err := generateRecord(temp)
		if err != nil {
			return "", nil, err
		}

		// Append record to records array
		records = append(records, record)
	}

	sqlQuery, args, err := query.Insert().Rows(records).ToSQL()
	if err != nil {
		return "", nil, err
	}

	sqlQuery = strings.Replace(sqlQuery, "\"", "", -1)
	return sqlQuery, args, nil
}
