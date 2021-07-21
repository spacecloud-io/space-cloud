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

func TestSQL_generateUpdateQuery(t *testing.T) {
	type fields struct {
		enabled    bool
		connection string
		client     *sqlx.DB
		dbType     string
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
		// #######################################################################################
		// ###################################  MySQL  ###########################################
		// #######################################################################################
		{
			name:   "msyql: valid set find on json",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$set",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}},
					Find: map[string]interface{}{
						"Obj2": map[string]interface{}{"$contains": map[string]interface{}{"obj1": "value1"}},
					},
				},
			},
			want:    "UPDATE col SET String1=? WHERE json_contains(Obj2,?)",
			want1:   []interface{}{"1", `{"obj1":"value1"}`},
			wantErr: false,
		},
		{
			name:   "sql: valid $set query",
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
					},
				},
			},
			want:    "UPDATE col SET String1=? WHERE (FindString1 = ?)",
			want1:   []interface{}{"1", "1"},
			wantErr: false,
		},
		{
			name:   "sql: invalid $set query",
			fields: fields{dbType: "mysql"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$set",
				req: model.UpdateRequest{
					Find: map[string]interface{}{
						"FindString1": "1",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql: valid $currentDate query",
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
					},
				},
			},
			want:    "UPDATE col SET String1=CURRENT_TIMESTAMP WHERE (today = ?)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "sql: invalid $currentDate query",
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
			name:   "sql: invalid $inc query",
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
			name:   "sql: $inc wrong input",
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
			name:   "sql: valid $mul query",
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
					},
				},
			},
			want:    "UPDATE col SET String1=String1*? WHERE (op1 = ?)",
			want1:   []interface{}{int64(6), int64(1)},
			wantErr: false,
		},
		{
			name:   "sql: valid max query",
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
					},
				},
			},
			want:    "UPDATE col SET String1=GREATEST(String1,?) WHERE (op1 = ?)",
			want1:   []interface{}{int64(6132), int64(121)},
			wantErr: false,
		},
		{
			name:   "sql: valid $min ",
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
					},
				},
			},
			want:    "UPDATE col SET String1=LEAST(String1,?) WHERE (op1 = ?)",
			want1:   []interface{}{int64(6), int64(1)},
			wantErr: false,
		},
		{
			name:   "sql: valid $min ",
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
					},
				},
			},
			want:    "UPDATE col SET String1=LEAST(String1,?) WHERE (op1 = ?)",
			want1:   []interface{}{float64(-6.54), int64(1)},
			wantErr: false,
		},
		{
			name:   "sql: valid $min query",
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
					},
				},
			},
			want:    "UPDATE col SET String1=LEAST(String1,?) WHERE (op1 = ?)",
			want1:   []interface{}{int64(18), int64(1)},
			wantErr: false,
		},
		{
			name:   "sql: invalid mul query",
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
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql: valid $currentDate query ",
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
			want:    "UPDATE col SET String1=CURRENT_DATE WHERE (today = ?)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "sql: invalid $currentDate query",
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
			// want1:   []interface{}{},
			wantErr: true,
		},
		{
			name:   "sql: different op",
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
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql: valid $inc query",
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
					},
				},
			},
			want:    "UPDATE col SET String1=String1+? WHERE (op1 = ?)",
			want1:   []interface{}{int64(18446), int64(67)},
			wantErr: false,
		},
		{
			name:   "sql: invalid input to max",
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
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql: invalid ip to min",
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
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql: invalid op",
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
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql: invalid ip to currentdate",
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
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql:currentdate invalid ip",
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
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sql: trying default op",
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
					},
				},
			},
			want:    "",
			wantErr: true,
		},

		// #######################################################################################
		// ###################################  Postgres  ########################################
		// #######################################################################################
		{
			name:   "postgres: valid set find on json",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$set",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}},
					Find: map[string]interface{}{
						"Obj2": map[string]interface{}{"$contains": map[string]interface{}{"obj1": "value1"}},
					},
				},
			},
			want:    "UPDATE project.col SET String1=$1 WHERE Obj2 @> $2",
			want1:   []interface{}{"1", `{"obj1":"value1"}`},
			wantErr: false,
		},
		{
			name:   "postgres: valid set",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$set",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}},
					Find: map[string]interface{}{
						"FindString1": "1",
					},
				},
			},
			want:    "UPDATE project.col SET String1=$1 WHERE (FindString1 = $2)",
			want1:   []interface{}{"1", "1"},
			wantErr: false,
		},
		{
			name:   "postgres: invalid set",
			fields: fields{dbType: "postgres"},
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
			name:   "postgres: valid current date",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": "timestamp"}}},
					Find: map[string]interface{}{
						"today": "1",
					},
				},
			},
			want:    "UPDATE project.col SET String1=CURRENT_TIMESTAMP WHERE (today = $1)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "postgres: invalid currentdate ",
			fields: fields{dbType: "postgres"},
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
			name:   "postgres: inc wrong query",
			fields: fields{dbType: "postgres"},
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
			name:   "postgres: inc wrong ip",
			fields: fields{dbType: "postgres"},
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
			name:   "postgres: valid mul",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$mul",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$mul": map[string]interface{}{"String1": 6}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=String1*$1 WHERE (op1 = $2)",
			want1:   []interface{}{int64(6), int64(1)},
			wantErr: false,
		},
		{
			name:   "postgres valid max query",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": 6132}},
					Find: map[string]interface{}{
						"op1": 121,
					},
				},
			},
			want:    "UPDATE project.col SET String1=GREATEST(String1,$1) WHERE (op1 = $2)",
			want1:   []interface{}{int64(6132), int64(121)},
			wantErr: false,
		}, {
			name:   "postgres: valid max 2 ip",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": 6132, "s2": 12}},
					Find: map[string]interface{}{
						"op1": 121,
					},
				},
			},
			want:    "UPDATE project.col SET String1=GREATEST(String1,$1),s2=GREATEST(s2,$2) WHERE (op1 = $3)",
			want1:   []interface{}{int64(6132), int64(12), int64(121)},
			wantErr: false,
		},
		{
			name:   "postgres: valid min query",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": 6}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,$1) WHERE (op1 = $2)",
			want1:   []interface{}{int64(6), int64(1)},
			wantErr: false,
		},
		{
			name:   "postgres: valid min query",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": -6.54}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,$1) WHERE (op1 = $2)",
			want1:   []interface{}{float64(-6.54), int64(1)},
			wantErr: false,
		},
		{
			name:   "postgres: valid min query int64",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": int64(18)}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,$1) WHERE (op1 = $2)",
			want1:   []interface{}{int64(18), int64(1)},
			wantErr: false,
		},
		{
			name:   "postgres: invalid different op query ",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$mul",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": int64(18446744)}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres: valid currentDate",
			fields: fields{dbType: "postgres"},
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
			want:    "UPDATE project.col SET String1=CURRENT_DATE WHERE (today = $1)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "postgres: currentdate wrong ip",
			fields: fields{dbType: "postgres"},
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
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres: wrong op",
			fields: fields{dbType: "postgres"},
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
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres: valid op ",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$inc",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$inc": map[string]interface{}{"String1": 18446}},
					Find: map[string]interface{}{
						"op1": 67,
					},
				},
			},
			want:    "UPDATE project.col SET String1=String1+$1 WHERE (op1 = $2)",
			want1:   []interface{}{int64(18446), int64(67)},
			wantErr: false,
		},
		{
			name:   "postgres: wrong ip max",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres:wrong ip min",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres:wrong op",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$maxjgf",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres: wrong ip currentdate",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres: current date wrong ip",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": 1}}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "postgres: checking default",
			fields: fields{dbType: "postgres"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDatefs",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDatshdge": map[string]interface{}{"String1": map[string]interface{}{"$type": 1}}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},

		// #######################################################################################
		// ###################################  SQLServer  #######################################
		// #######################################################################################
		{
			name:   "sqlserver: valid set ",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$set",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": "1"}},
					Find: map[string]interface{}{
						"FindString1": "1",
					},
				},
			},
			want:    "UPDATE project.col SET String1=@p1 WHERE (FindString1 = @p2)",
			want1:   []interface{}{"1", "1"},
			wantErr: false,
		},
		{
			name:   "sqlserver: invalid set",
			fields: fields{dbType: "sqlserver"},
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
			name:   "sqlserver: currentdate valid currentdate",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": "timestamp"}}},
					Find: map[string]interface{}{
						"today": "1",
					},
				},
			},
			want:    "UPDATE project.col SET String1=CURRENT_TIMESTAMP WHERE (today = @p1)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "sqlserver: invalid currentdate",
			fields: fields{dbType: "sqlserver"},
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
			name:   "sqlserver: invalid ip inc",
			fields: fields{dbType: "sqlserver"},
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
			name:   "sqlserver:invalid inc",
			fields: fields{dbType: "sqlserver"},
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
			name:   "sqlserver: valid mul",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$mul",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$mul": map[string]interface{}{"String1": 6}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=String1*@p1 WHERE (op1 = @p2)",
			want1:   []interface{}{int64(6), int64(1)},
			wantErr: false,
		},
		{
			name:   "sqlserver: max valid",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": 6132}},
					Find: map[string]interface{}{
						"op1": 121,
					},
				},
			},
			want:    "UPDATE project.col SET String1=GREATEST(String1,@p1) WHERE (op1 = @p2)",
			want1:   []interface{}{int64(6132), int64(121)},
			wantErr: false,
		},
		{
			name:   "sqlserver: max 2ip",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": 6132, "s2": 12}},
					Find: map[string]interface{}{
						"op1": 121,
					},
				},
			},
			want:    "UPDATE project.col SET String1=GREATEST(String1,@p1),s2=GREATEST(s2,@p2) WHERE (op1 = @p3)",
			want1:   []interface{}{int64(6132), int64(12), int64(121)},
			wantErr: false,
		},
		{
			name:   "sqlserver: int64 min",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": 6}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,@p1) WHERE (op1 = @p2)",
			want1:   []interface{}{int64(6), int64(1)},
			wantErr: false,
		},
		{
			name:   "sqlserver: valid min",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": -6.54}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,@p1) WHERE (op1 = @p2)",
			want1:   []interface{}{float64(-6.54), int64(1)},
			wantErr: false,
		},
		{
			name:   "sqlserver: int64 min",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": int64(18)}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "UPDATE project.col SET String1=LEAST(String1,@p1) WHERE (op1 = @p2)",
			want1:   []interface{}{int64(18), int64(1)},
			wantErr: false,
		},
		{
			name:   "sqlserver: different op",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$mul",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$set": map[string]interface{}{"String1": int64(18446744)}},
					Find: map[string]interface{}{
						"op1": 1,
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sqlserver: valid current date",
			fields: fields{dbType: "sqlserver"},
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
			want:    "UPDATE project.col SET String1=CAST( GETDATE() AS date ) WHERE (today = @p1)",
			want1:   []interface{}{"1"},
			wantErr: false,
		},
		{
			name:   "sqlserver: wrong ip currentdate",
			fields: fields{dbType: "sqlserver"},
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
			want:    "",
			wantErr: true,
		},
		{
			name:   "sqlserver: different op",
			fields: fields{dbType: "sqlserver"},
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
			want:    "",
			wantErr: true,
		},
		{
			name:   "sqlserver: valid inc",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$inc",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$inc": map[string]interface{}{"String1": 18446}},
					Find: map[string]interface{}{
						"op1": 67,
					},
				},
			},
			want:    "UPDATE project.col SET String1=String1+@p1 WHERE (op1 = @p2)",
			want1:   []interface{}{int64(18446), int64(67)},
			wantErr: false,
		},
		{
			name:   "sqlserver: max wrong ip",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$max",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sqlserver: wrong ip min",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$min",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$min": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sqlserver: wrong op",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$maxjgf",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$max": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sqlserver:wrong ip to current date",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": "s18446"}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "sqlserver: wrong ip to currentdate",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDate",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDate": map[string]interface{}{"String1": map[string]interface{}{"$type": 1}}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "checking default sqlserver",
			fields: fields{dbType: "sqlserver"},
			args: args{
				ctx:     context.TODO(),
				project: "project",
				col:     "col",
				op:      "$currentDatefs",
				req: model.UpdateRequest{
					Update: map[string]interface{}{"$currentDatshdge": map[string]interface{}{"String1": map[string]interface{}{"$type": 1}}},
					Find: map[string]interface{}{
						"op1": "67",
					},
				},
			},
			want:    "",
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

			got, got1, err := s.generateUpdateQuery(context.Background(), tt.args.col, &tt.args.req, tt.args.op)
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
