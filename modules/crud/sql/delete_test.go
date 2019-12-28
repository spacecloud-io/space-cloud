package sql

import (
	"testing"

	"github.com/spaceuptech/space-cloud/model"
)

func TestGenerateDeleteQuery(t *testing.T) {
	tests := []struct {
		name, dbType, tableName, wantThis, orThis string
		req                                       model.DeleteRequest
	}{
		{name: "Successfull Test", dbType: "mysql", tableName: "fooTable", orThis: "DELETE FROM projectName.fooTable WHERE ((String1 = ?) AND (String2 = ?))", wantThis: "DELETE FROM projectName.fooTable WHERE ((String1 = ?) AND (String2 = ?))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": "1", "String2": "2"}}},
		{name: "Nested Map Interface Equal To", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 = ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}}},
		{name: "Nested Map Interface Not Equal To", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 != ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}}},
		{name: "Nested Map Interface Greater than ", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 > ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}}},
		{name: "Nested Map Interface Greater than Equal To", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 >= ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}}},
		{name: "Nested Map Interface Less Than", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 < ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}}},
		{name: "Nested Map Interface Less Than Equal To", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 <= ?)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}}},
		{name: "Nested Map Interface In", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 IN (?))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
		{name: "Nested Map Interface Not in", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 NOT IN (?))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
		{name: "Nested Map Interface OR", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE ((string1ofstring1 = ?) OR (string1ofstring2 = ?))", req: model.DeleteRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
		{name: "When length is 0", dbType: "mysql", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable", req: model.DeleteRequest{Find: map[string]interface{}{}}},

		{name: "Successfull Test", dbType: "sqlserver", tableName: "fooTable", orThis: "DELETE FROM projectName.fooTable WHERE ((String2 = @p1) AND (String1 = @p2))", wantThis: "DELETE FROM projectName.fooTable WHERE ((String1 = @p1) AND (String2 = @p2))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": "1", "String2": "2"}}},
		{name: "Nested Map Interface Equal To", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 = @p1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}}},
		{name: "Nested Map Interface Not Equal To", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 != @p1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}}},
		{name: "Nested Map Interface Greater than ", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 > @p1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}}},
		{name: "Nested Map Interface Greater than Equal To", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 >= @p1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}}},
		{name: "Nested Map Interface Less Than", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 < @p1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}}},
		{name: "Nested Map Interface Less Than Equal To", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 <= @p1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}}},
		{name: "Nested Map Interface In", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 IN (@p1))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
		{name: "Nested Map Interface Not in", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 NOT IN (@p1))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
		{name: "Nested Map Interface OR", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE ((string1ofstring1 = @p1) OR (string1ofstring2 = @p2))", req: model.DeleteRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
		{name: "When length is 0", dbType: "sqlserver", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable", req: model.DeleteRequest{Find: map[string]interface{}{}}},

		{name: "Successfull Test", dbType: "postgres", tableName: "fooTable", orThis: "DELETE FROM projectName.fooTable WHERE ((String2 = $1) AND (String1 = $2))", wantThis: "DELETE FROM projectName.fooTable WHERE ((String1 = $1) AND (String2 = $2))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": "1", "String2": "2"}}},
		{name: "Nested Map Interface Equal To", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 = $1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}}},
		{name: "Nested Map Interface Not Equal To", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 != $1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}}},
		{name: "Nested Map Interface Greater than ", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 > $1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}}},
		{name: "Nested Map Interface Greater than Equal To", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 >= $1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}}},
		{name: "Nested Map Interface Less Than", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 < $1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}}},
		{name: "Nested Map Interface Less Than Equal To", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 <= $1)", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}}},
		{name: "Nested Map Interface In", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 IN ($1))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
		{name: "Nested Map Interface Not in", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE (String1 NOT IN ($1))", req: model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
		{name: "Nested Map Interface OR", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable WHERE ((string1ofstring1 = $1) OR (string1ofstring2 = $2))", req: model.DeleteRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
		{name: "When length is 0", dbType: "postgres", tableName: "fooTable", wantThis: "DELETE FROM projectName.fooTable", req: model.DeleteRequest{Find: map[string]interface{}{}}},
	}
	project := "projectName"

	for i, test := range tests {
		s := &SQL{dbType: test.dbType}
		t.Run(test.name, func(t *testing.T) {
			sqlString, _, err := s.generateDeleteQuery(project, test.tableName, &test.req)
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
