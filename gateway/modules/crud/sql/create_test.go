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

func TestSQL_generateCreateQuery(t *testing.T) {
	type fields struct {
		enabled    bool
		connection string
		client     *sqlx.DB
		dbType     string
	}
	type args struct {
		project string
		col     string
		req     *model.CreateRequest
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
			name:   "1",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "one",
					Document:  map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}},
			},
			want:    "INSERT INTO footable1 (string1, string2, string3) VALUES (?, ?, ?)",
			want1:   []interface{}{"1", "2", "3"},
			wantErr: false,
		},
		{
			name:   "2",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}},
				},
			},
			want:    "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)",
			want1:   []interface{}{"1", "2", "1", "2", "1", "2"},
			wantErr: false,
		},
		{
			name:   "3",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}},
			},
			wantErr: true,
		},
		{
			name:   "4",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{1, 2, 3}},
			},
			wantErr: true,
		},
		{
			name:   "5",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "one",
					Document:  map[string]interface{}{}},
			},
			want:    "INSERT INTO footable1 DEFAULT VALUES",
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "6",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{}, map[string]interface{}{"string1": "1"}}},
			},
			wantErr: true,
		},
		{
			name:   "7",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string6": "2"}}},
			},
			wantErr: true,
		},
		{
			name:   "8",
			fields: fields{dbType: "mysql"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}},
			},
			want:    "INSERT INTO footable1 (string1, string2) VALUES (?, ?), (?, ?), (?, ?)",
			want1:   []interface{}{"1", "2", "1", "2", "1", "2"},
			wantErr: false,
		},

		// #######################################################################################
		// ###################################  SQLServer  #######################################
		// #######################################################################################
		{
			name:   "1",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "one",
					Document:  map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}},
			},
			want:    "INSERT INTO foo.footable1 (string1, string2, string3) VALUES (@p1, @p2, @p3)",
			want1:   []interface{}{"1", "2", "3"},
			wantErr: false,
		},
		{
			name:   "2",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}},
				},
			},
			want:    "INSERT INTO foo.footable1 (string1, string2) VALUES (@p1, @p2), (@p3, @p4), (@p5, @p6)",
			want1:   []interface{}{"1", "2", "1", "2", "1", "2"},
			wantErr: false,
		},
		{
			name:   "3",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}},
			},
			wantErr: true,
		},
		{
			name:   "4",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{1, 2, 3}},
			},
			wantErr: true,
		},
		{
			name:   "5",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "one",
					Document:  map[string]interface{}{}},
			},
			want:    "INSERT INTO foo.footable1 DEFAULT VALUES",
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "6",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{}, map[string]interface{}{"string1": "1"}}},
			},
			wantErr: true,
		},
		{
			name:   "7",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string6": "2"}}},
			},
			wantErr: true,
		},
		{
			name:   "8",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}},
			},
			want:    "INSERT INTO foo.footable1 (string1, string2) VALUES (@p1, @p2), (@p3, @p4), (@p5, @p6)",
			want1:   []interface{}{"1", "2", "1", "2", "1", "2"},
			wantErr: false,
		},

		// #######################################################################################
		// ###################################  Postgres  ########################################
		// #######################################################################################
		{
			name:   "1",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "one",
					Document:  map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}},
			},
			want:    "INSERT INTO foo.footable1 (string1, string2, string3) VALUES ($1, $2, $3)",
			want1:   []interface{}{"1", "2", "3"},
			wantErr: false,
		},
		{
			name:   "2",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}},
				},
			},
			want:    "INSERT INTO foo.footable1 (string1, string2) VALUES ($1, $2), ($3, $4), ($5, $6)",
			want1:   []interface{}{"1", "2", "1", "2", "1", "2"},
			wantErr: false,
		},
		{
			name:   "3",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  map[string]interface{}{"string1": "1", "string2": "2", "string3": "3"}},
			},
			wantErr: true,
		},
		{
			name:   "4",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{1, 2, 3}},
			},
			wantErr: true,
		},
		{
			name:   "5",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "one",
					Document:  map[string]interface{}{}},
			},
			want:    "INSERT INTO foo.footable1 DEFAULT VALUES",
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "6",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{}, map[string]interface{}{"string1": "1"}}},
			},
			wantErr: true,
		},
		{
			name:   "7",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string6": "2"}}},
			},
			wantErr: true,
		},
		{
			name:   "8",
			fields: fields{dbType: "postgres"},
			args: args{project: "foo",
				col: "footable1",
				req: &model.CreateRequest{
					Operation: "all",
					Document:  []interface{}{map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}, map[string]interface{}{"string1": "1", "string2": "2"}}},
			},
			want:    "INSERT INTO foo.footable1 (string1, string2) VALUES ($1, $2), ($3, $4), ($5, $6)",
			want1:   []interface{}{"1", "2", "1", "2", "1", "2"},
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
			got, got1, err := s.generateCreateQuery(tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQL.generateCreateQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQL.generateCreateQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SQL.generateCreateQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
