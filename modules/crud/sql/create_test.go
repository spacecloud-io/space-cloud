package sql

import (
	"testing"

	"github.com/spaceuptech/space-cloud/model"
)

var tddStruct = []struct {
	project, dbType, col, want string
	req                        model.CreateRequest
}{
	{project: "foo", dbType: "mysql", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", dbType: "mysql", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},
	{project: "foo", dbType: "mysql", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", dbType: "mysql", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{1, 2, 3}}},
	{project: "foo", dbType: "mysql", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{}}},
	{project: "foo", dbType: "mysql", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{}, map[string]interface{}{"string1": "1"}}}},
	{project: "foo", dbType: "mysql", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string6": "2"}}}},
	{project: "foo", dbType: "mysql", col: "footable2", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},

	{project: "foo", dbType: "sqlserver", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", dbType: "sqlserver", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},
	{project: "foo", dbType: "sqlserver", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", dbType: "sqlserver", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{1, 2, 3}}},
	{project: "foo", dbType: "sqlserver", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{}}},
	{project: "foo", dbType: "sqlserver", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{}, map[string]interface{}{"string1": "1"}}}},
	{project: "foo", dbType: "sqlserver", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string6": "2"}}}},
	{project: "foo", dbType: "sqlserver", col: "footable2", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},

	{project: "foo", dbType: "postgres", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", dbType: "postgres", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},
	{project: "foo", dbType: "postgres", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}}},
	{project: "foo", dbType: "postgres", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{1, 2, 3}}},
	{project: "foo", dbType: "postgres", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2, string3) VALUES (?, ?, ?)", req: model.CreateRequest{Operation: "one", Document: map[string]interface{}{}}},
	{project: "foo", dbType: "postgres", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{}, map[string]interface{}{"string1": "1"}}}},
	{project: "foo", dbType: "postgres", col: "footable1", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string6": "2"}}}},
	{project: "foo", dbType: "postgres", col: "footable2", want: "INSERT INTO foo.footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)", req: model.CreateRequest{Operation: "all", Document: []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}}},
}

func TestGenerateCreateQuery(t *testing.T) {
	truecases := 2

	for i, structValue := range tddStruct {
		s := &SQL{dbType: structValue.dbType}
		sqlQuery, _, err := s.generateCreateQuery(structValue.project, structValue.col, &structValue.req)

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
