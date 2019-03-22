package sql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	goqu "gopkg.in/doug-martin/goqu.v4"

	_ "github.com/go-sql-driver/mysql"                 // Import for MySQL
	_ "github.com/lib/pq"                              // Import for postgres
	_ "gopkg.in/doug-martin/goqu.v4/adapters/postgres" // Adapter for postfres

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Read query document(s) from the database
func (s *SQL) Read(ctx context.Context, project, col string, req *model.ReadRequest) (interface{}, error) {
	// Generate a prepared query builder
	query := goqu.From(col).Prepared(true)
	query = query.SetAdapter(goqu.NewAdapter(s.dbType, query))

	if req.Find != nil {
		// Get the where clause from query object
		var err error
		query, err = generateWhereClause(query, req.Find)
		if err != nil {
			return nil, err
		}
	}

	if req.Options != nil {
		// Check if the select clause exists
		if req.Options.Select != nil {
			selArray := []interface{}{}
			for key := range req.Options.Select {
				selArray = append(selArray, key)
			}
			query = query.Select(selArray...)
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

			orderBys := []goqu.OrderedExpression{}

			// Iterate over order array
			for k, value := range orderMap {
				// Add order type based on type attribute of order element
				var exp goqu.OrderedExpression
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

	// Generate the sql string and arguments
	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)

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
	case utils.MySQL:
		rowTypes, _ = rows.ColumnTypes()
	}

	switch req.Operation {
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
		case utils.MySQL:
			mysqlTypeCheck(rowTypes, mapping)
		}

		return mapping, nil

	case "all":
		array := []interface{}{}
		for rows.Next() {
			mapping := make(map[string]interface{})
			err := rows.MapScan(mapping)
			if err != nil {
				return nil, err
			}

			switch s.GetDBType() {
			case utils.MySQL:
				rowTypes, _ = rows.ColumnTypes()
			}

			switch s.GetDBType() {
			case utils.MySQL:
				mysqlTypeCheck(rowTypes, mapping)
			}

			array = append(array, mapping)
		}

		return array, nil

	default:
		return nil, utils.ErrInvalidParams
	}
}
