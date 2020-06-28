package sql

import (
	"reflect"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/doug-martin/goqu/v8/dialect/postgres"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestSQL_generateReadQuery(t *testing.T) {
	// temp := "one"
	type fields struct {
		enabled    bool
		connection string
		client     *sqlx.DB
		dbType     string
	}
	type args struct {
		project string
		col     string
		req     *model.ReadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		want1   []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "String1 = ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 = ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 != ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 != ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 > ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 > ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 >= ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 >= ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 < ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 < ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 <= ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 <= ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 in ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 IN (?))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 not in ?",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
			want:    []string{"SELECT * FROM table WHERE (String1 NOT IN (?))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "string1 or string2",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
			want:    []string{"SELECT * FROM table WHERE ((string1ofstring1 = ?) OR (string1ofstring2 = ?))"},
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:    "regex",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"fieldName": map[string]interface{}{"$regex": "ss"}}}},
			want:    []string{"SELECT * FROM table WHERE (fieldName REGEXP ?)"},
			want1:   []interface{}{"ss"},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from one doc",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "one"}},
			want:    []string{"SELECT Column1, Column2 FROM table WHERE (Column1 = ?)", "SELECT Column2, Column1 FROM table WHERE (Column1 = ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from all doc",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM table WHERE (Column1 = ?)", "SELECT Column2, Column1 FROM table WHERE (Column1 = ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Column1 = ?, Limit = ? and offset = ?",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Skip: iti(2), Limit: iti(10)},
					Operation: "all"}},
			want:    []string{"SELECT * FROM table WHERE (Column1 = ?) LIMIT ? OFFSET ?"},
			want1:   []interface{}{int64(1), int64(10), int64(2)},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from all doc",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM table WHERE (Column1 = ?)", "SELECT Column2, Column1 FROM table WHERE (Column1 = ?)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Select column in asc and desc",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}, Sort: []string{"Column1", "-Column2"}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM table ORDER BY Column1 ASC, Column2 DESC", "SELECT Column2, Column1 FROM table ORDER BY Column1 ASC, Column2 DESC"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Count",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "count"}},
			want:    []string{"SELECT COUNT(*) FROM table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Select Distinct",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}, Distinct: str("Column1")},
					Operation: "distinct"}},
			want:    []string{"SELECT DISTINCT Column1 FROM table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:    "Select JSON",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"Obj1": map[string]interface{}{"$contains": map[string]interface{}{"obj1": "value1"}}}}},
			want:    []string{"SELECT * FROM table WHERE json_contains(Obj1,?)"},
			want1:   []interface{}{`{"obj1":"value1"}`},
			wantErr: false,
		},
		// postgres
		{
			name:    "String1 = ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 = $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 != ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 != $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 > ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 > $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 >= ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 >= $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 < ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 < $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 <= ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 <= $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 in ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 IN ($1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 not in ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 NOT IN ($1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "string1 or string2",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
			want:    []string{"SELECT * FROM test.table WHERE ((string1ofstring1 = $1) OR (string1ofstring2 = $2))"},
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:    "String1 in ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 IN ($1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 not in ?",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 NOT IN ($1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "string1 or string2",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
			want:    []string{"SELECT * FROM test.table WHERE ((string1ofstring1 = $1) OR (string1ofstring2 = $2))"},
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:    "regex",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"fieldName": map[string]interface{}{"$regex": "ss"}}}},
			want:    []string{"SELECT * FROM test.table WHERE (fieldName ~ $1)"},
			want1:   []interface{}{"ss"},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from one doc",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "one"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table WHERE (Column1 = $1)", "SELECT Column2, Column1 FROM test.table WHERE (Column1 = $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from all doc",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table WHERE (Column1 = $1)", "SELECT Column2, Column1 FROM test.table WHERE (Column1 = $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Column1 = ?, Limit = ? and offset = ?",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Skip: iti(2), Limit: iti(10)},
					Operation: "all"}},
			want:    []string{"SELECT * FROM test.table WHERE (Column1 = $1) LIMIT $2 OFFSET $3"},
			want1:   []interface{}{int64(1), int64(10), int64(2)},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from all doc",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table WHERE (Column1 = $1)", "SELECT Column2, Column1 FROM test.table WHERE (Column1 = $1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "Select JSON",
			fields:  fields{dbType: "postgres"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"Obj1": map[string]interface{}{"$contains": map[string]interface{}{"obj1": "value1"}}}}},
			want:    []string{"SELECT * FROM test.table WHERE Obj1 @> $1"},
			want1:   []interface{}{`{"obj1":"value1"}`},
			wantErr: false,
		},

		// sqlserver
		{
			name:    "String1 = ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 = @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 != ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 != @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 > ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 > @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 >= ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 >= @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 < ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 < @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 <= ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 <= @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 in ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 IN (@p1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 not in ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 NOT IN (@p1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "string1 or string2",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
			want:    []string{"SELECT * FROM test.table WHERE ((string1ofstring1 = @p1) OR (string1ofstring2 = @p2))"},
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:    "String1 in ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 IN (@p1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "String1 not in ?",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}}},
			want:    []string{"SELECT * FROM test.table WHERE (String1 NOT IN (@p1))"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:    "string1 or string2",
			fields:  fields{dbType: "sqlserver"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}}},
			want:    []string{"SELECT * FROM test.table WHERE ((string1ofstring1 = @p1) OR (string1ofstring2 = @p2))"},
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from one doc",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "one"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table WHERE (Column1 = @p1)", "SELECT Column2, Column1 FROM test.table WHERE (Column1 = @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from all doc",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table WHERE (Column1 = @p1)", "SELECT Column2, Column1 FROM test.table WHERE (Column1 = @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Column1 = ?, Limit = ? and offset = ?",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Skip: iti(2), Limit: iti(10)},
					Operation: "all"}},
			want:    []string{"SELECT * FROM test.table WHERE (Column1 = @p1) LIMIT @p2 OFFSET @p3"},
			want1:   []interface{}{int64(1), int64(10), int64(2)},
			wantErr: false,
		},
		{
			name:   "Column1 = ? and select Column1 from all doc",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table WHERE (Column1 = @p1)", "SELECT Column2, Column1 FROM test.table WHERE (Column1 = @p1)"},
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQL{
				enabled:    tt.fields.enabled,
				connection: tt.fields.connection,
				client:     tt.fields.client,
				dbType:     tt.fields.dbType,
				name:       tt.args.project,
			}
			got, got1, err := s.generateReadQuery(tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQL.generateReadQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if found := find(tt.want, got); !found {
				t.Errorf("SQL.generateReadQuery() got = %v, want %v", got, tt.want[0])
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SQL.generateReadQuery() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func str(s string) *string {
	return &s
}
func iti(s int64) *int64 {
	return &s
}
