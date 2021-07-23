package sql

import (
	"context"
	"strings"

	"github.com/doug-martin/goqu/v8"

	_ "github.com/denisenkom/go-mssqldb"                // Import for MsSQL
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postgres
	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for postgres

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Delete removes the document(s) from the database which match the condition
func (s *SQL) Delete(ctx context.Context, col string, req *model.DeleteRequest) (int64, error) {
	sqlString, args, err := s.generateDeleteQuery(ctx, req, col)
	if err != nil {
		return 0, err
	}
	res, err := doExecContext(ctx, sqlString, args, s.getClient())
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// genrateDeleteQuery makes query for delete operation
func (s *SQL) generateDeleteQuery(ctx context.Context, req *model.DeleteRequest, col string) (string, []interface{}, error) {
	// Generate a prepared query builder

	dbType := s.dbType
	if dbType == string(model.SQLServer) {
		dbType = string(model.Postgres)
	}

	dialect := goqu.Dialect(dbType)
	query := dialect.From(s.getColName(col)).Prepared(true)

	if req.Find != nil {
		// Get the where clause from query object
		query = s.generateWhereClause(ctx, query, req.Find, nil)
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.Delete().ToSQL()
	if err != nil {
		return "", nil, err
	}
	sqlString = strings.Replace(sqlString, "\"", "", -1)

	if s.dbType == string(model.SQLServer) {
		sqlString = s.generateQuerySQLServer(sqlString)
	}
	return sqlString, args, nil
}

// DeleteCollection drops a table
func (s *SQL) DeleteCollection(ctx context.Context, col string) error {
	query := "DROP TABLE " + s.getColName(col)
	_, err := s.getClient().ExecContext(ctx, query, []interface{}{}...)
	return err
}
