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

func TestSQL_generateDeleteQuery(t *testing.T) {
	type fields struct {
		enabled    bool
		connection string
		client     *sqlx.DB
		dbType     string
	}
	type args struct {
		project string
		col     string
		req     *model.DeleteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		// #######################################################################################
		// ###################################  MySQL  ###########################################
		// #######################################################################################
		{
			name:   "Successfull Test json in where clause",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"Obj1": map[string]interface{}{"$contains": map[string]interface{}{"obj1": "value1"}}}},
			},
			want:    "DELETE FROM fooTable WHERE json_contains(Obj1,?)",
			want1:   []interface{}{`{"obj1":"value1"}`},
			wantErr: false,
		},
		{
			name:   "Successfull Test",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": "1"}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 = ?)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "Successfull Test",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": "1"}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 = ?)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Equal To",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 = ?)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface not Equal To",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 != ?)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Greater than ",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 > ?)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Greater than equal to",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 >= ?)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Less than ",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 < ?)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Less than equal to",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 <= ?)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface IN",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 IN (?))",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface NOT IN",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}},
			},
			want:    "DELETE FROM fooTable WHERE (String1 NOT IN (?))",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map OR",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}},
			},
			want:    "DELETE FROM fooTable WHERE ((string1ofstring1 = ?) OR (string1ofstring2 = ?))",
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:   "When length is 0",
			fields: fields{dbType: "mysql"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{}},
			},
			want:    "DELETE FROM fooTable",
			want1:   []interface{}{},
			wantErr: false,
		},

		// #######################################################################################
		// ###################################  SQLServer  #######################################
		// #######################################################################################
		{
			name:   "Successfull Test",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": "1"}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 = @p1)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Equal To",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 = @p1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface not Equal To",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 != @p1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Greater than ",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 > @p1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Greater than equal to",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 >= @p1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Less than ",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 < @p1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Less than equal to",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 <= @p1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface IN",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 IN (@p1))",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface NOT IN",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 NOT IN (@p1))",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map OR",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE ((string1ofstring1 = @p1) OR (string1ofstring2 = @p2))",
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:   "When length is 0",
			fields: fields{dbType: "sqlserver"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{}},
			},
			want:    "DELETE FROM projectName.fooTable",
			want1:   []interface{}{},
			wantErr: false,
		},

		// #######################################################################################
		// ###################################  Postgres  ########################################
		// #######################################################################################
		{
			name:   "Successfull Test",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"Obj1": map[string]interface{}{"$contains": map[string]interface{}{"obj1": "value1"}}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE Obj1 @> $1",
			want1:   []interface{}{`{"obj1":"value1"}`},
			wantErr: false,
		},
		{
			name:   "Successfull Test",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String2": "2"}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String2 = $1)",
			want1:   []interface{}{"2"},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Equal To",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$eq": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 = $1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface not Equal To",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$ne": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 != $1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Greater than ",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gt": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 > $1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Greater than equal to",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$gte": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 >= $1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Less than ",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lt": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 < $1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface Less than equal to",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$lte": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 <= $1)",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface IN",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$in": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 IN ($1))",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map Interface NOT IN",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"String1": map[string]interface{}{"$nin": 1}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE (String1 NOT IN ($1))",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name:   "Nested Map OR",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"string1ofstring1": "1"}, map[string]interface{}{"string1ofstring2": "2"}}}},
			},
			want:    "DELETE FROM projectName.fooTable WHERE ((string1ofstring1 = $1) OR (string1ofstring2 = $2))",
			want1:   []interface{}{"1", "2"},
			wantErr: false,
		},
		{
			name:   "When length is 0",
			fields: fields{dbType: "postgres"},
			args: args{
				project: "projectName",
				col:     "fooTable",
				req:     &model.DeleteRequest{Find: map[string]interface{}{}},
			},
			want:    "DELETE FROM projectName.fooTable",
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
			got, got1, err := s.generateDeleteQuery(context.Background(), tt.args.req, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQL.generateDeleteQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQL.generateDeleteQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SQL.generateDeleteQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
