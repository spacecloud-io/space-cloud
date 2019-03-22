package sql

import (
	"context"
	"strings"

	goqu "gopkg.in/doug-martin/goqu.v4"

	_ "github.com/go-sql-driver/mysql"                 // Import for MySQL
	_ "github.com/lib/pq"                              // Import for postgres
	_ "gopkg.in/doug-martin/goqu.v4/adapters/postgres" // Adapter for postgres

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Update updates the document(s) which match the condition provided.
func (s *SQL) Update(ctx context.Context, project, col string, req *model.UpdateRequest) error {
	sqlString, args, err := s.GenerateUpdateQuery(ctx, project, col, req)
	if err != nil {
		return err
	}
	return s.doExec(sqlString, args)
}

//generateUpdateQuery makes query for update operation
func (s *SQL) GenerateUpdateQuery(ctx context.Context, project, col string, req *model.UpdateRequest) (string, []interface{}, error) {
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

	if req.Update == nil {
		return "", nil, utils.ErrInvalidParams
	}

	record, err := generateRecord(req.Update["$set"])
	if err != nil {
		return "", nil, err
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.ToUpdateSql(record)
	if err != nil {
		return "", nil, err
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)
	return sqlString, args, nil
}
