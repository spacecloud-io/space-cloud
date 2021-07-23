package sql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v8"
	"github.com/spaceuptech/helpers"

	_ "github.com/denisenkom/go-mssqldb"                // Import for MsSQL
	_ "github.com/doug-martin/goqu/v8/dialect/postgres" // Dialect for postgres
	_ "github.com/go-sql-driver/mysql"                  // Import for MySQL
	_ "github.com/lib/pq"                               // Import for

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Update updates the document(s) which match the condition provided.
func (s *SQL) Update(ctx context.Context, col string, req *model.UpdateRequest) (int64, error) {
	tx, err := s.getClient().BeginTxx(ctx, nil) // TODO - Write *sqlx.TxOption instead of nil
	if err != nil {
		return 0, err
	}
	count, err := s.update(ctx, col, req, tx)
	if err != nil {
		return 0, err
	}
	return count, tx.Commit() // commit the Batch
}

func (s *SQL) update(ctx context.Context, col string, req *model.UpdateRequest, executor executor) (int64, error) {
	if req == nil {
		return 0, utils.ErrInvalidParams
	}
	if req.Update == nil {
		return 0, utils.ErrInvalidParams
	}
	switch req.Operation {
	case utils.All:
		var count int64
		for k := range req.Update {
			switch k {
			case "$set", "$inc", "$mul", "$max", "$min", "$currentDate":
				sqlQuery, args, err := s.generateUpdateQuery(ctx, col, req, k)
				if err != nil {
					return 0, err
				}
				helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Update Query", map[string]interface{}{"sqlQuery": sqlQuery, "queryArgs": args})
				res, err := doExecContext(ctx, sqlQuery, args, executor)
				if err != nil {
					return 0, err
				}

				c, _ := res.RowsAffected()
				count += c

			default: // (case "$push", "$unset", "$rename")
				return 0, utils.ErrInvalidParams
			}
		}

		return count, nil

	case utils.Upsert:
		count, _, _, _, err := s.read(ctx, col, &model.ReadRequest{Find: req.Find, Operation: utils.All}, executor)
		if err != nil {
			return 0, err
		}

		if count == 0 {
			// not found. So, insert
			doc := make(map[string]interface{})
			dates := make(map[string]interface{})
			for k, v := range req.Find {
				obj, ok := v.(map[string]interface{})
				if !ok {
					doc[k] = v
					continue
				}

				// Required when v is a map. Eg. {"id": { "$eq": "some-value" }}
				for _, colValue := range obj {
					doc[k] = colValue
				}
			}
			for op := range req.Update {
				m, ok := req.Update[op].(map[string]interface{})
				if !ok {
					return 0, utils.ErrInvalidParams
				}
				if op == "$currentDate" {
					err := s.flattenForDate(ctx, &m)
					if err != nil {
						return 0, err
					}
					for k, v := range m { // k -> column name, v -> function name
						dates[k] = v
					}
				} else {
					for k, v := range m { // k -> column name
						doc[k] = v
					}
				}
			}
			sqlQuery, args, err := s.generateCreateQuery(col, &model.CreateRequest{Document: doc, Operation: utils.One})
			if err != nil {
				return 0, err
			}

			for k, v := range dates {
				f := strings.Index(sqlQuery, ")")
				sqlQuery = sqlQuery[:f] + ", " + k + sqlQuery[f:]
				l := strings.LastIndex(sqlQuery, ")")
				sqlQuery = sqlQuery[:l] + ", " + v.(string) + sqlQuery[l:]
			}

			res, err := doExecContext(ctx, sqlQuery, args, executor)
			if err != nil {
				return 0, err
			}

			return res.RowsAffected()
		}
		req.Operation = utils.All
		return s.update(ctx, col, req, executor)
	default: // (case utils.One)
		return 0, errors.New("sql databases do no support op type `one` for updates")
	}
}

// generateUpdateQuery makes query for update operations
func (s *SQL) generateUpdateQuery(ctx context.Context, col string, req *model.UpdateRequest, op string) (string, []interface{}, error) {
	// Generate a prepared query builder

	dbType := s.dbType
	if dbType == string(model.SQLServer) {
		dbType = string(model.Postgres)
	}
	dialect := goqu.Dialect(dbType)
	query := dialect.From(s.getColName(col)).Prepared(true)

	if req.Find != nil {
		// Get the where clause from query object
		query = s.generateWhereClause(ctx, query, req.Find, nil)
	}

	if req.Update == nil {
		return "", nil, utils.ErrInvalidParams
	}
	m, ok := req.Update[op].(map[string]interface{})
	if !ok {
		return "", nil, utils.ErrInvalidParams
	}

	if op == "$currentDate" {
		err := s.flattenForDate(ctx, &m)
		if err != nil {
			return "", nil, err
		}
	}

	record, err := generateRecord(req.Update[op])
	if err != nil {
		return "", nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error generating update query unable to generate record %s", op), err, nil)
	}

	// Generate SQL string and arguments
	sqlString, args, err := query.Update().Set(record).ToSQL()
	if err != nil {
		return "", nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Error generating update query unable generate sql string", err, nil)
	}

	sqlString = strings.Replace(sqlString, "\"", "", -1)
	switch op {
	case "$set":
		for k := range m {
			if s.dbType == string(model.Postgres) {
				arr := strings.Split(k, ".")
				if len(arr) >= 2 {
					sqlString = strings.Replace(sqlString, k+"=$", fmt.Sprintf("%s=jsonb_set(%s, '{%s}', $", arr[0], arr[0], strings.Join(arr[1:], ",")), -1)
					sqlString = s.sanitiseUpdateQuery(sqlString)
				}
			}
		}
	case "$inc":
		for k, v := range m {
			_, err := checkIfNum(v)
			if err != nil {
				return "", nil, err
			}
			if s.dbType == string(model.MySQL) {
				sqlString = strings.Replace(sqlString, k+"=?", k+"="+k+"+?", -1)
			}
			if dbType == string(model.Postgres) {
				sqlString = strings.Replace(sqlString, k+"=$", k+"="+k+"+$", -1)
			}

		}

	case "$mul":
		for k, v := range m {
			_, err := checkIfNum(v)
			if err != nil {
				return "", nil, err
			}
			if dbType == string(model.MySQL) {
				sqlString = strings.Replace(sqlString, k+"=?", k+"="+k+"*?", -1)
			}
			if dbType == string(model.Postgres) {
				sqlString = strings.Replace(sqlString, k+"=$", k+"="+k+"*$", -1)
			}

		}

	case "$max":
		for k, v := range m {
			_, err := checkIfNum(v)
			if err != nil {
				return "", nil, err
			}
			if s.dbType == string(model.MySQL) {
				sqlString = strings.Replace(sqlString, k+"=?", k+"=GREATEST("+k+","+"?"+")", -1)
			}
			if dbType == string(model.Postgres) {
				sqlString = strings.Replace(sqlString, k+"=$", k+"=GREATEST("+k+","+"$"+"", -1)
			}
		}
		sqlString = s.sanitiseUpdateQuery(sqlString)

	case "$min":
		for k, v := range m {
			_, err := checkIfNum(v)
			if err != nil {
				return "", nil, err
			}
			if dbType == string(model.MySQL) {
				sqlString = strings.Replace(sqlString, k+"=?", k+"=LEAST("+k+","+"?"+")", -1)
			}
			if dbType == string(model.Postgres) {
				sqlString = strings.Replace(sqlString, k+"=$", k+"=LEAST("+k+","+"$", -1)
			}
		}
		sqlString = s.sanitiseUpdateQuery(sqlString)

	case "$currentDate":
		for k, v := range m {
			val, ok := v.(string)
			if !ok {
				return "", nil, utils.ErrInvalidParams
			}
			if dbType == string(model.MySQL) {
				sqlString = strings.Replace(sqlString, k+"=?", k+"="+val, -1)
			}
			if dbType == string(model.Postgres) {
				sqlString = strings.Replace(sqlString, k+"=$", k+"="+val, -1)
			}

			args = args[1:]

		}
		sqlString = s.sanitiseUpdateQuery2(sqlString)
	default:
		return "", nil, utils.ErrInvalidParams
	}

	if s.dbType == string(model.SQLServer) {
		sqlString = s.generateQuerySQLServer(sqlString)
	}

	return sqlString, args, nil
}

func checkIfNum(v interface{}) (string, error) {
	switch val := v.(type) {
	case float64:

		return strconv.FormatFloat(val, 'f', -1, 64), nil

	case int64:
		return strconv.FormatInt(val, 10), nil
	case int:
		return strconv.Itoa(val), nil
	}

	return "", errors.New("invalid data format provided")
}

func (s *SQL) flattenForDate(ctx context.Context, m *map[string]interface{}) error {
	for k, v := range *m {
		mm, ok := v.(map[string]interface{})
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid current date format (%v) provided", reflect.TypeOf(v)), nil, nil)
		}
		for _, valTemp := range mm {
			val, ok := valTemp.(string)
			if !ok {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid current date type (%v) provided", valTemp), nil, nil)
			}
			switch val {
			case "date":
				(*m)[k] = "CURRENT_DATE"
				// CURRENT_DATE is not supported in sql-server
				if model.DBType(s.dbType) == model.SQLServer {
					(*m)[k] = "CAST( GETDATE() AS date )"
				}
			case "timestamp":
				(*m)[k] = "CURRENT_TIMESTAMP"
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid current date value (%s) provided", val), nil, nil)
			}
		}
	}
	return nil
}

func (s *SQL) sanitiseUpdateQuery(sqlString string) string {
	var placeholder byte
	if (model.DBType(s.dbType) == model.Postgres) || (model.DBType(s.dbType) == model.SQLServer) {
		placeholder = '$'
	}
	var start bool
	for i := 0; i < len(sqlString); i++ {
		c := sqlString[i]
		if c == placeholder {
			start = true
		}
		if start && (c == ' ' || c == ',') {
			sqlString = sqlString[0:i] + ")" + sqlString[i:]
			start = false
		}
		if strings.HasPrefix(sqlString[i:], "WHERE") {
			break
		}
	}
	return sqlString
}
func (s *SQL) sanitiseUpdateQuery2(sqlString string) string {

	// counts number of occurrences of date and time stamps
	noOfStamps := 0

	for i := 0; i < len(sqlString); i++ {

		if strings.HasPrefix(sqlString[i:], "CURRENT_TIMESTAMP") {
			i += len("CURRENT_TIMESTAMP")
			if sqlString[i] != ' ' && sqlString[i] != ',' {
				sqlString = sqlString[:i] + sqlString[i+1:]
			}
			noOfStamps++
		}
		if strings.HasPrefix(sqlString[i:], "CURRENT_DATE") {
			i += len("CURRENT_DATE")
			if sqlString[i] != ' ' && sqlString[i] != ',' {
				sqlString = sqlString[:i] + sqlString[i+1:]
			}
			noOfStamps++
		}
		if strings.HasPrefix(sqlString[i:], "CAST( GETDATE() AS date )") {
			i += len("CAST( GETDATE() AS date )")
			if sqlString[i] != ' ' && sqlString[i] != ',' {
				sqlString = sqlString[:i] + sqlString[i+1:]
			}
			noOfStamps++
		}

	}

	// reduces the parameter $ and @p value by no. of stamps occurrence
	if model.DBType(s.dbType) == model.Postgres || model.DBType(s.dbType) == model.SQLServer {
		for i := 1; i < len(sqlString); i++ {
			c, _ := strconv.Atoi(string(sqlString[i]))
			if c > 1 && (sqlString[i-1] == '$' || (sqlString[i-1] == 'p' && sqlString[i-2] == '@')) {
				sqlString = sqlString[:i] + strconv.Itoa(c-noOfStamps) + sqlString[i+1:]
			}
		}
	}

	return sqlString
}
