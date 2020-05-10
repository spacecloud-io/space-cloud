package sql

import (
	"context"
	"strings"

	"github.com/doug-martin/goqu/v8"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"

	_ "github.com/denisenkom/go-mssqldb"                // Import for MsSQL
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postgres
	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for postgres
)

// Delete removes the document(s) from the database which match the condition
func (s *SQL) Delete(ctx context.Context, col string, req *model.DeleteRequest) (int64, error) {
	sqlString, args, err := s.generateDeleteQuery(col, req)
	if err != nil {
		return 0, err
	}
	res, err := doExecContext(ctx, sqlString, args, s.client)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// genrateDeleteQuery makes query for delete operation
func (s *SQL) generateDeleteQuery(col string, req *model.DeleteRequest) (string, []interface{}, error) {
	// Generate a prepared query builder

	dbType := s.dbType
	if dbType == string(utils.SQLServer) {
		dbType = string(utils.Postgres)
	}

	dialect := goqu.Dialect(dbType)
	query := dialect.From(s.getDBName(col)).Prepared(true)

	if req.Find != nil {
		// Get the where clause from query object
		query, _ = s.generateWhereClause(query, req.Find)
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.Delete().ToSQL()
	if err != nil {
		return "", nil, err
	}
	sqlString = strings.Replace(sqlString, "\"", "", -1)

	if s.dbType == string(utils.SQLServer) {
		sqlString = s.generateQuerySQLServer(sqlString)
	}
	return sqlString, args, nil
}

// DeleteCollection drops a table
func (s *SQL) DeleteCollection(ctx context.Context, col string) error {
	query := "DROP TABLE " + s.getDBName(col)
	_, err := s.client.ExecContext(ctx, query, []interface{}{}...)
	return err
}
