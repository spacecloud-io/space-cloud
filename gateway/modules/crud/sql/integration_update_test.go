// +build integration

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
	type test struct {
		name            string
		insertQuery     []string
		readQuery       string
		args            args
		want            int64
		wantErr         bool
		wantReadResult  []interface{}
		readResult      []interface{}
		isMysqlSkip     bool
		isPostgresSkip  bool
		isSQLServerSkip bool
	}
	var testCases []test
	tests := []test{
		{
			name: "$set operation on all the supported types of SC",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',509,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1'`,
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
					"is_public":        true,
					"name":             "tata motors",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           5.5,
				},
			},
		},
		{
			name: "$inc operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
			isPostgresSkip: true,
		},
		{
			name: "$mul operation on type Integer Float postgres specifc case",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
					"volume": 48.449999999999996,
				},
			},
			isSQLServerSkip: true,
			isMysqlSkip:     true,
		},
		{
			name: "$max successful operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
			isPostgresSkip: true,
		},
		{
			name: "$max un successful operation on type Integer Float for postgres",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(20),
					"volume": 9.5,
				},
			},
			isMysqlSkip:     true,
			isSQLServerSkip: true,
		},
		{
			name: "$min successful operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(20),
					"volume": 9.5,
				},
			},
			isMysqlSkip: true,
		},
		{
			name: "$min un successful operation on type Integer Float for postgres",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"kind":   int64(20),
					"volume": 9.5,
				},
			},
			isMysqlSkip:     true,
			isSQLServerSkip: true,
		},
		{
			name: "$unset operation on type ID",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1'`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						// TODO: default value of date time not known
						"$unset": map[string]interface{}{
							"parent":      "",
							"kind":        0,
							"is_public":   false,
							"name":        "",
							"description": "{}",
							"volume":      0.0,
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
					"established_date": "2001-11-01T14:29:36Z",
					"parent":           "",
					"kind":             int64(0),
					"is_public":        false,
					"name":             "",
					"description":      "{}",
					"volume":           0.0,
				},
			},
			isPostgresSkip: true,
		},
		{
			name: "$unset operation on type ID for postgres",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1'`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						// TODO: default value of date time not known
						"$unset": map[string]interface{}{
							"parent": "",
							"kind":   0,
							// "is_public": false,
							"name":        "",
							"description": "{}",
							"volume":      0.0,
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
					"established_date": "2001-11-01T14:29:36Z",
					"parent":           "",
					"kind":             int64(0),
					"is_public":        false,
					"name":             "",
					"description":      "{}",
					"volume":           0.0,
				},
			},
			isMysqlSkip:     true,
			isSQLServerSkip: true,
		},
		// TODO:Current Date Operator Remaining
		{
			name: "upsert operation data doesn exists so insert operation",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5)`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1' or id = '2'`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "2",
					},
					Operation: utils.Upsert,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"parent":           "reliance",
							"established_date": "2002-11-01 14:29:36",
							"kind":             10,
							"is_public":        true,
							"name":             "jio",
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
					"parent":           "tata",
					"established_date": "2001-11-01T14:29:36Z",
					"kind":             int64(20),
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance",
					"established_date": "2002-11-01T14:29:36Z",
					"kind":             int64(10),
					"is_public":        true,
					"name":             "jio",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           5.5,
				},
			},
		},
		{
			name: "upsert operation data exists so update operation",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,description,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,false,'tata salt','{"city":"india", "pinCode": 400014}',9.5),
				('2','reliance','2002-11-01 14:29:36',10,true,'jio','{"city":"india", "pinCode": 400014}',5.5)
				`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1' or id = '2'`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "2",
					},
					Operation: utils.Upsert,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"parent": "reliance telecom",
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
					"parent":           "tata",
					"established_date": "2001-11-01T14:29:36Z",
					"kind":             int64(20),
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance telecom",
					"established_date": "2002-11-01T14:29:36Z",
					"kind":             int64(10),
					"is_public":        true,
					"name":             "jio",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           5.5,
				},
			},
		},
	}

	mssqlCases := []test{
		{
			name: "$set operation on all the supported types of SC",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',509,0,'tata salt',9.5)`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1'`,
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
							"is_public":        1,
							"name":             "tata motors",
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
					"is_public":        true,
					"name":             "tata motors",
					"volume":           5.5,
				},
			},
		},
		{
			name: "$inc operation on type Integer Float",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
					"volume": 48.449999999999996,
				},
			},
			isPostgresSkip: true,
		},
		{
			name: "$mul operation on type Integer Float postgres specifc case",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
			},
			readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
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
					"volume": 48.449999999999996,
				},
			},
			isSQLServerSkip: true,
			isMysqlSkip:     true,
		},
		//{
		//	name: "$max successful operation on type Integer Float",
		//	insertQuery: []string{
		//		`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
		//		VALUES
		//		('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
		//	},
		//	readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
		//	args: args{
		//		ctx: context.Background(),
		//		col: "companies",
		//		req: &model.UpdateRequest{
		//			Find: map[string]interface{}{
		//				"id": "1",
		//			},
		//			Operation: utils.All,
		//			Update: map[string]interface{}{
		//				"$max": map[string]interface{}{
		//					"kind":   40,
		//					"volume": 15.5,
		//				},
		//			},
		//		},
		//	},
		//	want:           1,
		//	wantErr:        false,
		//	wantReadResult: []interface{}{},
		//	readResult: []interface{}{
		//		map[string]interface{}{
		//			"kind":   int64(40),
		//			"volume": 15.5,
		//		},
		//	},
		//},
		//{
		//	name: "$max un successful operation on type Integer Float",
		//	insertQuery: []string{
		//		`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
		//		VALUES
		//		('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
		//	},
		//	readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
		//	args: args{
		//		ctx: context.Background(),
		//		col: "companies",
		//		req: &model.UpdateRequest{
		//			Find: map[string]interface{}{
		//				"id": "1",
		//			},
		//			Operation: utils.All,
		//			Update: map[string]interface{}{
		//				"$max": map[string]interface{}{
		//					"kind":   10,
		//					"volume": 5.5,
		//				},
		//			},
		//		},
		//	},
		//	want:           0,
		//	wantErr:        false,
		//	wantReadResult: []interface{}{},
		//	readResult: []interface{}{
		//		map[string]interface{}{
		//			"kind":   int64(20),
		//			"volume": 9.5,
		//		},
		//	},
		//	isPostgresSkip: true,
		//},
		//{
		//	name: "$max un successful operation on type Integer Float for postgres",
		//	insertQuery: []string{
		//		`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
		//		VALUES
		//		('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
		//	},
		//	readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
		//	args: args{
		//		ctx: context.Background(),
		//		col: "companies",
		//		req: &model.UpdateRequest{
		//			Find: map[string]interface{}{
		//				"id": "1",
		//			},
		//			Operation: utils.All,
		//			Update: map[string]interface{}{
		//				"$max": map[string]interface{}{
		//					"kind":   10,
		//					"volume": 5.5,
		//				},
		//			},
		//		},
		//	},
		//	want:           1,
		//	wantErr:        false,
		//	wantReadResult: []interface{}{},
		//	readResult: []interface{}{
		//		map[string]interface{}{
		//			"kind":   int64(20),
		//			"volume": 9.5,
		//		},
		//	},
		//	isMysqlSkip: true,
		//	isSQLServerSkip: true,
		//},
		//{
		//	name: "$min successful operation on type Integer Float",
		//	insertQuery: []string{
		//		`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
		//		VALUES
		//		('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
		//	},
		//	readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
		//	args: args{
		//		ctx: context.Background(),
		//		col: "companies",
		//		req: &model.UpdateRequest{
		//			Find: map[string]interface{}{
		//				"id": "1",
		//			},
		//			Operation: utils.All,
		//			Update: map[string]interface{}{
		//				"$min": map[string]interface{}{
		//					"kind":   10,
		//					"volume": 5.5,
		//				},
		//			},
		//		},
		//	},
		//	want:           1,
		//	wantErr:        false,
		//	wantReadResult: []interface{}{},
		//	readResult: []interface{}{
		//		map[string]interface{}{
		//			"kind":   int64(10),
		//			"volume": 5.5,
		//		},
		//	},
		//},
		//{
		//	name: "$min un successful operation on type Integer Float",
		//	insertQuery: []string{
		//		`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
		//		VALUES
		//		('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
		//	},
		//	readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
		//	args: args{
		//		ctx: context.Background(),
		//		col: "companies",
		//		req: &model.UpdateRequest{
		//			Find: map[string]interface{}{
		//				"id": "1",
		//			},
		//			Operation: utils.All,
		//			Update: map[string]interface{}{
		//				"$min": map[string]interface{}{
		//					"kind":   40,
		//					"volume": 20.5,
		//				},
		//			},
		//		},
		//	},
		//	want:           1,
		//	wantErr:        false,
		//	wantReadResult: []interface{}{},
		//	readResult: []interface{}{
		//		map[string]interface{}{
		//			"kind":   int64(20),
		//			"volume": 9.5,
		//		},
		//	},
		//	isMysqlSkip: true,
		//},
		//{
		//	name: "$min un successful operation on type Integer Float for postgres",
		//	insertQuery: []string{
		//		`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
		//		VALUES
		//		('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
		//	},
		//	readQuery: `SELECT kind,volume FROM myproject.companies where  id = '1'`,
		//	args: args{
		//		ctx: context.Background(),
		//		col: "companies",
		//		req: &model.UpdateRequest{
		//			Find: map[string]interface{}{
		//				"id": "1",
		//			},
		//			Operation: utils.All,
		//			Update: map[string]interface{}{
		//				"$min": map[string]interface{}{
		//					"kind":   40,
		//					"volume": 20.5,
		//				},
		//			},
		//		},
		//	},
		//	want:           1,
		//	wantErr:        false,
		//	wantReadResult: []interface{}{},
		//	readResult: []interface{}{
		//		map[string]interface{}{
		//			"kind":   int64(20),
		//			"volume": 9.5,
		//		},
		//	},
		//	isMysqlSkip: true,
		//	isSQLServerSkip: true,
		//},
		{
			name: "$unset operation on type ID",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES 
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1'`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						// TODO: default value of date time not known
						"$unset": map[string]interface{}{
							"parent":    "",
							"kind":      0,
							"is_public": false,
							"name":      "",
							"volume":    0.0,
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
					"established_date": "2001-11-01T14:29:36Z",
					"parent":           "",
					"kind":             int64(0),
					"is_public":        false,
					"name":             "",
					"volume":           0.0,
				},
			},
			isPostgresSkip: true,
		},
		// TODO:Current Date Operator Remaining
		{
			name: "upsert operation data doesn exists so insert operation",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5)`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1' or id = '2'`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "2",
					},
					Operation: utils.Upsert,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"parent":           "reliance",
							"established_date": "2002-11-01 14:29:36",
							"kind":             10,
							"is_public":        true,
							"name":             "jio",
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
					"parent":           "tata",
					"established_date": "2001-11-01T14:29:36Z",
					"kind":             int64(20),
					"is_public":        false,
					"name":             "tata salt",
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance",
					"established_date": "2002-11-01T14:29:36Z",
					"kind":             int64(10),
					"is_public":        true,
					"name":             "jio",
					"volume":           5.5,
				},
			},
		},
		{
			name: "upsert operation data exists so update operation",
			insertQuery: []string{
				`INSERT INTO myproject.companies (id,parent,established_date,kind,is_public,name,volume)
				VALUES
				('1','tata','2001-11-01 14:29:36',20,0,'tata salt',9.5),
				('2','reliance','2002-11-01 14:29:36',10,1,'jio',5.5)
				`,
			},
			readQuery: `SELECT * FROM myproject.companies where  id = '1' or id = '2'`,
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "2",
					},
					Operation: utils.Upsert,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"parent": "reliance telecom",
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
					"parent":           "tata",
					"established_date": "2001-11-01T14:29:36Z",
					"kind":             int64(20),
					"is_public":        false,
					"name":             "tata salt",
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance telecom",
					"established_date": "2002-11-01T14:29:36Z",
					"kind":             int64(10),
					"is_public":        true,
					"name":             "jio",
					"volume":           5.5,
				},
			},
		},
	}

	switch utils.DBType(*dbType) {
	case utils.MySQL, utils.Postgres:
		testCases = tests
	case utils.SQLServer:
		testCases = mssqlCases
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("Update() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			if *dbType == string(utils.MySQL) && tt.isMysqlSkip {
				return
			}
			if *dbType == string(utils.Postgres) && tt.isPostgresSkip {
				return
			}
			if *dbType == string(utils.SQLServer) && tt.isSQLServerSkip {
				return
			}
			// clear the mutated data in db
			if _, err := db.client.Exec("TRUNCATE TABLE myproject.companies"); err != nil {
				t.Log("Update() Couldn't truncate table", err)
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
			// read the data to check if it is actually updated
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
			if _, err := db.client.Exec("TRUNCATE TABLE myproject.companies"); err != nil {
				t.Log("Update() Couldn't truncate table", err)
			}
		})
	}
}
