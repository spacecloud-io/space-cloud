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
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postgres
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

	if req.Options == nil {
		req.Options = &model.ReadOptions{}
	}
	if req.Options.Limit == nil {
		req.Options.Limit = s.queryFetchLimit
		req.Options.HasOptions = true
	}

	dialect := goqu.Dialect(dbType)
	query := dialect.From(s.getColName(col)).Prepared(true)
	var regexArr []string
	// Get the where clause from query object
	query, regexArr = s.generateWhereClause(ctx, query, req.Find, req.MatchWhere)

	selArray := make([]interface{}, 0)
	if req.Options != nil {

		isJoin := len(req.Options.Join) > 0

		// Throw error if select is not provided during joins
		if isJoin && len(req.Options.Select) == 0 && len(req.Aggregate) == 0 {
			return "", nil, errors.New("select cannot be nil when using joins")
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

		q, err := s.processJoins(ctx, query, req.Options.Join, req.Options.Select, len(req.Aggregate) > 0)
		if err != nil {
			return "", nil, err
		}
		query = q
		// Check if the select clause exists
		if req.Options.Select != nil {
			for key := range req.Options.Select {
				if !isJoin {
					selArray = append(selArray, key)
					continue
				}

				arr := strings.Split(key, ".")
				selArray = append(selArray, goqu.I(key).As(strings.Join(arr, "__")))
			}
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
	case utils.One:
		query = query.Select(selArray...)
	case utils.All:
		for function, colArray := range req.Aggregate {
			for _, column := range colArray {
				asColumnName := getAggregateAsColumnName(function, column)
				switch function {
				case "sum":
					selArray = append(selArray, goqu.SUM(getAggregateColumnName(column)).As(asColumnName))
				case "max":
					selArray = append(selArray, goqu.MAX(getAggregateColumnName(column)).As(asColumnName))
				case "min":
					selArray = append(selArray, goqu.MIN(getAggregateColumnName(column)).As(asColumnName))
				case "avg":
					selArray = append(selArray, goqu.AVG(getAggregateColumnName(column)).As(asColumnName))
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

	if model.DBType(s.dbType) == model.SQLServer {
		sqlString = mutateSQLServerLimitAndOffsetOperation(sqlString, req)
	}

	for _, v := range regexArr {
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

func getAggregateColumnName(column string) string {
	return strings.Split(column, ":")[1]
}

func getAggregateAsColumnName(function, column string) string {
	format := "nested"
	arr := strings.Split(column, ":")

	returnField := arr[0]
	column = arr[1]
	if len(arr) == 3 && arr[2] == "table" {
		format = "table"
	}

	return fmt.Sprintf("%s___%s___%s___%s___%s", utils.GraphQLAggregate, format, returnField, function, strings.Join(strings.Split(column, "."), "__"))
}

func splitAggregateAsColumnName(asColumnName string) (format, returnField, functionName, columnName string, isAggregateColumn bool) {
	v := strings.Split(asColumnName, "___")
	if len(v) != 5 || !strings.HasPrefix(asColumnName, utils.GraphQLAggregate) {
		return "", "", "", "", false
	}
	return v[1], v[2], v[3], v[4], true
}

// Read query document(s) from the database
func (s *SQL) Read(ctx context.Context, col string, req *model.ReadRequest) (int64, interface{}, map[string]map[string]string, error) {
	return s.read(ctx, col, req, s.getClient())
}

func (s *SQL) read(ctx context.Context, col string, req *model.ReadRequest, executor executor) (int64, interface{}, map[string]map[string]string, error) {
	sqlString, args, err := s.generateReadQuery(ctx, col, req)
	if err != nil {
		return 0, nil, nil, err
	}

	if col != utils.TableInvocationLogs && col != utils.TableEventingLogs {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Executing sql read query", map[string]interface{}{"sqlQuery": sqlString, "queryArgs": args})
	}

	return s.readExec(ctx, col, sqlString, args, executor, req)
}

func (s *SQL) readExec(ctx context.Context, col, sqlString string, args []interface{}, executor executor, req *model.ReadRequest) (int64, interface{}, map[string]map[string]string, error) {
	operation := req.Operation
	isAggregate := len(req.Aggregate) > 0

	stmt, err := executor.PreparexContext(ctx, sqlString)
	if err != nil {
		return 0, nil, nil, err
	}
	defer func() { _ = stmt.Close() }()

	rows, err := stmt.QueryxContext(ctx, args...)
	if err != nil {
		return 0, nil, nil, err
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
			return 0, nil, nil, errors.New("SQL: No response from db")
		}

		err := rows.MapScan(mapping)
		if err != nil {
			return 0, nil, nil, err
		}

		switch s.GetDBType() {
		case model.MySQL, model.Postgres, model.SQLServer:
			mysqlTypeCheck(ctx, s.GetDBType(), rowTypes, mapping)
		}

		for _, v := range mapping {
			return v.(int64), v, make(map[string]map[string]string), nil
		}

		return 0, nil, nil, errors.New("unknown error occurred")

	case utils.One:
		mapping := make(map[string]interface{})
		if !rows.Next() {
			return 0, nil, nil, errors.New("SQL: No response from db")
		}

		err := rows.MapScan(mapping)
		if err != nil {
			return 0, nil, nil, err
		}

		switch s.GetDBType() {
		case model.MySQL, model.Postgres, model.SQLServer:
			mysqlTypeCheck(ctx, s.GetDBType(), rowTypes, mapping)
		}

		processAggregate(mapping, mapping, col, isAggregate)
		if req.PostProcess != nil {
			_ = s.auth.PostProcessMethod(ctx, req.PostProcess[col], mapping)
		}

		return 1, mapping, make(map[string]map[string]string), nil

	case utils.All, utils.Distinct:
		array := make([]interface{}, 0)
		mapping := make(map[string]map[string]interface{})
		jointMapping := make(map[string]map[string]string)
		var count int64
		for rows.Next() {

			// Increment the counter
			count++

			row := make(map[string]interface{})
			err := rows.MapScan(row)
			if err != nil {
				return 0, nil, nil, err
			}

			switch s.GetDBType() {
			case model.MySQL, model.Postgres, model.SQLServer:
				mysqlTypeCheck(ctx, s.GetDBType(), rowTypes, row)
			}

			if req.Options == nil || req.Options.ReturnType == "table" || len(req.Options.Join) == 0 {
				processAggregate(row, row, "*", isAggregate)
				array = append(array, row)
				continue
			}

			s.processRows(ctx, []string{col}, isAggregate, row, req.Options.Join, mapping, &array, req.PostProcess, jointMapping)
		}

		return count, array, jointMapping, nil

	default:
		return 0, nil, nil, utils.ErrInvalidParams
	}
}
func processAggregate(row, m map[string]interface{}, tableName string, isAggregate bool) {
	if isAggregate {
		funcMap := map[string]interface{}{}
		for asColumnName, value := range row {
			format, returnField, functionName, columnName, isAggregateColumn := splitAggregateAsColumnName(asColumnName)
			if isAggregateColumn {
				// Only process aggregated field if it belongs to the current table
				if arr := strings.Split(columnName, "__"); len(arr) == 2 {
					if arr[0] != tableName && tableName != "*" {
						continue
					}
				}
				delete(row, asColumnName)

				if format == "table" {
					m[returnField] = value
					continue
				}

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
			m[utils.GraphQLAggregate] = funcMap
		}
	}
}
func (s *SQL) processRows(ctx context.Context, table []string, isAggregate bool, row map[string]interface{}, join []*model.JoinOption, mapping map[string]map[string]interface{}, finalArray *[]interface{}, postProcess map[string]*model.PostProcess, joinMapping map[string]map[string]string) {
	// row obtained from database contains flattened result of all tables(if join was specified)
	// m variable will only store result of specific table
	m := map[string]interface{}{}
	keyMap := map[string]interface{}{}

	length := len(table) - 1

	// Get keys of this table
	for k, v := range row {
		a := strings.Split(k, "__")
		if a[0] == utils.GraphQLAggregate {
			_, _, _, columnName, _ := splitAggregateAsColumnName(k)
			if arr := strings.Split(columnName, "__"); len(arr) == 2 {
				a[0] = arr[0] // table name
				a[1] = arr[1] // column name
			}
		}
		if utils.StringExists(table, a[0]) {
			keyMap[a[1]] = v
		}
		if table[length] == a[0] {
			m[a[1]] = v
		}
	}

	// Generate unique key for row. fmt.Sprintf internally sorts all keys
	// hence returns a deterministic key.
	key := fmt.Sprintf("%v", keyMap)

	// Check if key exists in mapping. This can happen if the row has multiple
	// sub rows else append self to final array.
	var mapLength int
	if m2, p := mapping[key]; p {
		mapLength = len(m2)
		m = m2
	} else {
		mapLength = len(m)
		mapping[key] = m
		*finalArray = append(*finalArray, m)

		// Perform post processing
		if postProcess != nil {
			_ = s.auth.PostProcessMethod(ctx, postProcess[table[length]], m)
		}

		// Process aggregate field only if its the root table that we are processing
		processAggregate(row, m, table[length], isAggregate)
	}

	if mapLength == 0 {
		return
	}

	// Process joint tables for rows
	for _, j := range join {
		var arr []interface{}
		utils.GenerateJoinKeys(j.Table, j.On, row, joinMapping)
		// Check if table name is already present in parent row. If not, create a new array
		if arrTemp, p := m[j.Table]; p {
			switch t := arrTemp.(type) {
			case []interface{}:
				arr = t
			case map[string]interface{}:
				arr = []interface{}{t}
			}
		} else {
			arr = []interface{}{}
		}

		// Recursively call the same function again
		s.processRows(ctx, append(table, j.Table), isAggregate, row, j.Join, mapping, &arr, postProcess, joinMapping)
		if j.Op == utils.All || j.Op == "" {
			m[j.Table] = arr
		} else {
			if len(arr) > 0 {
				m[j.Table] = arr[0]
			} else {
				m[j.Table] = map[string]interface{}{}
			}
		}
	}
}
