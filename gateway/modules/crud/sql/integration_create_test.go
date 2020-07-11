// +build integration
// use "go test -tags integration" to run integration tests

package sql

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestSQL_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		col string
		req *model.CreateRequest
	}
	type test struct {
		name           string
		readQuery      string
		args           args
		want           int64
		wantErr        bool
		wantReadResult []interface{}
	}
	var testCases []test
	mssqlCases := []test{
		{
			name: "Single Insert",
			args: args{
				ctx: context.Background(),
				col: "customers",
				req: &model.CreateRequest{
					Document: map[string]interface{}{
						"id":         "1",
						"name":       "John",
						"age":        20,
						"height":     5.8,
						"is_prime":   1,
						"birth_date": "2015-11-05 14:29:36",
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
					"is_prime":   true,
					"birth_date": "2015-11-05T14:29:36Z",
				},
			},
			readQuery: `SELECT * FROM myproject.customers WHERE id = '1'`,
		},
		{
			name: "Multiple Insert",
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
							"is_prime":   1,
							"birth_date": "2015-11-05 14:29:36",
						},
						map[string]interface{}{
							"id":         "3",
							"name":       "Amy",
							"age":        int64(40),
							"height":     5.0,
							"is_prime":   0,
							"birth_date": "2015-11-05 14:29:36",
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
					"is_prime":   true,
					"birth_date": "2015-11-05T14:29:36Z",
				},
				map[string]interface{}{
					"id":         "3",
					"name":       "Amy",
					"age":        int64(40),
					"height":     5.0,
					"is_prime":   false,
					"birth_date": "2015-11-05T14:29:36Z",
				},
			},
			readQuery: `SELECT * FROM myproject.customers WHERE id = '2' or id = '3'`,
		},
	}
	sqlCases := []test{
		{
			name: "Single Insert",
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
					"is_prime":   true,
					"birth_date": "2015-11-05T14:29:36Z",
					"address":    `{"city": "pune", "pinCode": 123456}`,
				},
			},
			readQuery: `SELECT * FROM myproject.customers WHERE id = '1'`,
		},
		{
			name: "Multiple Insert",
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
							"is_prime":   true,
							"birth_date": "2015-11-05 14:29:36",
							"address":    `{"city": "california", "pinCode": 567890}`,
						},
						map[string]interface{}{
							"id":         "3",
							"name":       "Amy",
							"age":        int64(40),
							"height":     5.0,
							"is_prime":   false,
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
					"is_prime":   true,
					"birth_date": "2015-11-05T14:29:36Z",
					"address":    `{"city": "california", "pinCode": 567890}`,
				},
				map[string]interface{}{
					"id":         "3",
					"name":       "Amy",
					"age":        int64(40),
					"height":     5.0,
					"is_prime":   false,
					"birth_date": "2015-11-05T14:29:36Z",
					"address":    `{"city": "newYork", "pinCode": 654321}`,
				},
			},
			readQuery: `SELECT * FROM myproject.customers WHERE id = '2' or id = '3'`,
		},
	}

	if utils.DBType(*dbType) == utils.SQLServer {
		testCases = mssqlCases
	} else {
		testCases = sqlCases
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("Create() Couldn't establishing connection with database", dbType)
	}

	// ensure that the table is empty
	if _, err := db.client.Exec("TRUNCATE TABLE myproject.customers"); err != nil {
		t.Log("Create() Couldn't truncate table", err)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// insert data in db
			got, err := db.Create(tt.args.ctx, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}

			// query the database to ensure that data is inserted in database
			rows, err := db.client.Queryx(tt.readQuery)
			if err != nil {
				t.Error("Create() query error", err)
				return
			}

			// store query result
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
						t.Errorf("Create() mismatch in result got %v \n want %v", readValue, value)
					}
				}
			}
		})
	}

	// clear data
	if _, err := db.client.Exec("TRUNCATE TABLE myproject.customers"); err != nil {
		t.Log("Create() Couldn't truncate table", err)
	}
}
