package sql

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/model"

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
		s := SQL{dbType: string(dbTypeValue)}
		for i, structValue := range tddStruct {
			sqlQuery, _, err := s.generateCreateQuery(ctx, structValue.project, structValue.col, &structValue.req)

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
