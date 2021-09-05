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

func TestSQL_generateReadAggregateQuery(t *testing.T) {
	// temp := "one"
	testAggregateValue := []string{"table:Column1"}
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
			name:   "Sum",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"sum": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT SUM(Column1) AS aggregate___nested___table___sum___Column1 FROM table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Max",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"max": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT MAX(Column1) AS aggregate___nested___table___max___Column1 FROM table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "min",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"min": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT MIN(Column1) AS aggregate___nested___table___min___Column1 FROM table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "avg",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"avg": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT AVG(Column1) AS aggregate___nested___table___avg___Column1 FROM table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Count",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"count": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT COUNT(Column1) AS aggregate___nested___table___count___Column1 FROM table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Wrong Operation",
			fields: fields{dbType: "mysql"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1}},
					Aggregate: map[string][]string{"wrongOp": testAggregateValue},
					Operation: "all"}},
			want:    []string{""},
			want1:   nil,
			wantErr: true,
		},

		// #######################################################################################
		// ###################################  Postgres  ########################################
		// #######################################################################################
		{
			name:   "Sum",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"sum": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT SUM(Column1) AS aggregate___nested___table___sum___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Max",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"max": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT MAX(Column1) AS aggregate___nested___table___max___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "min",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"min": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT MIN(Column1) AS aggregate___nested___table___min___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "avg",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"avg": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT AVG(Column1) AS aggregate___nested___table___avg___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Count",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"count": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT COUNT(Column1) AS aggregate___nested___table___count___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Wrong Operation",
			fields: fields{dbType: "postgres"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1}},
					Aggregate: map[string][]string{"wrongOp": testAggregateValue},
					Operation: "all"}},
			want:    []string{""},
			want1:   nil,
			wantErr: true,
		},

		// #######################################################################################
		// ###################################  SQLServer  #######################################
		// #######################################################################################
		{
			name:   "Sum",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"sum": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT SUM(Column1) AS aggregate___nested___table___sum___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Max",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"max": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT MAX(Column1) AS aggregate___nested___table___max___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "min",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"min": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT MIN(Column1) AS aggregate___nested___table___min___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "avg",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"avg": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT AVG(Column1) AS aggregate___nested___table___avg___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Count",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{},
					Aggregate: map[string][]string{"count": testAggregateValue},
					Operation: "all"}},
			want:    []string{"SELECT COUNT(Column1) AS aggregate___nested___table___count___Column1 FROM test.table"},
			want1:   []interface{}{},
			wantErr: false,
		},
		{
			name:   "Wrong Operation",
			fields: fields{dbType: "sqlserver"},
			args: args{project: "test", col: "table",
				req: &model.ReadRequest{
					Find:      map[string]interface{}{},
					Options:   &model.ReadOptions{Select: map[string]int32{"Column1": 1}},
					Aggregate: map[string][]string{"wrongOp": testAggregateValue},
					Operation: "all"}},
			want:    []string{""},
			want1:   nil,
			wantErr: true,
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
