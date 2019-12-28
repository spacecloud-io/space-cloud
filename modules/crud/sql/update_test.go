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
	"github.com/spaceuptech/space-cloud/model"
)

func TestSQL_generateUpdateQuery(t *testing.T) {
	type fields struct {
		enabled            bool
		connection         string
		client             *sqlx.DB
		dbType             string
		removeProjectScope bool
	}
	type args struct {
		ctx     context.Context
		project string
		col     string
		req     model.UpdateRequest
		op      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$set",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}},
					Find: map[string]interface{}{
						"FindString1": "1",
						"FindString2": "2",
					},
				},
			},
			want:    "UPDATE project.col SET String1=? WHERE ((FindString1 = ?) AND (FindString2 = ?))",
			want1:   []interface{}{"1", "1", "2"},
			wantErr: false,
		},
		{
			name:   "test2",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$set",
				req: model.UpdateRequest{

					Find: map[string]interface{}{
						"FindString1": "1",
						"FindString2": "2",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "test3",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": "timestamp"}}},
					Find: map[string]interface{}{
						"today": "1",
						"op2":   "2",
					},
				},
			},
			want: "UPDATE project.col SET String1=CURRENT_TIMESTAMP WHERE ((today = ?) AND (op2 = ?))",

			wantErr: false,
		},
		{
			name:   "test4",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "test5",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$inc",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$inc": map[string]interface{}{"String1": "1"}},
					Find: map[string]interface{}{
						"today": "1",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "test6",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$inc",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$inc": map[string]interface{}{"String1": "r"}},
					Find: map[string]interface{}{
						"today": "d",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "test7",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$mul",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$mul": map[string]interface{}{"String1": 6}},
					Find: map[string]interface{}{
						"op1": 1,
						"op2": 2,
					},
				},
			},
			want:    "UPDATE project.col SET String1=String1*? WHERE ((op1 = ?) AND (op2 = ?))",
			want1:   []interface{}{int64(6), int64(1), int64(2)},
			wantErr: false,
		},
		{
			name:   "test8",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": 6132}},
					Find: map[string]interface{}{
						"op1": 121,
						"op2": 21,
					},
				},
			},
			want:    "UPDATE project.col SET String1=GREATEST(String1,?) WHERE ((op1 = ?) AND (op2 = ?))",
			want1:   []interface{}{int64(6132), int64(121), int64(21)},
			wantErr: false,
		},
		{
			name:   "test9",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": 6}},
					Find: map[string]interface{}{
						"op1": 1,
						"op2": 2,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,?) WHERE ((op1 = ?) AND (op2 = ?))",
			want1:   []interface{}{int64(6), int64(1), int64(2)},
			wantErr: false,
		},
		{
			name:   "test10",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": -6.54}},
					Find: map[string]interface{}{
						"op1": 1,
						"op2": 2,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,?) WHERE ((op1 = ?) AND (op2 = ?))",
			want1:   []interface{}{float64(-6.54), int64(1), int64(2)},
			wantErr: false,
		},
		{
			name:   "test11",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": int64(18)}},
					Find: map[string]interface{}{
						"op1": 1,
						"op2": 2,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,?) WHERE ((op1 = ?) AND (op2 = ?))",
			want1:   []interface{}{int64(18), int64(1), int64(2)},
			wantErr: false,
		},
		{
			name:   "test12",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$mul",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": int64(18446744)}},
					Find: map[string]interface{}{
						"op1": 1,
						"op2": 2,
					},
				},
			},
			want: "",

			wantErr: true,
		},
		{
			name:   "test13",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": "date"}}},
					Find: map[string]interface{}{
						"today": "1",
					},
				},
			},
			want: "UPDATE project.col SET String1=CURRENT_DATE WHERE (today = ?)",

			wantErr: false,
		},
		{
			name:   "test14",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": ""}}},
					Find: map[string]interface{}{
						"today": "1",
					},
				},
			},
			want: "",
			//want1:   []interface{}{},
			wantErr: true,
		},
		{
			name:   "test15",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$mul",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": int64(18446744)}},
					Find:   map[string]interface{}{},
				},
			},
			want: "",

			wantErr: true,
		},
		{
			name:   "test16",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$inc",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$inc": map[string]interface{}{"String1": 18446}},
					Find: map[string]interface{}{
						"op1": 67,
						"op2": 78,
					},
				},
			},
			want:    "UPDATE project.col SET String1=String1+? WHERE ((op1 = ?) AND (op2 = ?))",
			want1:   []interface{}{int64(18446), int64(67), int64(78)},
			wantErr: false,
		},
		{
			name:   "test17",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
						"op2": "78",
					},
				},
			},
			want: "",

			wantErr: true,
		},
		{
			name:   "test18",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
						"op2": "78",
					},
				},
			},
			want: "",

			wantErr: true,
		},
		{
			name:   "test19",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$maxjgf",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
						"op2": "78",
					},
				},
			},
			want: "",

			wantErr: true,
		},
		{
			name:   "test20",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
						"op2": "78",
					},
				},
			},
			want: "",

			wantErr: true,
		},
		{
			name:   "test21",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": 1}}},
					Find: map[string]interface{}{
						"op1": "67",
						"op2": "78",
					},
				},
			},
			want: "",

			wantErr: true,
		},
		{
			name:   "test22",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDatefs",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDatshdge": map[string]interface{}{"String1": map[string]interface{}{"$type": 1}}},
					Find: map[string]interface{}{
						"op1": "67",
						"op2": "78",
					},
				},
			},
			want: "",

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQL{
				enabled:            tt.fields.enabled,
				connection:         tt.fields.connection,
				client:             tt.fields.client,
				dbType:             tt.fields.dbType,
				removeProjectScope: tt.fields.removeProjectScope,
			}
			got, got1, err := s.generateUpdateQuery(tt.args.ctx, tt.args.project, tt.args.col, &tt.args.req, tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("name = %v, SQL.generateUpdateQuery() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQL.generateUpdateQuery() got = %v, want %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SQL.generateUpdateQuery() got1 = %v, want1 %v", got1, tt.want1)
				return
			}

		})
	}
}
