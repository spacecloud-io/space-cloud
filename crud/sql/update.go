package sql

import (
	"context"
	"strings"

	goqu "gopkg.in/doug-martin/goqu.v4"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Update updates the document(s) which match the condition provided.
func (s *SQL) Update(ctx context.Context, project, col string, req *model.UpdateRequest) error {
	// Generate a prepared query builder
	query := goqu.From(col).Prepared(true)
	query = query.SetAdapter(goqu.NewAdapter(s.dbType, query))

	if req.Find != nil {

		// Get the where clause from query object
		var err error
		query, err = generateWhereClause(query, req.Find)
		if err != nil {
			return err
		}

	}

	if req.Update == nil {
		return utils.ErrInvalidParams
	}

	record, err := generateRecord(req.Update)
	if err != nil {
		return err
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.ToUpdateSql(record)
	if err != nil {
		return err
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)

	return s.doExec(sqlString, args)
}
