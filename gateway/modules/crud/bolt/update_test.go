package bolt

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestBolt_Update(t *testing.T) {
	type fields struct {
		enabled    bool
		connection string
	}
	type args struct {
		ctx context.Context
		col string
		req *model.UpdateRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        int64
		want1       interface{}
		wantErr     bool
		readRequest *model.ReadRequest
	}{
		{
			name: "update single document",
			want: 1,
			want1: map[string]interface{}{
				"_id":           "1",
				"name":          "sharad regoti",
				"team":          "admin",
				"project_count": float64(15),
				"isPrimary":     false,
				"project_details": map[string]interface{}{
					"project_name": "project1",
				},
			},
			fields: fields{
				enabled:    true,
				connection: "update.db",
			},
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"_id": "1",
					},
					Operation: utils.One,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"name": "sharad regoti",
						},
					},
				},
			},
		},
		{
			name: "update multiple document",
			want: 3,
			want1: []interface{}{
				map[string]interface{}{
					"_id":           "2",
					"name":          "jayesh",
					"project_count": float64(10),
					"team":          "admin",
					"isPrimary":     true,
					"project_details": map[string]interface{}{
						"project_name": "project2",
					},
				},
				map[string]interface{}{
					"_id":           "3",
					"name":          "noorain",
					"team":          "admin",
					"project_count": float64(52),
					"isPrimary":     true,
					"project_details": map[string]interface{}{
						"project_name": "project2",
					},
				},
				map[string]interface{}{
					"_id":           "4",
					"name":          "ali",
					"team":          "admin",
					"project_count": float64(100),
					"isPrimary":     true,
					"project_details": map[string]interface{}{
						"project_name": "project2",
					},
				},
			},
			fields: fields{
				enabled:    true,
				connection: "update.db",
			},
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"isPrimary": true,
					},
					Operation: utils.All,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"project_details": map[string]interface{}{
								"project_name": "project2",
							},
						},
					},
				},
			},
		},

		{
			name: "upsert single document which doesn't exists",
			want: 1,
			want1: []interface{}{
				map[string]interface{}{
					"_id":  "5",
					"team": "new",
					"project_details": map[string]interface{}{
						"project_name": "project4",
					},
				},
			},
			fields: fields{
				enabled:    true,
				connection: "update.db",
			},
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"team": "new",
					},
					Operation: utils.Upsert,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"_id": "5",
							"project_details": map[string]interface{}{
								"project_name": "project4",
							},
						},
					},
				},
			},
		},

		{
			name: "update only single document where the find clause is such that it has multiple documents",
			want: 1,
			want1: map[string]interface{}{
				"_id":           "2",
				"name":          "sharad regoti",
				"team":          "admin",
				"project_count": float64(10),
				"isPrimary":     true,
				"project_details": map[string]interface{}{
					"project_name": "project2",
				},
			},
			fields: fields{
				enabled:    true,
				connection: "update.db",
			},
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.UpdateRequest{
					Find: map[string]interface{}{
						"isPrimary": true,
					},
					Operation: utils.One,
					Update: map[string]interface{}{
						"$set": map[string]interface{}{
							"name": "sharad regoti",
						},
					},
				},
			},
		},
	}

	b, err := Init(true, "update.db", "bucketName", nil)
	if err != nil {
		t.Fatal("error initializing database")
	}

	if err := createDatabaseWithTestData(b); err != nil {
		log.Fatal("error test data cannot be created for executing read test")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.Update(context.Background(), tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}

			if tt.args.req.Operation == utils.Upsert {
				tt.args.req.Operation = utils.All
			}
			readCount, readResult, err := b.Read(context.Background(), tt.args.col, &model.ReadRequest{Operation: tt.args.req.Operation, Find: tt.args.req.Find})
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if readCount != tt.want {
				t.Errorf("Read() readCount = %v, want %v", readCount, tt.want)
			}
			if !reflect.DeepEqual(readResult, tt.want1) {
				t.Log(reflect.TypeOf(readResult), reflect.TypeOf(tt.want1))
				t.Errorf("Read() readResult = %v, want %v", readResult, tt.want1)
			}
		})
	}

	utils.CloseTheCloser(b)
	if err := os.Remove("update.db"); err != nil {
		t.Error("error removing database file")
	}
}
