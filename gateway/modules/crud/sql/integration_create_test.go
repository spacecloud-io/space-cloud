// +build integration
// use "go test -tags integration" to run integration tests

package sql

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

var dbType = flag.String("db_type", "", "db_type of test case to be run")
var connection = flag.String("conn", "", "connection string of the database")

func TestSQL_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		col string
		req *model.CreateRequest
	}
	tests := []struct {
		name           string
		readQuery      string
		args           args
		want           int64
		wantErr        bool
		wantReadResult []interface{}
	}{
		{
			name: "Single Simple Insert",
			args: args{
				ctx: context.Background(),
				col: "customers",
				req: &model.CreateRequest{
					Document: map[string]interface{}{
						"id":         "1",
						"name":       "John",
						"age":        20,
						"height":     5.8,
						"is_prime":   true,
						"birth_date": "2015-11-05 14:29:36",
						"address":    `{"city": "pune", "pinCode": 123456}`,
					},
					Operation: utils.One,
				},
			},
			want:    1,
			wantErr: false,
			wantReadResult: []interface{}{
				map[string]interface{}{
					"id":         "1",
					"name":       "John",
					"age":        int64(20),
					"height":     5.8,
					"is_prime":   int64(1),
					"birth_date": "2015-11-05T14:29:36Z",
					"address":    `{"city": "pune", "pinCode": 123456}`,
				},
			},
			readQuery: `SELECT * FROM customers WHERE id = "1"`,
		},
		{
			name: "Multiple Simple Insert",
			args: args{
				ctx: context.Background(),
				col: "customers",
				req: &model.CreateRequest{
					Document: []interface{}{
						map[string]interface{}{
							"id":         "2",
							"name":       "Sam",
							"age":        int64(30),
							"height":     6.2,
							"is_prime":   int64(1),
							"birth_date": "2015-11-05 14:29:36",
							"address":    `{"city": "california", "pinCode": 567890}`,
						},
						map[string]interface{}{
							"id":         "3",
							"name":       "Amy",
							"age":        int64(40),
							"height":     5.0,
							"is_prime":   int64(0),
							"birth_date": "2015-11-05 14:29:36",
							"address":    `{"city": "newYork", "pinCode": 654321}`,
						},
					},
					Operation: utils.All,
				},
			},
			want:    2,
			wantErr: false,
			wantReadResult: []interface{}{
				map[string]interface{}{
					"id":         "2",
					"name":       "Sam",
					"age":        int64(30),
					"height":     6.2,
					"is_prime":   int64(1),
					"birth_date": "2015-11-05T14:29:36Z",
					"address":    `{"city": "california", "pinCode": 567890}`,
				},
				map[string]interface{}{
					"id":         "3",
					"name":       "Amy",
					"age":        int64(40),
					"height":     5.0,
					"is_prime":   int64(0),
					"birth_date": "2015-11-05T14:29:36Z",
					"address":    `{"city": "newYork", "pinCode": 654321}`,
				},
			},
			readQuery: `SELECT * FROM customers WHERE id = "2" or id = "3"`,
		},
	}
	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("Couldn't establishing connection with database", dbType)
	}
	if _, err := db.client.Exec("TRUNCATE TABLE customers"); err != nil {
		t.Log("Couldn't truncate table", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Create(tt.args.ctx, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			rows, err := db.client.Queryx(tt.readQuery)
			if err != nil {
				t.Error("Create() query error", err)
				return
			}
			readResult := []interface{}{}
			rowTypes, _ := rows.ColumnTypes()
			for rows.Next() {
				v := map[string]interface{}{}
				if err := rows.MapScan(v); err != nil {
					t.Error("Create() Scanning error", err)
				}
				mysqlTypeCheck(utils.DBType(*dbType), rowTypes, v)
				readResult = append(readResult, v)
			}
			if len(tt.wantReadResult) != len(readResult) {
				t.Errorf("Create() mismatch in result lenght got %v want %v", len(readResult), len(tt.wantReadResult))
			}
			for index, result := range tt.wantReadResult {
				for key, value := range result.(map[string]interface{}) {
					readValue, ok := readResult[index].(map[string]interface{})[key]
					if !ok {
						t.Errorf("Create() missing field key %v at index %v", key, value)
					}
					if !reflect.DeepEqual(readValue, value) {
						// switch value.(type) {
						// case :
						// 	t.Errorf("Create() mismatch in result got %v \n want %v", string(readValue.([]byte)), string(value.([]byte)))
						// }
						t.Errorf("Create() mismatch in result got %v \n want %v", readValue, value)
					}
				}
			}
		})
	}
	if _, err := db.client.Exec("TRUNCATE TABLE customers"); err != nil {
		t.Log("Couldn't truncate table", err)
	}
}
