package sql

import (
	"context"
	"testing"

	"github.com/spaceuptech/space-cloud/model"
)

func TestGenerateDeleteQuery(t *testing.T) {
	tests := []struct {
		name, tableName, wantThis, orThis string
		req                               model.DeleteRequest
	}{
		{name: "Successfull Test", tableName: "fooTable", orThis: "DELETE FROM fooTable WHERE ((String2 = ?) AND (String1 = ?))", wantThis: "DELETE FROM fooTable WHERE ((String1 = ?) AND (String2 = ?))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": "1", "String2": "2"}}},
		{name: "Nested Map Interface Equal To", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 = ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}}},
		{name: "Nested Map Interface Not Equal To", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 != ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}}},
		{name: "Nested Map Interface Greater than ", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 > ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}}},
		{name: "Nested Map Interface Greater than Equal To", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 >= ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}}},
		{name: "Nested Map Interface Less Than", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 < ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}}},
		{name: "Nested Map Interface Less Than Equal To", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 <= ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}}},
		{name: "Nested Map Interface In", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 IN (?))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
		{name: "Nested Map Interface Not in", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE (String1 NOT IN (?))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
		{name: "Nested Map Interface OR", tableName: "fooTable", wantThis: "DELETE FROM fooTable WHERE ((string1ofstring1 = ?) OR (string1ofstring2 = ?))", req: model.DeleteRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
		{name: "When length is 0", tableName: "fooTable", wantThis: "DELETE FROM fooTable", req: model.DeleteRequest{Find: map[string]interface{}{}}},
	}
	var ctx context.Context
	project := "projectName"
	s := SQL{dbType: string("sql-mysql")}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sqlString, _, err := s.generateDeleteQuery(ctx, project, test.tableName, &test.req)
			if i == 0 {
				if ((sqlString != test.wantThis) && (sqlString != test.orThis)) || err != nil {
					t.Errorf("|Got| %s |But wantThis| %s |But orThis| %s \n %v", sqlString, test.wantThis, test.orThis, err)
				}
			} else if (sqlString != test.wantThis) || err != nil {
				t.Errorf("|Got| %s |But wantThis| %s \n %v", sqlString, test.wantThis, err)
			}
		})
	}

}
