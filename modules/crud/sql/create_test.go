package sql

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/spaceuptech/space-cloud/model"

	"github.com/jmoiron/sqlx"
	"github.com/spaceuptech/space-cloud/utils"
)

var tddStruct = []struct {
	project, col, want string
	req                model.CreateRequest
}{
	{project: "foo", col: "footable1", want: "INSERT INTO footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", col: "footable1", want: "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},
	{project: "foo", col: "footable1", want: "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", col: "footable1", want: "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{1, 2, 3}}},
	{project: "foo", col: "footable1", want: "INSERT INTO footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{}}},
	{project: "foo", col: "footable1", want: "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{}, map[string]interface{}{"string1": "1"}}}},
	{project: "foo", col: "footable1", want: "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string6": "2"}}}},
	{project: "foo", col: "footable2", want: "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},
}

func TestGenerateCreateQuery(t *testing.T) {
	truecases := 2
	var dbTypes = []utils.DBType{"sql-mysql"}
	var ctx context.Context

	for _, dbTypeValue := range dbTypes {
		s, err := InitializeDatabase(dbTypeValue)
		if err != nil {
			fmt.Println("initialization ", err)
			return
		}
		for i, structValue := range tddStruct {
			sqlQuery, _, err := s.GenerateCreateQuery(ctx, structValue.project, structValue.col, &structValue.req)

			if i < truecases {
				if (sqlQuery != structValue.want) || err != nil {
					t.Error(i+1, "incorrect match1", "got : ", sqlQuery, "want: ", structValue.want, "\n Err : ", err)
				}
				continue
			} else if (sqlQuery == structValue.want) && err == nil {
				t.Error(i+1, "incorrect match2", "got : ", sqlQuery, "want: ", structValue.want, "\n Err : ", err)
			}
		}
	}

}

func InitializeDatabase(dbType utils.DBType) (*SQL, error) {
	var sql *sqlx.DB
	var err error
	s := &SQL{}
	switch dbType {
	case utils.Postgres:
		sql, err = sqlx.Open("postgres", "postgres://myuser:password@localhost/testdb?sslmode=disable")
		s.dbType = "postgres"

	case utils.MySQL:
		sql, err = sqlx.Open("mysql", "testuser:password@(localhost:3306)/testdb")
		s.dbType = "mysql"

	default:
		return nil, errors.New("SQL: Invalid driver provided")
	}

	if err != nil {
		return nil, err
	}

	err = sql.Ping()
	if err != nil {
		return nil, err
	}

	s.client = sql
	return s, nil
}
