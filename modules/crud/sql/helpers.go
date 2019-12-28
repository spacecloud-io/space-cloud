package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	goqu "github.com/doug-martin/goqu/v8"

	"github.com/spaceuptech/space-cloud/utils"
)

func (s *SQL) generator(find map[string]interface{}) (goqu.Expression, []string) {
	var regxarr []string
	array := []goqu.Expression{}
	for k, v := range find {
		if k == "$or" {
			orArray := v.([]interface{})
			orFinalArray := []goqu.Expression{}
			for _, item := range orArray {
				exp, a := s.generator(item.(map[string]interface{}))
				orFinalArray = append(orFinalArray, exp)
				regxarr = append(regxarr, a...)
			}

			array = append(array, goqu.Or(orFinalArray...))
			continue
		}
		val, isObj := v.(map[string]interface{})
		if isObj {
			for k2, v2 := range val {
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
			array = append(array, goqu.I(k).Eq(v))
		}
	}
	return goqu.And(array...), regxarr
}

func (s *SQL) generateWhereClause(q *goqu.SelectDataset, find map[string]interface{}) (query *goqu.SelectDataset, arr []string) {
	query = q
	if len(find) == 0 {
		return
	}
	exp, arr := s.generator(find)
	query = query.Where(exp)
	return query, arr
}

func generateRecord(temp interface{}) (goqu.Record, error) {
	insertObj, ok := temp.(map[string]interface{})
	if !ok {
		return nil, errors.New("Incorrect insert object provided")
	}

	record := make(goqu.Record, len(insertObj))
	for k, v := range insertObj {
		record[k] = v
	}
	return record, nil
}

func (s *SQL) getDBName(project, col string) string {
	if s.removeProjectScope {
		return col
	}

	return project + "." + col
}

func (s *SQL) generateQuerySQLServer(query string) string {
	return strings.Replace(query, "$", "@p", -1)
}

func mysqlTypeCheck(dbType utils.DBType, types []*sql.ColumnType, mapping map[string]interface{}) {
	for _, colType := range types {
		typeName := colType.DatabaseTypeName()
		if typeName == "VARCHAR" || typeName == "TEXT" {
			val, ok := mapping[colType.Name()].([]byte)
			if ok {
				mapping[colType.Name()] = string(val)
			}
		}
		if data, ok := mapping[colType.Name()].([]byte); ok && (typeName == "BIGINT" || typeName == "INT" || typeName == "SMALLINT") {
			var err error
			mapping[colType.Name()], err = strconv.ParseInt(string(data), 10, 64)
			if err != nil {
				log.Println("Error:", err)
			}
		}
		if data, ok := mapping[colType.Name()].([]byte); ok && (typeName == "DECIMAL" || typeName == "NUMERIC") {
			var err error
			mapping[colType.Name()], err = strconv.ParseFloat(string(data), 64)
			if err != nil {
				log.Println("Error:", err)
			}
		}
		if data, ok := mapping[colType.Name()].([]byte); ok && (typeName == "DATE" || typeName == "DATETIME") {
			if dbType == utils.MySQL {
				d, _ := time.Parse("2006-01-02 15:04:05", string(data))
				mapping[colType.Name()] = d.Format(time.RFC3339)
				continue
			}
			mapping[colType.Name()] = string(data)
		}
	}
}
