package sql

import (
	"context"
	"strings"

	"github.com/spaceuptech/space-cloud/model"
	goqu "github.com/doug-martin/goqu/v8"

	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for postgres
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postgres
)

// Delete removes the document(s) from the database which match the condition
func (s *SQL) Delete(ctx context.Context, project, col string, req *model.DeleteRequest) error {
	sqlString, args, err := s.generateDeleteQuery(ctx, project, col, req)
	if err != nil {
		return err
	}
	_, err = doExecContext(ctx, sqlString, args, s.client)
	return err
}

//genrateDeleteQuery makes query for delete operation
func (s *SQL) generateDeleteQuery(ctx context.Context, project, col string, req *model.DeleteRequest) (string, []interface{}, error) {
	// Generate a prepared query builder
	dialect := goqu.Dialect(s.dbType)
	query := dialect.From(col).Prepared(true)

	if req.Find != nil {
		// Get the where clause from query object
		var err error
		query, err = generateWhereClause(query, req.Find)
		if err != nil {
			return "", nil, err
		}
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.Delete().ToSQL()
	if err != nil {
		return "", nil, err
	}
	sqlString = strings.Replace(sqlString, "\"", "", -1)
	return sqlString, args, nil
}
