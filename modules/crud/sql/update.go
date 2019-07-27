package sql

import (
	"context"
	"strings"

	goqu "github.com/doug-martin/goqu/v8"

	_ "github.com/go-sql-driver/mysql"                 // Import for MySQL
	_ "github.com/lib/pq"                              // Import for postgres
	_ "github.com/doug-martin/goqu/v8/dialect/postgres"  // Dialect for postgres

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Update updates the document(s) which match the condition provided.
func (s *SQL) Update(ctx context.Context, project, col string, req *model.UpdateRequest) error {
	if req == nil {
		return utils.ErrInvalidParams
	}
	sqlString, args, err := s.generateUpdateQuery(ctx, project, col, req)
	if err != nil {
		return err
	}
	return s.doExecContext(ctx, sqlString, args)
}

//generateUpdateQuery makes query for update operation
func (s *SQL) generateUpdateQuery(ctx context.Context, project, col string, req *model.UpdateRequest) (string, []interface{}, error) {
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

	if req.Update == nil {
		return "", nil, utils.ErrInvalidParams
	}

	record, err := generateRecord(req.Update["$set"])
	if err != nil {
		return "", nil, err
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.Update().Set(record).ToSQL()
	if err != nil {
		return "", nil, err
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)
	return sqlString, args, nil
}
