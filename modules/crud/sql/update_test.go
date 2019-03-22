package sql

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/model"
)

func TestGenerateUpdateQuery(t *testing.T) {
	truecases := 10
	var ctx context.Context
	project := "projectName"
	tests := []struct {
		name, tableName, want string
		req                   model.UpdateRequest
	}{
		{name: "Successfull Test", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
		{name: "where clause is 0", tableName: "fooTable", want: "UPDATE fooTable SET String1=?", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}}},
		{name: "Successfull Test EQUAL", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 = ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$eq": 1}}}},
		{name: "Successfull Test NOT EQUAL", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 != ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$ne": 1}}}},
		{name: "Successfull Test GREATER THAN", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 > ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$gt": 1}}}},
		{name: "Successfull Test GREATER THAN EQUAL TO", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 >= ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$gte": 1}}}},
		{name: "Successfull Test LESS THAN", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 < ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$lt": 1}}}},
		{name: "Successfull Test LESS THAN EQUAL TO", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 <= ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$lte": 1}}}},
		{name: "Successfull Test IN", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 IN (?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$in": 1}}}},
		{name: "Successfull Test NOT IN", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE (key1 NOT IN (?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$nin": 1}}}},

		//
		{name: "Error Update is NIL", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
		{name: "Error No $set", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Update: map[string]interface{}{}, Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
		{name: "Query Update sql", tableName: "fooTable", want: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{}}, Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
	}
	s, _ := InitializeDatabase("sql-mysql")
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sqlString, _, err := s.GenerateUpdateQuery(ctx, project, test.tableName, &test.req)
			if i < truecases {
				if (sqlString != test.want) || err != nil {
					t.Errorf("|Got| %s |But Want| %s \n Error %v", sqlString, test.want, err)
				}
			} else if (sqlString == test.want) || err == nil {
				t.Errorf("|Got| %s |But Want| %s \n Error %v", sqlString, test.want, err)
			}
		})
	}
}
