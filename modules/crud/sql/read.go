package sql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	goqu "github.com/doug-martin/goqu/v8"
	"github.com/doug-martin/goqu/v8/exp"

	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for postgres
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postfres

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// generateReadQuery makes a query for read operation
func (s *SQL) generateReadQuery(ctx context.Context, project, col string, req *model.ReadRequest) (string, []interface{}, error) {
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
	
	selArray := []interface{}{}

	if req.Options != nil {
		// Check if the select clause exists
		if req.Options.Select != nil {
			for key := range req.Options.Select {
				selArray = append(selArray, key)
			}
		}

		if req.Options.Skip != nil {
			query = query.Offset(uint(*req.Options.Skip))
		}

		if req.Options.Limit != nil {
			query = query.Limit(uint(*req.Options.Limit))
		}

		if req.Options.Sort != nil {
			// Format the order array to a suitable type
			orderMap := req.Options.Sort

			orderBys := []exp.OrderedExpression{}

			// Iterate over order array
			for k, value := range orderMap {
				// Add order type based on type attribute of order element
				var exp exp.OrderedExpression
				if value < 0 {
					exp = goqu.I(k).Desc()
				} else {
					exp = goqu.I(k).Asc()
				}

				// Append the order expression to the order expression array
				orderBys = append(orderBys, exp)
			}
			query = query.Order(orderBys...)
		}
	}
	
	switch req.Operation {
	case utils.Count:
		query = query.Select(goqu.COUNT("*"))
	case utils.Distinct:
		distinct := req.Options.Distinct
		if distinct == nil {
			return "", nil, utils.ErrInvalidParams
		}
		query = query.SelectDistinct(*distinct)
	case utils.One, utils.All:
		query = query.Select(selArray...)
	}

	// Generate the sql string and arguments
	sqlString, args, err := query.ToSQL()
	if err != nil {
		return "", nil, err
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)
	return sqlString, args, nil
}
// Read query document(s) from the database
func (s *SQL) Read(ctx context.Context, project, col string, req *model.ReadRequest) (interface{}, error) {
	sqlString, args, err := s.generateReadQuery(ctx, project, col, req)
	if err != nil {
		return nil, err
	}
	
	stmt, err := s.client.PreparexContext(ctx, sqlString)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rowTypes []*sql.ColumnType

	switch s.GetDBType() {
	case utils.MySQL, utils.Postgres:
		rowTypes, _ = rows.ColumnTypes()
	}

	switch req.Operation {
	case utils.Count:
		mapping := make(map[string]interface{})
		if !rows.Next() {
			return nil, errors.New("SQL: No response from db")
		}

		err := rows.MapScan(mapping)
		if err != nil {
			return nil, err
		}

		switch s.GetDBType() {
		case utils.MySQL, utils.Postgres:
			mysqlTypeCheck(rowTypes, mapping)
		}

		for _, v := range mapping {
			return v, nil
		}
		return nil, nil

	case utils.One:
		mapping := make(map[string]interface{})
		if !rows.Next() {
			return nil, errors.New("SQL: No response from db")
		}

		err := rows.MapScan(mapping)
		if err != nil {
			return nil, err
		}

		switch s.GetDBType() {
		case utils.MySQL, utils.Postgres:
			mysqlTypeCheck(rowTypes, mapping)
		}

		return mapping, nil

	case utils.All, utils.Distinct:
		array := []interface{}{}
		for rows.Next() {
			mapping := make(map[string]interface{})
			err := rows.MapScan(mapping)
			if err != nil {
				return nil, err
			}

			switch s.GetDBType() {
			case utils.MySQL, utils.Postgres:
				mysqlTypeCheck(rowTypes, mapping)
			}

			array = append(array, mapping)
		}

		return array, nil

	default:
		return nil, utils.ErrInvalidParams
	}
}
