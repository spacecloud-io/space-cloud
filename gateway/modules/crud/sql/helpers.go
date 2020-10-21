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

	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/doug-martin/goqu/v8"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (s *SQL) generator(ctx context.Context, find map[string]interface{}, isJoin bool) (goqu.Expression, []string) {
	var regxarr []string
	array := []goqu.Expression{}
	for k, v := range find {
		if k == "$or" {
			orArray := v.([]interface{})
			orFinalArray := []goqu.Expression{}
			for _, item := range orArray {
				exp, a := s.generator(ctx, item.(map[string]interface{}), isJoin)
				orFinalArray = append(orFinalArray, exp)
				regxarr = append(regxarr, a...)
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
						regxarr = append(regxarr, fmt.Sprintf("%s = $", k))
					case "mysql":
						regxarr = append(regxarr, fmt.Sprintf("%s = ?", k))
					}
					array = append(array, goqu.I(k).Eq(v2))
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
	return goqu.And(array...), regxarr
}

func (s *SQL) generateWhereClause(ctx context.Context, q *goqu.SelectDataset, find map[string]interface{}) (query *goqu.SelectDataset, arr []string) {
	query = q
	if len(find) == 0 {
		return
	}
	exp, arr := s.generator(ctx, find, false)
	query = query.Where(exp)
	return query, arr
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

func (s *SQL) getDBName(col string) string {
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
		case []byte:
			switch typeName {
			case "VARCHAR", "TEXT", "JSON", "JSONB":
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
			case "DATE", "DATETIME":
				if dbType == model.MySQL {
					d, _ := time.Parse("2006-01-02 15:04:05", string(v))
					mapping[colType.Name()] = d.Format(time.RFC3339)
					continue
				}
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
			mapping[colType.Name()] = v.UTC().Format(time.RFC3339)

		case primitive.DateTime:
			mapping[colType.Name()] = v.Time().UTC().Format(time.RFC3339)
		}
	}
}

func (s *SQL) processJoins(ctx context.Context, table string, query *goqu.SelectDataset, join []model.JoinOption) (*goqu.SelectDataset, error) {
	for _, j := range join {
		on, _ := s.generator(ctx, j.On, true)
		switch j.Type {
		case "", "LEFT":
			query = query.LeftJoin(goqu.T(j.Table), goqu.On(on))
		case "RIGHT":
			query = query.RightJoin(goqu.T(j.Table), goqu.On(on))
		case "INNER":
			query = query.InnerJoin(goqu.T(j.Table), goqu.On(on))
		case "OUTER":
			query = query.FullOuterJoin(goqu.T(j.Table), goqu.On(on))
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid join type (%s) provided", j.Type), nil, nil)
		}

		if j.Join != nil {
			q, err := s.processJoins(ctx, j.Table, query, j.Join)
			if err != nil {
				return nil, err
			}

			query = q
		}
	}

	return query, nil
}
