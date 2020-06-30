// +build integration
// use "go test -tags integration" to run integration tests

package mgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		find           map[string]interface{}
		args           args
		want           int64
		wantErr        bool
		wantReadResult []interface{}
	}
	testCases := []test{
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
			find:    map[string]interface{}{"id": "1"},
			wantReadResult: []interface{}{
				map[string]interface{}{
					"id":         "1",
					"name":       "John",
					"age":        int32(20),
					"height":     5.8,
					"is_prime":   true,
					"birth_date": "2015-11-05 14:29:36",
					"address":    `{"city": "pune", "pinCode": 123456}`,
				},
			},
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
			find:    map[string]interface{}{"$or": []interface{}{map[string]interface{}{"id": "3"}, map[string]interface{}{"id": "2"}}},
			wantReadResult: []interface{}{
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
		},
	}

	db, err := Init(true, *connection, "myproject")
	if err != nil {
		t.Fatal("Create() Couldn't establishing connection with database", dbType)
		return
	}

	// ensure that the table is empty
	coll := db.client.Database("myproject").Collection("customers")

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// clear data
			if err := coll.Drop(context.Background()); err != nil {
				t.Log("Create() Couldn't truncate table", err)
			}

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
			results := []interface{}{}
			findOptions := options.Find()
			cur, err := coll.Find(tt.args.ctx, tt.find, findOptions)
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

				results = append(results, doc)
			}

			if err := cur.Err(); err != nil {
				t.Log("Create() got error", err)
			}

			// store query result
			if len(tt.wantReadResult) != len(results) {
				t.Errorf("Create() mismatch in result lenght got %v want %v", len(results), len(tt.wantReadResult))
				return
			}
			for index, result := range tt.wantReadResult {
				for key, value := range result.(map[string]interface{}) {
					readValue, ok := results[index].(map[string]interface{})[key]
					if !ok {
						t.Errorf("Create() missing field key %v at index %v", key, value)
					}
					if !reflect.DeepEqual(readValue, value) {
						t.Errorf("Create() mismatch in result got %v type %v \n want %v", readValue, reflect.TypeOf(readValue), value)
					}
				}
			}
		})
	}

	// clear data
	if err := coll.Drop(context.Background()); err != nil {
		t.Log("Create() Couldn't truncate table", err)
	}
}
