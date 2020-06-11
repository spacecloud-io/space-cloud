package sql

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestSQL_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		col string
		req *model.UpdateRequest
	}
	tests := []struct {
		name           string
		insertQuery    []string
		readQuery      string
		args           args
		want           int64
		wantErr        bool
		wantReadResult []interface{}
		readResult     []interface{}
	}{
		{
			name: "$set operation on all the supported types of SC",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",509,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT * FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"parent":           "reliance",
							"established_date": "2002-11-01 14:29:36",
							"kind":             10,
							"is_public":        true,
							"name":             "tata motors",
							"description":      `{"city":"india", "pinCode": 400014}`,
							"volume":           5.5,
						},
					},
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "reliance",
					"established_date": "2002-11-01T14:29:36Z",
					"kind":             int64(10),
					"is_public":        int64(1),
					"name":             "tata motors",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           5.5,
				},
			},
		},
		{
			name: "$inc operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",20,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT kind,volume FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$inc": map[string]interface{}{
							"kind":   10,
							"volume": 5.1,
						},
					},
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(30),
					"volume": 14.6,
				},
			},
		},
		{
			name: "decrement using $inc operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",20,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT kind,volume FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$inc": map[string]interface{}{
							"kind":   -5,
							"volume": -5.1,
						},
					},
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(15),
					"volume": 4.4,
				},
			},
		},
		{
			name: "$mul operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",20,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT kind,volume FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$mul": map[string]interface{}{
							"kind":   10,
							"volume": 5.1,
						},
					},
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(200),
					"volume": 48.45,
				},
			},
		},
		{
			name: "$max successful operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",20,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT kind,volume FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$max": map[string]interface{}{
							"kind":   40,
							"volume": 15.5,
						},
					},
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(40),
					"volume": 15.5,
				},
			},
		},
		{
			name: "$max un successful operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",20,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT kind,volume FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$max": map[string]interface{}{
							"kind":   10,
							"volume": 5.5,
						},
					},
				},
			},
			want:           0,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(20),
					"volume": 9.5,
				},
			},
		},
		{
			name: "$min successful operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",20,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT kind,volume FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$min": map[string]interface{}{
							"kind":   10,
							"volume": 5.5,
						},
					},
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(10),
					"volume": 5.5,
				},
			},
		},
		{
			name: "$min un successful operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				("1","tata","2001-11-01 14:29:36",20,0,"tata salt",'{"city":"india", "pinCode": 400014}',9.5)
				`,
			},
			readQuery: `SELECT kind,volume FROM companies where  id = "1"`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$min": map[string]interface{}{
							"kind":   40,
							"volume": 20.5,
						},
					},
				},
			},
			want:           0,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(20),
					"volume": 9.5,
				},
			},
		},
	}
	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear the mutated data in db
			if _, err := db.client.Exec("TRUNCATE TABLE companies"); err != nil {
				t.Log("Couldn't truncate table", err)
			}
			// insert data in db
			if err := db.RawBatch(context.Background(), tt.insertQuery); err != nil {
				t.Errorf("Update() couldn't insert rows got error - (%v)", err)
			}
			// update the data
			got, err := db.Update(tt.args.ctx, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}

			rows, err := db.client.Queryx(tt.readQuery)
			if err != nil {
				t.Error("Update() query error", err)
				return
			}
			readResult := []interface{}{}
			rowTypes, _ := rows.ColumnTypes()
			for rows.Next() {
				v := map[string]interface{}{}
				if err := rows.MapScan(v); err != nil {
					t.Error("Update() Scanning error", err)
				}
				mysqlTypeCheck(utils.DBType(*dbType), rowTypes, v)
				readResult = append(readResult, v)
			}
			if !reflect.DeepEqual(tt.readResult, readResult) {
				t.Errorf("Update() got (%v) want (%v)", readResult, tt.readResult)
			}

			// clear the mutated data in db
			if _, err := db.client.Exec("TRUNCATE TABLE companies"); err != nil {
				t.Log("Couldn't truncate table", err)
			}
		})
	}
}
