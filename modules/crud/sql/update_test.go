package sql

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func TestGenerateUpdateQuery(t *testing.T) {
	truecases := 10
	project := "projectName"
	tests := []struct {
		name, tableName, wantThis, orThis string
		req                               model.UpdateRequest
	}{
		{name: "Successfull Test", tableName: "fooTable", orThis: "UPDATE fooTable SET String1=? WHERE ((FindString2 = ?) AND (FindString1 = ?))", wantThis: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
		{name: "where clause is 0", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=?", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}}},
		{name: "Successfull Test EQUAL", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 = ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$eq": 1}}}},
		{name: "Successfull Test NOT EQUAL", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 != ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$ne": 1}}}},
		{name: "Successfull Test GREATER THAN", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 > ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$gt": 1}}}},
		{name: "Successfull Test GREATER THAN EQUAL TO", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 >= ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$gte": 1}}}},
		{name: "Successfull Test LESS THAN", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 < ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$lt": 1}}}},
		{name: "Successfull Test LESS THAN EQUAL TO", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 <= ?)", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$lte": 1}}}},
		{name: "Successfull Test IN", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 IN (?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$in": 1}}}},
		{name: "Successfull Test NOT IN", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE (key1 NOT IN (?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}}, Find: map[string]interface{}{"key1": map[string]interface{}{"$nin": 1}}}},

		{name: "Error Update is NIL", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
		{name: "Error No $set", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Update: map[string]interface{}{}, Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
		{name: "Query Update sql", tableName: "fooTable", wantThis: "UPDATE fooTable SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))", req: model.UpdateRequest{Update: map[string]interface{}{"$set": map[string]interface{}{}}, Find: map[string]interface{}{"FindString1": "1", "FindString2": "2"}}},
	}
	s := SQL{dbType: string(utils.MySQL)}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sqlString, _, err := s.generateUpdateQuery(context.TODO(), project, test.tableName, &test.req)
			if i < truecases {
				if i == 0 {
					if ((sqlString != test.wantThis) && (sqlString != test.orThis)) || err != nil {
						t.Errorf("|Got| %s |But wantThis| %s |But orThis| %s \n %v", sqlString, test.wantThis, test.orThis, err)
					}
				} else if (sqlString != test.wantThis) || err != nil {
					t.Errorf("|Got| %s |But wantThis| %s \n Error %v", sqlString, test.wantThis, err)
				}
			} else if (sqlString == test.wantThis) || err == nil {
				t.Errorf("|Got| %s |But wantThis| %s \n Error %v", sqlString, test.wantThis, err)
			}
		})
	}
}
