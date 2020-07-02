// +build integration

package mgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		name           string
		insertQuery    []interface{}
		find           map[string]interface{}
		args           args
		want           int64
		wantErr        bool
		wantReadResult []interface{}
		readResult     []interface{}
		selectFields   map[string]int32
	}
	testCases := []test{
		{
			name: "$set operation on all the supported types of SC",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{"id": "1"},
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
							"description":      `{"city": "india", "pinCode": 400014}`,
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
					"established_date": "2002-11-01 14:29:36",
					"kind":             int32(10),
					"is_public":        true,
					"name":             "tata motors",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           5.5,
				},
			},
		},
		{
			name: "$inc operation on type Integer Float",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			selectFields: map[string]int32{"kind": 1, "volume": 1},
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
							"kind":   int32(10),
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
					"kind":   int32(30),
					"volume": 14.6,
				},
			},
		},
		{
			name: "decrement using $inc operation on type Integer Float",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			selectFields: map[string]int32{"kind": 1, "volume": 1},
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
					"kind":   int32(15),
					"volume": 4.4,
				},
			},
		},
		{
			name: "$mul operation on type Integer Float",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			selectFields: map[string]int32{"kind": 1, "volume": 1},
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
					"kind":   int32(200),
					"volume": 48.449999999999996,
				},
			},
		},
		{
			name: "$max successful operation on type Integer Float",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			selectFields: map[string]int32{"kind": 1, "volume": 1},
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
					"kind":   int32(40),
					"volume": 15.5,
				},
			},
		},
		{
			name: "$max un successful operation on type Integer Float ",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			selectFields: map[string]int32{"kind": 1, "volume": 1},
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
					"kind":   int32(20),
					"volume": 9.5,
				},
			},
		},
		{
			name: "$min successful operation on type Integer Float",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			selectFields: map[string]int32{"kind": 1, "volume": 1},
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
					"kind":   int32(10),
					"volume": 5.5,
				},
			},
		},
		{
			name: "$min un successful operation on type Integer Float",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			selectFields: map[string]int32{"kind": 1, "volume": 1},
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
					"kind":   int32(20),
					"volume": 9.5,
				},
			},
		},
		{
			name: "$unset operation on type ID",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{
				"id": "1",
			},
			args: args{
				ctx: context.Background(),
				col: "companies",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$unset": map[string]interface{}{
							"parent":           "",
							"kind":             0,
							"is_public":        false,
							"name":             "",
							"description":      "{}",
							"volume":           0.0,
							"established_date": "",
						},
					},
				},
			},
			want:           1,
			wantErr:        false,
			wantReadResult: []interface{}{},
			readResult: []interface{}{
				map[string]interface{}{
					"id": "1",
				},
			},
		},
		// TODO:Current Date Operator Remaining
		{
			name: "upsert operation data doesn exists so insert operation",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{},
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
							"description":      `{"city": "india", "pinCode": 400014}`,
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
					"established_date": "2001-11-01 14:29:36",
					"kind":             int32(20),
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance",
					"established_date": "2002-11-01 14:29:36",
					"kind":             int32(10),
					"is_public":        true,
					"name":             "jio",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           5.5,
				},
			},
		},
		{
			name: "upsert operation data exists so update operation",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tata",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
			},
			find: map[string]interface{}{},
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
					"established_date": "2001-11-01 14:29:36",
					"kind":             int32(20),
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city": "india", "pinCode": 400014}`,
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":     "2",
					"parent": "reliance telecom",
				},
			},
		},
	}

	db, err := Init(true, *connection, "myproject")
	if err != nil {
		t.Fatal("Update() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// clear the mutated data in db
			if err := db.client.Database("myproject").Collection("companies").Drop(context.Background()); err != nil {
				t.Log("Create() Couldn't truncate table", err)
			}

			// insert data in db
			if _, err := db.client.Database("myproject").Collection("companies").InsertMany(context.Background(), tt.insertQuery); err != nil {
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
			results := []interface{}{}
			findOptions := options.Find()
			if tt.selectFields != nil {
				findOptions = findOptions.SetProjection(tt.selectFields)
			}
			cur, err := db.client.Database("myproject").Collection("companies").Find(tt.args.ctx, tt.find, findOptions)
			if err != nil {
				t.Log("Create() got error", err)
			}
			defer func() { _ = cur.Close(tt.args.ctx) }()

			var count int64
			// Finding multiple documents returns a cursor
			// Iterating through the cursor allows us to decode documents one at a time
			for cur.Next(tt.args.ctx) {
				// Increment the counter
				count++

				// Read the document
				var doc map[string]interface{}
				err := cur.Decode(&doc)
				if err != nil {
					t.Log("Create() got error", err)
				}
				delete(doc, "_id")
				results = append(results, doc)
			}

			if err := cur.Err(); err != nil {
				t.Log("Create() got error", err)
			}

			if !reflect.DeepEqual(tt.readResult, results) {
				t.Errorf("Update() got (%v) \n want (%v)", results, tt.readResult)
			}

			// clear the mutated data in db
			if err := db.client.Database("myproject").Collection("companies").Drop(context.Background()); err != nil {
				t.Log("Create() Couldn't truncate table", err)
			}
		})
	}
}
