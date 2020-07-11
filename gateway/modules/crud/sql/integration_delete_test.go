// +build integration

package sql

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestSQL_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		col string
		req *model.DeleteRequest
	}
	type test struct {
		name           string
		insertQuery    []string
		readQuery      string
		args           args
		want           int64
		wantErr        bool
		wantReadResult []interface{}
	}
	var testCases []test
	mssqlCases := []test{
		{
			name: "Simple Delete",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5),
				('2','reliance','2002-11-01 14:29:36',30,1,'jio',18.5)
				`,
			},
			readQuery: "SELECT * FROM myproject.companies",
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.DeleteRequest{
					Operation: utils.All,
				},
			},
			want:           2,
			wantErr:        false,
			wantReadResult: []interface{}{},
		},
		{
			name: "Delete with where clause",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5),
				('2','reliance','2002-11-01 14:29:36',30,1,'jio',18.5)
				`,
			},
			readQuery: "SELECT * FROM myproject.companies where id = '1'",
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.DeleteRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
		},
	}
	sqlCases := []test{
		{
			name: "Simple Delete",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5),
				('2','reliance','2002-11-01 14:29:36',30,true,'jio','{"city":"india", "pinCode": 400014}',18.5)
				`,
			},
			readQuery: "SELECT * FROM myproject.companies",
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.DeleteRequest{
					Operation: utils.All,
				},
			},
			want:           2,
			wantErr:        false,
			wantReadResult: []interface{}{},
		},
		{
			name: "Delete with where clause",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5),
				('2','reliance','2002-11-01 14:29:36',30,true,'jio','{"city":"india", "pinCode": 400014}',18.5)
				`,
			},
			readQuery: "SELECT * FROM myproject.companies where id = '1'",
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.DeleteRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
		},
	}

	if utils.DBType(*dbType) == utils.SQLServer {
		testCases = mssqlCases
	} else {
		testCases = sqlCases
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("Delete() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// insert data in db
			if err := db.RawBatch(context.Background(), tt.insertQuery); err != nil {
				t.Errorf("Delete() couldn't insert rows got error - (%v)", err)
			}
			got, err := db.Delete(tt.args.ctx, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}

			rows, err := db.client.Queryx(tt.readQuery)
			if err != nil {
				t.Error("Delete() query error", err)
				return
			}

			readResult := []interface{}{}
			rowTypes, _ := rows.ColumnTypes()
			for rows.Next() {
				v := map[string]interface{}{}
				if err := rows.MapScan(v); err != nil {
					t.Error("Delete() Scanning error", err)
				}
				mysqlTypeCheck(utils.DBType(*dbType), rowTypes, v)
				readResult = append(readResult, v)
			}
			if !reflect.DeepEqual(tt.wantReadResult, readResult) {
				t.Errorf("Delete() got (%v) want (%v)", readResult, tt.wantReadResult)
			}

			// clear the mutated data in db
			if _, err := db.client.Exec("TRUNCATE TABLE myproject.companies"); err != nil {
				t.Log("Delete() Couldn't truncate table", err)
			}
		})
	}
}

func TestSQL_DeleteCollection(t *testing.T) {
	type args struct {
		ctx context.Context
		col string
	}
	tests := []struct {
		name        string
		createQuery []string
		readQuery   string
		wantResult  []interface{}
		args        args
		wantErr     bool
	}{
		{
			name:        "Delete table",
			createQuery: []string{"CREATE TABLE myproject.abcd (id VARCHAR(50))"},
			readQuery:   "SELECT * FROM information_schema.TABLES where TABLE_SCHEMA = 'myproject' AND TABLE_NAME = 'abcd';",
			args: args{
				ctx: context.Background(),
				col: "abcd",
			},
			wantErr:    false,
			wantResult: []interface{}{},
		},
	}
	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("DeleteCollection Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create table
			if err := db.RawBatch(context.Background(), tt.createQuery); err != nil {
				t.Errorf("DeleteCollection couldn't insert rows got error - (%v)", err)
			}

			// delete table
			if err := db.DeleteCollection(tt.args.ctx, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}

			// check if table is actually deleted
			rows, err := db.client.Queryx(tt.readQuery)
			if err != nil {
				t.Error("DeleteCollection query error", err)
				return
			}

			readResult := []interface{}{}
			rowTypes, _ := rows.ColumnTypes()
			for rows.Next() {
				v := map[string]interface{}{}
				if err := rows.MapScan(v); err != nil {
					t.Error("DeleteCollection Scanning error", err)
				}
				mysqlTypeCheck(utils.DBType(*dbType), rowTypes, v)
				readResult = append(readResult, v)
			}
			if !reflect.DeepEqual(tt.wantResult, readResult) {
				t.Errorf("DeleteCollection got (%v) want (%v)", readResult, tt.wantResult)
			}
		})
	}

	// delete abcd table
	if _, err := db.client.Exec("DROP TABLE IF EXISTS abcd"); err != nil {
		t.Log("DeleteCollection() Couldn't drop table", err)
	}
}
