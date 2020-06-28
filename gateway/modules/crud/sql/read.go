package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v8"
	"github.com/doug-martin/goqu/v8/exp"
	"github.com/sirupsen/logrus"

	_ "github.com/denisenkom/go-mssqldb"                // Import for MsSQL
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postfres
	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for postgres

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// generateReadQuery makes a query for read operation
func (s *SQL) generateReadQuery(col string, req *model.ReadRequest) (string, []interface{}, error) {
	dbType := s.dbType
	if dbType == string(utils.SQLServer) {
		dbType = string(utils.Postgres)
	}

	dialect := goqu.Dialect(dbType)
	query := dialect.From(s.getDBName(col)).Prepared(true)
	var tarr []string
	if req.Find != nil {
		// Get the where clause from query object
		query, tarr = s.generateWhereClause(query, req.Find)
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
					return "", nil, utils.LogError(fmt.Sprintf(`Unknown aggregate funcion "%s"`, function), "sql", "generateReadQuery", nil)
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

	if s.dbType == string(utils.SQLServer) {
		sqlString = s.generateQuerySQLServer(sqlString)
	}
	// logrus.Println("Log Query", sqlString)
	return sqlString, args, nil
}
func generateAggregateAsColumnName(function, column string) string {
	return fmt.Sprintf("%s__%s__%s", utils.GraphQLAggregate, function, column)
}

func splitAggregateAsColumnName(asColumnName string) (functionName string, columnName string, isAggregateColumn bool) {
	v := strings.Split(asColumnName, "__")
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
	sqlString, args, err := s.generateReadQuery(col, req)
	if err != nil {
		return 0, nil, err
	}

	logrus.Debugf("Executing sql read query: %s - %v", sqlString, args)

	return s.readexec(ctx, sqlString, args, req.Operation, executor, len(req.Aggregate) > 0)
}

func (s *SQL) readexec(ctx context.Context, sqlString string, args []interface{}, operation string, executor executor, isAggregate bool) (int64, interface{}, error) {
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
	case utils.MySQL, utils.Postgres, utils.SQLServer:
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
		case utils.MySQL, utils.Postgres, utils.SQLServer:
			mysqlTypeCheck(s.GetDBType(), rowTypes, mapping)
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
		case utils.MySQL, utils.Postgres, utils.SQLServer:
			mysqlTypeCheck(s.GetDBType(), rowTypes, mapping)
		}

		return 1, mapping, nil

	case utils.All, utils.Distinct:
		array := []interface{}{}
		var count int64
		for rows.Next() {

			// Increment the counter
			count++

			mapping := make(map[string]interface{})
			err := rows.MapScan(mapping)
			if err != nil {
				return 0, nil, err
			}
			switch s.GetDBType() {
			case utils.MySQL, utils.Postgres, utils.SQLServer:
				mysqlTypeCheck(s.GetDBType(), rowTypes, mapping)
			}
			if isAggregate {
				funcMap := map[string]interface{}{}
				for asColumnName, value := range mapping {
					functionName, columnName, isAggregateColumn := splitAggregateAsColumnName(asColumnName)
					if isAggregateColumn {
						delete(mapping, asColumnName)
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
					mapping[utils.GraphQLAggregate] = funcMap
				}
			}

			array = append(array, mapping)
		}

		return count, array, nil

	default:
		return 0, nil, utils.ErrInvalidParams
	}
}
