package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/doug-martin/goqu/v8"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (s *SQL) generator(ctx context.Context, find map[string]interface{}, isJoin bool) goqu.Expression {
	array := []goqu.Expression{}
	for k, v := range find {
		if strings.HasPrefix(k, "$or") {
			orArray := v.([]interface{})
			orFinalArray := []goqu.Expression{}
			for _, item := range orArray {
				f2 := item.(map[string]interface{})

				// Add an always match case if or had an empty find. We do this so that sql generator
				// doesn't ignore something like this
				if len(f2) == 0 {
					orFinalArray = append(orFinalArray, goqu.I("1").Eq(goqu.I("1")))
					continue
				}

				exp := s.generator(ctx, f2, isJoin)
				orFinalArray = append(orFinalArray, exp)
			}

			array = append(array, goqu.Or(orFinalArray...))
			continue
		}

		val, isObj := v.(map[string]interface{})
		if isObj {
			for k2, v2 := range val {
				if vString, p := v2.(string); p && isJoin {
					v2 = goqu.I(vString)
				}
				switch k2 {
				case "$regex":
					switch s.dbType {
					case "postgres":
						array = append(array, goqu.L(fmt.Sprintf("(%s ~ ?)", k), v2))
					case "mysql":
						array = append(array, goqu.L(fmt.Sprintf("(%s REGEXP ?)", k), v2))
					}

				case "$like":
					array = append(array, goqu.I(k).Like(v2))
				case "$eq":
					array = append(array, goqu.I(k).Eq(v2))
				case "$ne":
					array = append(array, goqu.I(k).Neq(v2))
				case "$contains":
					data, err := json.Marshal(v2)
					if err != nil {
						_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error marshalling data $contains data", err, nil)
						break
					}
					switch s.dbType {
					case string(model.MySQL):
						array = append(array, goqu.L(fmt.Sprintf("json_contains(%s,?)", k), string(data)))
					case string(model.Postgres):
						array = append(array, goqu.L(fmt.Sprintf("%s @> ?", k), string(data)))
					default:
						_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("_contains not supported for database (%s)", s.dbType), nil, nil)
					}
				case "$gt":
					array = append(array, goqu.I(k).Gt(v2))

				case "$gte":
					array = append(array, goqu.I(k).Gte(v2))

				case "$lt":
					array = append(array, goqu.I(k).Lt(v2))

				case "$lte":
					array = append(array, goqu.I(k).Lte(v2))

				case "$in":
					array = append(array, goqu.I(k).In(v2))

				case "$nin":
					array = append(array, goqu.I(k).NotIn(v2))
				}
			}
		} else {
			if vString, ok := v.(string); ok && isJoin {
				v = goqu.I(vString)
			}
			array = append(array, goqu.I(k).Eq(v))
		}
	}

	return goqu.And(array...)
}

func (s *SQL) generateWhereClause(ctx context.Context, q *goqu.SelectDataset, find map[string]interface{}, matchWhere []map[string]interface{}) (query *goqu.SelectDataset) {
	query = q

	exps := make([]goqu.Expression, len(matchWhere))
	for i, f := range matchWhere {
		exps[i] = s.generator(ctx, f, false)
	}

	if len(find) > 0 {
		exp := s.generator(ctx, find, false)
		exps = append(exps, exp)
	}

	if len(exps) > 0 {
		query = query.Where(goqu.And(exps...))
	}

	return query
}

func generateRecord(temp interface{}) (goqu.Record, error) {
	insertObj, ok := temp.(map[string]interface{})
	if !ok {
		return nil, errors.New("incorrect insert object provided")
	}

	record := make(goqu.Record, len(insertObj))
	for k, v := range insertObj {
		record[k] = v
	}
	return record, nil
}

func (s *SQL) getColName(col string) string {
	switch model.DBType(s.dbType) {
	case model.Postgres, model.SQLServer:
		return fmt.Sprintf("%s.%s", s.name, col)
	}
	return col
}

func (s *SQL) generateQuerySQLServer(query string) string {
	return strings.Replace(query, "$", "@p", -1)
}

func mysqlTypeCheck(ctx context.Context, dbType model.DBType, types []*sql.ColumnType, mapping map[string]interface{}) {
	var err error
	for _, colType := range types {
		typeName := colType.DatabaseTypeName()
		switch v := mapping[colType.Name()].(type) {
		case string:
			switch typeName {
			case "JSON", "JSONB":
				var val interface{}
				if err := json.Unmarshal([]byte(v), &val); err == nil {
					mapping[colType.Name()] = val
				}
			}
			if dbType == model.SQLServer || typeName == "NVARCHAR" {
				if (strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")) || (strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]")) {
					var val interface{}
					if err := json.Unmarshal([]byte(v), &val); err == nil {
						mapping[colType.Name()] = val
						continue
					}
				}
				mapping[colType.Name()] = v
			}
		case []byte:
			switch typeName {
			case "BIT":
				if len(v) > 0 {
					if v[0] == byte(1) {
						mapping[colType.Name()] = true
					} else {
						mapping[colType.Name()] = false
					}
				}

			case "JSON", "JSONB":
				var val interface{}
				if err := json.Unmarshal(v, &val); err == nil {
					mapping[colType.Name()] = val
				}
			case "VARCHAR", "CHAR", "TEXT", "NAME", "BPCHAR":
				// NOTE: The NAME data type is only valid for Postgres database, as it exists for Postgres only (Name is a 63 byte (varchar) type used for storing system identifiers.)
				val, ok := mapping[colType.Name()].([]byte)
				if ok {
					mapping[colType.Name()] = string(val)
				}
			case "TINYINT":
				mapping[colType.Name()], err = strconv.ParseBool(string(v))
				if err != nil {
					helpers.Logger.LogInfo(helpers.GetRequestID(ctx), fmt.Sprintf("Error:%v", err), nil)
				}
			case "BIGINT", "INT", "SMALLINT":
				mapping[colType.Name()], err = strconv.ParseInt(string(v), 10, 64)
				if err != nil {
					helpers.Logger.LogInfo(helpers.GetRequestID(ctx), fmt.Sprintf("Error:%v", err), nil)
				}
			case "DECIMAL", "NUMERIC", "FLOAT":
				mapping[colType.Name()], err = strconv.ParseFloat(string(v), 64)
				if err != nil {
					helpers.Logger.LogInfo(helpers.GetRequestID(ctx), fmt.Sprintf("Error:%v", err), nil)
				}
			case "DATETIME":
				if dbType == model.MySQL {
					d, _ := time.Parse("2006-01-02 15:04:05", string(v))
					mapping[colType.Name()] = d.Format(time.RFC3339Nano)
					continue
				}
				mapping[colType.Name()] = string(v)
			case "TIMESTAMP", "TIME", "DATE": // For mysql
				mapping[colType.Name()] = string(v)
			}
		case int64:
			if typeName == "TINYINT" {
				// this case occurs for mysql database with column type tinyint during the upsert operation
				if v == int64(1) {
					mapping[colType.Name()] = true
				} else {
					mapping[colType.Name()] = false
				}
			}
		case time.Time:
			switch typeName {
			// For postgres & SQL server
			case "TIME":
				mapping[colType.Name()] = v.Format("15:04:05.999999999")
				continue
			case "DATE":
				mapping[colType.Name()] = v.Format("2006-01-02")
				continue
			}
			mapping[colType.Name()] = v.UTC().Format(time.RFC3339Nano)
		case primitive.DateTime:
			mapping[colType.Name()] = v.Time().UTC().Format(time.RFC3339Nano)
		}
	}
}

func (s *SQL) processJoins(ctx context.Context, query *goqu.SelectDataset, join []*model.JoinOption, sel map[string]int32, isAggregate bool) (*goqu.SelectDataset, error) {
	for _, j := range join {
		on := s.generator(ctx, j.On, true)
		switch j.Type {
		case "", "LEFT":
			query = query.LeftJoin(goqu.T(s.getColName(j.Table)), goqu.On(on))
		case "RIGHT":
			query = query.RightJoin(goqu.T(s.getColName(j.Table)), goqu.On(on))
		case "INNER":
			query = query.InnerJoin(goqu.T(s.getColName(j.Table)), goqu.On(on))
		case "OUTER":
			query = query.FullOuterJoin(goqu.T(s.getColName(j.Table)), goqu.On(on))
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid join type (%s) provided", j.Type), nil, nil)
		}

		// Don't put any fields in select clause if it's an aggregation request
		if !isAggregate {
			isValidJoin, columnName := utils.IsValidJoin(j.On, j.Table)
			if isValidJoin {
				// forcing a join field in select clause, required for caching
				sel[j.Table+"."+columnName] = 1
			}
		}

		if j.Join != nil {
			q, err := s.processJoins(ctx, query, j.Join, sel, isAggregate)
			if err != nil {
				return nil, err
			}

			query = q
		}
	}

	return query, nil
}

// replaceSQLOperationWithPlaceHolder
// e.g-> sql string -> select * from users limit $1
// this function will replace (limit $1) to an value that you specify
func replaceSQLOperationWithPlaceHolder(replace, sqlString string, replaceWith func(value string) string) (string, string) {
	startIndex := strings.Index(sqlString, replace)
	if startIndex == -1 {
		return "", sqlString
	}
	endIndex := 0
	dollarEndIndex := startIndex + len(replace) + 1
	if dollarEndIndex > len(sqlString) {
		return "", sqlString
	}
	tempArr := sqlString[dollarEndIndex:]
	dollarValue := ""
	for index, value := range tempArr {
		dollarValue += string(value)
		if unicode.IsSpace(value) || len(tempArr)-1 == index {
			endIndex = dollarEndIndex + index + 1
			break
		}
	}
	dollarValue = strings.TrimSpace(dollarValue)
	arr1 := sqlString[:startIndex]
	arr2 := sqlString[endIndex:]
	sqlString = arr1 + replaceWith(dollarValue) + " " + arr2
	return dollarValue, strings.TrimSpace(sqlString)
}

func mutateSQLServerLimitAndOffsetOperation(sqlString string, req *model.ReadRequest) (string, error) {
	if req.Options.Skip != nil && req.Options.Limit != nil {
		if len(req.Options.Sort) == 0 {
			return "", fmt.Errorf("sql server cannot process skip operation, sort option is mandatory with skip")
		}
		offsetValue, sqlString := replaceSQLOperationWithPlaceHolder("OFFSET", sqlString, func(value string) string {
			return ""
		})

		_, sqlString = replaceSQLOperationWithPlaceHolder("LIMIT", sqlString, func(value string) string {
			return fmt.Sprintf("OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", offsetValue, value)
		})
		return sqlString, nil
	}
	if req.Options.Limit != nil {
		_, sqlString = replaceSQLOperationWithPlaceHolder("LIMIT", sqlString, func(value string) string {
			return ""
		})

		if strings.HasPrefix(sqlString, "SELECT DISTINCT") {
			sqlString = strings.Replace(sqlString, "SELECT DISTINCT", fmt.Sprintf("SELECT DISTINCT TOP %d", uint(*req.Options.Limit)), 1)
		} else {
			sqlString = strings.Replace(sqlString, "SELECT", fmt.Sprintf("SELECT TOP %d", uint(*req.Options.Limit)), 1)
		}
		return sqlString, nil
	}
	return sqlString, nil
}
