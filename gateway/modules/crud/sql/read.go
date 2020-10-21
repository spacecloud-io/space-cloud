package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v8"
	"github.com/doug-martin/goqu/v8/exp"
	"github.com/spaceuptech/helpers"

	_ "github.com/denisenkom/go-mssqldb"                // Import for MsSQL
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postfres
	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for postgres

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// generateReadQuery makes a query for read operation
func (s *SQL) generateReadQuery(ctx context.Context, col string, req *model.ReadRequest) (string, []interface{}, error) {
	dbType := s.dbType
	if dbType == string(model.SQLServer) {
		dbType = string(model.Postgres)
	}

	dialect := goqu.Dialect(dbType)
	query := dialect.From(s.getDBName(col)).Prepared(true)
	var tarr []string
	if req.Find != nil {
		// Get the where clause from query object
		query, tarr = s.generateWhereClause(ctx, query, req.Find)
	}

	selArray := []interface{}{}
	if req.Options != nil {

		isJoin := len(req.Options.Join) > 0

		// Throw error if select is not provided during joins
		if isJoin && len(req.Options.Select) == 0 {
			return "", nil, errors.New("select cannot be nil when using joins")
		}

		// Check if the select clause exists
		if req.Options.Select != nil {
			for key := range req.Options.Select {
				if !isJoin {
					selArray = append(selArray, key)
					continue
				}

				selArray = append(selArray, goqu.I(key).As(strings.Join(strings.Split(key, "."), "__")))
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
			orderBys := make([]exp.OrderedExpression, len(req.Options.Sort))

			// Iterate over order array
			for i, value := range req.Options.Sort {
				// Add order type based on type attribute of order element
				var e exp.OrderedExpression
				if strings.HasPrefix(value, "-") {
					e = goqu.I(strings.TrimPrefix(value, "-")).Desc()
				} else {
					e = goqu.I(value).Asc()
				}

				// Append the order expression to the order expression array
				orderBys[i] = e
			}
			query = query.Order(orderBys...)
		}

		q, err := s.processJoins(ctx, col, query, req.Options.Join)
		if err != nil {
			return "", nil, err
		}
		query = q
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
	case utils.One:
		query = query.Select(selArray...)
	case utils.All:
		for function, colArray := range req.Aggregate {
			for _, column := range colArray {
				asColumnName := generateAggregateAsColumnName(function, column)
				switch function {
				case "sum":
					selArray = append(selArray, goqu.SUM(column).As(asColumnName))
				case "max":
					selArray = append(selArray, goqu.MAX(column).As(asColumnName))
				case "min":
					selArray = append(selArray, goqu.MIN(column).As(asColumnName))
				case "avg":
					selArray = append(selArray, goqu.AVG(column).As(asColumnName))
				case "count":
					selArray = append(selArray, goqu.COUNT("*").As(asColumnName))
				default:
					return "", nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unknown aggregate funcion (%s)", function), nil, map[string]interface{}{})
				}
			}
		}
		query = query.Select(selArray...)
		if len(req.GroupBy) > 0 {
			query = query.GroupBy(req.GroupBy...)
		}
	}

	// Generate the sql string and arguments
	sqlString, args, err := query.ToSQL()
	if err != nil {
		return "", nil, err
	}
	sqlString = strings.Replace(sqlString, "\"", "", -1)

	for _, v := range tarr {
		switch s.dbType {
		case "mysql":
			vReplaced := strings.Replace(v, "=", "REGEXP", -1)
			sqlString = strings.Replace(sqlString, v, vReplaced, -1)
		case "postgres":
			vReplaced := strings.Replace(v, "=", "~", -1)
			sqlString = strings.Replace(sqlString, v, vReplaced, -1)
		}

	}

	if s.dbType == string(model.SQLServer) {
		sqlString = s.generateQuerySQLServer(sqlString)
	}
	return sqlString, args, nil
}

func generateAggregateAsColumnName(function, column string) string {
	return fmt.Sprintf("%s___%s___%s", utils.GraphQLAggregate, function, column)
}

func splitAggregateAsColumnName(asColumnName string) (functionName string, columnName string, isAggregateColumn bool) {
	v := strings.Split(asColumnName, "___")
	if len(v) != 3 || !strings.HasPrefix(asColumnName, utils.GraphQLAggregate) {
		return "", "", false
	}
	return v[1], v[2], true
}

// Read query document(s) from the database
func (s *SQL) Read(ctx context.Context, col string, req *model.ReadRequest) (int64, interface{}, error) {
	return s.read(ctx, col, req, s.client)
}

func (s *SQL) read(ctx context.Context, col string, req *model.ReadRequest, executor executor) (int64, interface{}, error) {
	sqlString, args, err := s.generateReadQuery(ctx, col, req)
	if err != nil {
		return 0, nil, err
	}

	if col != utils.TableInvocationLogs && col != utils.TableEventingLogs {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Executing sql read query", map[string]interface{}{"sqlQuery": sqlString, "queryArgs": args})
	}

	return s.readexec(ctx, col, sqlString, args, executor, req)
}

func (s *SQL) readexec(ctx context.Context, col, sqlString string, args []interface{}, executor executor, req *model.ReadRequest) (int64, interface{}, error) {
	operation := req.Operation
	isAggregate := len(req.Aggregate) > 0

	stmt, err := executor.PreparexContext(ctx, sqlString)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = stmt.Close() }()

	rows, err := stmt.QueryxContext(ctx, args...)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = rows.Close() }()

	var rowTypes []*sql.ColumnType

	switch s.GetDBType() {
	case model.MySQL, model.Postgres, model.SQLServer:
		rowTypes, _ = rows.ColumnTypes()
	}

	switch operation {
	case utils.Count:
		mapping := make(map[string]interface{})
		if !rows.Next() {
			return 0, nil, errors.New("SQL: No response from db")
		}

		err := rows.MapScan(mapping)
		if err != nil {
			return 0, nil, err
		}

		switch s.GetDBType() {
		case model.MySQL, model.Postgres, model.SQLServer:
			mysqlTypeCheck(ctx, s.GetDBType(), rowTypes, mapping)
		}

		for _, v := range mapping {
			return v.(int64), v, nil
		}

		return 0, nil, errors.New("unknown error occurred")

	case utils.One:
		mapping := make(map[string]interface{})
		if !rows.Next() {
			return 0, nil, errors.New("SQL: No response from db")
		}

		err := rows.MapScan(mapping)
		if err != nil {
			return 0, nil, err
		}

		switch s.GetDBType() {
		case model.MySQL, model.Postgres, model.SQLServer:
			mysqlTypeCheck(ctx, s.GetDBType(), rowTypes, mapping)
		}

		return 1, mapping, nil

	case utils.All, utils.Distinct:
		array := []interface{}{}
		mapping := make(map[string]map[string]interface{})
		var count int64
		for rows.Next() {

			// Increment the counter
			count++

			row := make(map[string]interface{})
			err := rows.MapScan(row)
			if err != nil {
				return 0, nil, err
			}
			switch s.GetDBType() {
			case model.MySQL, model.Postgres, model.SQLServer:
				mysqlTypeCheck(ctx, s.GetDBType(), rowTypes, row)
			}
			if isAggregate {
				funcMap := map[string]interface{}{}
				for asColumnName, value := range row {
					functionName, columnName, isAggregateColumn := splitAggregateAsColumnName(asColumnName)
					if isAggregateColumn {
						delete(row, asColumnName)
						// check if function name already exists
						funcValue, ok := funcMap[functionName]
						if !ok {
							// set new function
							// NOTE: This case occurs for count function with no column name (using * operator instead)
							if columnName == "" {
								funcMap[functionName] = value
							} else {
								funcMap[functionName] = map[string]interface{}{columnName: value}
							}
							continue
						}
						// add new column to existing function
						funcValue.(map[string]interface{})[columnName] = value
					}
				}
				if len(funcMap) > 0 {
					row[utils.GraphQLAggregate] = funcMap
				}
			}

			if req.Options == nil || req.Options.ReturnType == "flat" {
				array = append(array, row)
				continue
			}

			processRows(col, row, req.Options.Join, mapping, &array)
		}

		return count, array, nil

	default:
		return 0, nil, utils.ErrInvalidParams
	}
}

func processRows(table string, row map[string]interface{}, join []model.JoinOption, mapping map[string]map[string]interface{}, finalArray *[]interface{}) {
	m := map[string]interface{}{}

	// Get keys of this table
	for k, v := range row {
		if strings.Split(k, "__")[0] == table {
			m[k] = v
		}
	}

	// Generate unique key for row. fmt.Sprintf internally sorts all keys
	// hence returns a deterministic key.
	key := fmt.Sprintf("%v", m)

	// Check if key exists in mapping. This can happen if the row has multiple
	// sub rows else append self to final array.
	if m2, p := mapping[key]; p {
		m = m2
	} else {
		mapping[key] = m
		*finalArray = append(*finalArray, m)
	}

	// Process joint tables for rows
	for _, j := range join {
		var arr []interface{}

		// Check if table name is already present in parent row. If not, create a new array
		if arrTemp, p := m[j.Table]; p {
			arr = arrTemp.([]interface{})
		} else {
			arr = []interface{}{}
		}

		// Recursively call the same function again
		processRows(j.Table, row, j.Join, mapping, &arr)
	}
}
