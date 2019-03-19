package sql

import (
	"context"
	"strings"

	"github.com/spaceuptech/space-cloud/model"
	goqu "gopkg.in/doug-martin/goqu.v4"

	_ "github.com/go-sql-driver/mysql"                 // Import for MySQL
	_ "github.com/lib/pq"                              // Import for postgres
	_ "gopkg.in/doug-martin/goqu.v4/adapters/postgres" // Adapter for postgres
)

// Delete removes the document(s) from the database which match the condition
func (s *SQL) Delete(ctx context.Context, project, col string, req *model.DeleteRequest) error {
	sqlString, args, err := s.GenerateDeleteQuery(ctx, project, col, req)
	if err != nil {
		return err
	}
	return s.doExec(sqlString, args)
}

//GenrateDeleteQuery makes query for delete operation
func (s *SQL) GenerateDeleteQuery(ctx context.Context, project, col string, req *model.DeleteRequest) (string, []interface{}, error) {
	// Generate a prepared query builder
	query := goqu.From(col).Prepared(true)
	query = query.SetAdapter(goqu.NewAdapter(s.dbType, query))

	if req.Find != nil {
		// Get the where clause from query object
		var err error
		query, err = generateWhereClause(query, req.Find)
		if err != nil {
			return "", nil, err
		}
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.ToDeleteSql()
	if err != nil {
		return "", nil, err
	}
	sqlString = strings.Replace(sqlString, "\"", "", -1)
	return sqlString, args, nil
}
