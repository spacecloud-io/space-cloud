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

func TestSQL_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		col string
		req *model.DeleteRequest
	}
	type test struct {
		name           string
		insertQuery    []interface{}
		find           map[string]interface{}
		args           args
		want           int64
		wantErr        bool
		wantReadResult []interface{}
	}
	testCases := []test{
		{
			name: "Simple Delete",
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tate",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city":"india", "pinCode": 400014}`,
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance",
					"established_date": "2002-11-01 14:29:36",
					"kind":             30,
					"is_public":        true,
					"name":             "jio",
					"description":      `{"city":"india", "pinCode": 400014}`,
					"volume":           18.5,
				},
			},
			find: map[string]interface{}{},
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
			insertQuery: []interface{}{
				map[string]interface{}{
					"id":               "1",
					"parent":           "tate",
					"established_date": "2001-11-01 14:29:36",
					"kind":             20,
					"is_public":        false,
					"name":             "tata salt",
					"description":      `{"city":"india", "pinCode": 400014}`,
					"volume":           9.5,
				},
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance",
					"established_date": "2002-11-01 14:29:36",
					"kind":             30,
					"is_public":        true,
					"name":             "jio",
					"description":      `{"city":"india", "pinCode": 400014}`,
					"volume":           18.5,
				},
			},
			find: map[string]interface{}{
				"$or": []interface{}{map[string]interface{}{"id": "1"}, map[string]interface{}{"id": "2"}},
			},
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
			want:    1,
			wantErr: false,
			wantReadResult: []interface{}{
				map[string]interface{}{
					"id":               "2",
					"parent":           "reliance",
					"established_date": "2002-11-01 14:29:36",
					"kind":             int32(30),
					"is_public":        true,
					"name":             "jio",
					"description":      `{"city":"india", "pinCode": 400014}`,
					"volume":           18.5,
				},
			},
		},
	}

	db, err := Init(true, *connection, "myproject")
	if err != nil {
		t.Fatal("Create() Couldn't establishing connection with database", dbType)
	}

	// ensure that the table is empty
	coll := db.client.Database("myproject").Collection("companies")
	if err := coll.Drop(context.Background()); err != nil {
		t.Log("Create() Couldn't truncate table", err)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// insert data in db
			if _, err := coll.InsertMany(tt.args.ctx, tt.insertQuery); err != nil {
				t.Errorf("Delete() cannot insert data error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := db.Delete(tt.args.ctx, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}

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
				delete(doc, "_id")
				results = append(results, doc)
			}

			if err := cur.Err(); err != nil {
				t.Log("Create() got error", err)
			}

			if !reflect.DeepEqual(tt.wantReadResult, results) {
				t.Errorf("Delete() got (%v) \n want (%v)", results, tt.wantReadResult)
			}

			// clear the mutated data in db
			if err := coll.Drop(context.Background()); err != nil {
				t.Log("Create() Couldn't truncate table", err)
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
		createQuery map[string]interface{}
		collection  string
		readQuery   string
		wantResult  []interface{}
		args        args
		wantErr     bool
	}{
		{
			name:        "Delete table",
			createQuery: map[string]interface{}{"id": "1", "name": "string"},
			collection:  "abcd",
			args: args{
				ctx: context.Background(),
				col: "abcd",
			},
			wantErr: false,
		},
	}
	db, err := Init(true, *connection, "myproject")
	if err != nil {
		t.Fatal("DeleteCollection Couldn't establishing connection with database", dbType)
	}
	coll := db.client.Database("myproject").Collection("companies")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create table
			if _, err := db.client.Database("myproject").Collection(tt.collection).InsertOne(tt.args.ctx, tt.createQuery); err != nil {
				t.Errorf("DeleteCollection couldn't insert rows got error - (%v)", err)
			}

			// delete table
			if err := db.DeleteCollection(tt.args.ctx, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}

			// check if table is actually deleted
			collections, err := db.client.Database("myproject").ListCollectionNames(tt.args.ctx, map[string]interface{}{})
			if err != nil {
				t.Error("DeleteCollection query error", err)
				return
			}
			isFound := false
			for _, collection := range collections {
				if collection == tt.collection {
					isFound = true
				}
			}

			if isFound {
				t.Errorf("Delete() collection not deleted")
			}

		})
	}

	// delete abcd table
	if err := coll.Drop(context.Background()); err != nil {
		t.Log("Create() Couldn't truncate table", err)
	}
}
