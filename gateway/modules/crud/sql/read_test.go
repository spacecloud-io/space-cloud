package sql

import (
	"context"
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
		// #######################################################################################
		// ###################################  MySQL  ###########################################
		// #######################################################################################
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
			name:    "Select JSON",
			fields:  fields{dbType: "mysql"},
			args:    args{project: "test", col: "table", req: &model.ReadRequest{Find: map[string]interface{}{"Obj1": map[string]interface{}{"$contains": map[string]interface{}{"obj1": "value1"}}}}},
			want:    []string{"SELECT * FROM table WHERE json_contains(Obj1,?)"},
			want1:   []interface{}{`{"obj1":"value1"}`},
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
			name:   "simple join without select",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "t1",
				req: &model.ReadRequest{
					Find: map[string]interface{}{"t1.col1": map[string]interface{}{"$eq": 1}},
					Options: &model.ReadOptions{
						Join: []*model.JoinOption{
							{Table: "t2", Type: "LEFT", On: map[string]interface{}{"t1.col1": "t2.col2"}},
						}},
					Operation: "all"}},
			want:    []string{""},
			want1:   nil,
			wantErr: true,
		},
		// This is a valid test case, but we are commenting it out because,
		// the resultant sql string can have more than 16 combination, all of them are valid
		// {
		// 	name:   "nested join with select",
		// 	fields: fields{dbType: "mysql"},
		// 	args: args{project: "test", col: "t1",
		// 		req: &model.ReadRequest{
		// 			Find: map[string]interface{}{"t1.col1": map[string]interface{}{"$eq": 1}},
		// 			Options: &model.ReadOptions{
		// 				Select: map[string]int32{"t1.col1": 1},
		// 				Join: []*model.JoinOption{
		// 					{Table: "t2", Type: "LEFT", On: map[string]interface{}{"t1.col1": "t2.col2"}, Join: []*model.JoinOption{
		// 						{Table: "t3", Type: "LEFT", On: map[string]interface{}{"t2.col3": map[string]interface{}{"$eq": "t3.col4"}}},
		// 					}},
		// 				},
		// 			},
		// 			Operation: "all"}},
		// 	want:    []string{"SELECT t1.col1 AS t1__col1, t2.col2 AS t2__col2, t3.col4 AS t3__col4 FROM t1 LEFT JOIN t2 ON (t1.col1 = t2.col2) LEFT JOIN t3 ON (t2.col3 = t3.col4) WHERE (t1.col1 = ?)"},
		// 	want1:   []interface{}{int64(1)},
		// 	wantErr: false,
		// },

		// #######################################################################################
		// ###################################  Postgres  ########################################
		// #######################################################################################
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
		{
			name:   "Select column in asc and desc",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}, Sort: []string{"Column1", "-Column2"}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table ORDER BY Column1 ASC, Column2 DESC", "SELECT Column2, Column1 FROM test.table ORDER BY Column1 ASC, Column2 DESC"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Count",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "count"}},
			want:    []string{"SELECT COUNT(*) FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Select Distinct",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}, Distinct: str("Column1")},
					Operation: "distinct"}},
			want:    []string{"SELECT DISTINCT Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},

		// #######################################################################################
		// ###################################  SQLServer  #######################################
		// #######################################################################################
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
					Options:   &model.ReadOptions{Skip: iti(2), Limit: iti(10), Sort: []string{"age"}},
					Operation: "all"}},
			want:    []string{"SELECT * FROM test.table WHERE (Column1 = @p1) ORDER BY age ASC OFFSET @p3 ROWS FETCH NEXT @p2 ROWS ONLY"},
			want1:   []interface{}{int64(1), int64(10), int64(2)},
			wantErr: false,
		},
		{
			name:   "Column1 = ?, limit = ?",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{"Column1": map[string]interface{}{"$eq": 1}},
					Options:   &model.ReadOptions{Limit: iti(20)},
					Operation: "all"}},
			want:    []string{"SELECT TOP 20 * FROM test.table WHERE (Column1 = @p1)"},
			want1:   []interface{}{int64(1), int64(20)},
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
			name:   "Select column in asc and desc",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}, Sort: []string{"Column1", "-Column2"}},
					Operation: "all"}},
			want:    []string{"SELECT Column1, Column2 FROM test.table ORDER BY Column1 ASC, Column2 DESC", "SELECT Column2, Column1 FROM test.table ORDER BY Column1 ASC, Column2 DESC"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Count",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}},
					Operation: "count"}},
			want:    []string{"SELECT COUNT(*) FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Select Distinct",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1, "Column2": 1}, Distinct: str("Column1")},
					Operation: "distinct"}},
			want:    []string{"SELECT DISTINCT Column1 FROM test.table"},
			want1:   []interface{}{},
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
			got, got1, err := s.generateReadQuery(context.Background(), tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQL.generateReadQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if found := find(tt.want, got); !found {
				t.Errorf("SQL.generateReadQuery() got = %v,\n want %v", got, tt.want[0])
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SQL.generateReadQuery() got1 = %v,\n want1 %v", got1, tt.want1)
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

func Test_processRows(t *testing.T) {
	type args struct {
		rows        []interface{}
		table       string
		join        []*model.JoinOption
		isAggregate bool
	}
	tests := []struct {
		name   string
		args   args
		result []interface{}
	}{
		{
			name: "normal test case",
			args: args{
				table: "t1",
				rows: []interface{}{
					map[string]interface{}{"t1__c1": "a1", "t1__c2": "a2", "t2__c3": "b1"},
					map[string]interface{}{"t1__c1": "a1", "t1__c2": "a2", "t2__c3": "b2"},
					map[string]interface{}{"t1__c1": "c1", "t1__c2": "c2", "t2__c3": "d1"},
					map[string]interface{}{"t1__c1": "c1", "t1__c2": "c2", "t2__c3": "d2"},
				},
				join: []*model.JoinOption{{Table: "t2"}},
			},
			result: []interface{}{
				map[string]interface{}{"c1": "a1", "c2": "a2", "t2": []interface{}{
					map[string]interface{}{"c3": "b1"},
					map[string]interface{}{"c3": "b2"},
				}},
				map[string]interface{}{"c1": "c1", "c2": "c2", "t2": []interface{}{
					map[string]interface{}{"c3": "d1"},
					map[string]interface{}{"c3": "d2"},
				}},
			},
		}, {
			name: "2 level nested test case",
			args: args{
				table: "t1",
				rows: []interface{}{
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b11", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b11", "t3__tc4": "c21"},
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b21", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b21", "t3__tc4": "c21"},

					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b11", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b11", "t3__tc4": "c21"},
					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b21", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b21", "t3__tc4": "c21"},
				},
				join: []*model.JoinOption{{Table: "t2", Join: []*model.JoinOption{{Table: "t3"}}}},
			},
			result: []interface{}{
				map[string]interface{}{"tc1": "a11", "tc2": "a12", "t2": []interface{}{
					map[string]interface{}{"tc3": "b11", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
					map[string]interface{}{"tc3": "b21", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
				}},
				map[string]interface{}{"tc1": "a21", "tc2": "a22", "t2": []interface{}{
					map[string]interface{}{"tc3": "b11", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
					map[string]interface{}{"tc3": "b21", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
				}},
			},
		}, {
			name: "2 level nested parallel branch test case",
			args: args{
				table: "t1",
				rows: []interface{}{
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b11", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b11", "t3__tc4": "c21"},
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b21", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b21", "t3__tc4": "c21"},

					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b11", "t3__tc4": "c11", "t4__tc5": "c11"},
					map[string]interface{}{"t1__tc1": "a11", "t1__tc2": "a12", "t2__tc3": "b11", "t3__tc4": "c21", "t4__tc5": "c21"},

					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b11", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b11", "t3__tc4": "c21"},
					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b21", "t3__tc4": "c11"},
					map[string]interface{}{"t1__tc1": "a21", "t1__tc2": "a22", "t2__tc3": "b21", "t3__tc4": "c21"},
				},
				join: []*model.JoinOption{{Table: "t2", Join: []*model.JoinOption{{Table: "t3"}}}, {Table: "t4"}},
			},
			result: []interface{}{
				map[string]interface{}{"tc1": "a11", "tc2": "a12", "t2": []interface{}{
					map[string]interface{}{"tc3": "b11", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
					map[string]interface{}{"tc3": "b21", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
				}, "t4": []interface{}{
					map[string]interface{}{"tc5": "c11"},
					map[string]interface{}{"tc5": "c21"},
				}},
				map[string]interface{}{"tc1": "a21", "tc2": "a22", "t2": []interface{}{
					map[string]interface{}{"tc3": "b11", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
					map[string]interface{}{"tc3": "b21", "t3": []interface{}{
						map[string]interface{}{"tc4": "c11"},
						map[string]interface{}{"tc4": "c21"},
					}},
				}, "t4": []interface{}{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQL{}

			finalArray := make([]interface{}, 0)
			mapping := map[string]map[string]interface{}{}
			for _, row := range tt.args.rows {
				s.processRows(context.Background(), false, []string{tt.args.table}, tt.args.isAggregate, row.(map[string]interface{}), tt.args.join, mapping, &finalArray, nil, map[string]map[string]string{})
			}

			for _, elem := range finalArray {
				delete(elem.(map[string]interface{}), "_dbFetchTs")
			}

			if !reflect.DeepEqual(finalArray, tt.result) {
				t.Errorf("processRows() = %v; wanted = %v", finalArray, tt.result)
			}
		})
	}
}
