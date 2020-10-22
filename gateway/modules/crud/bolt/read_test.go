package bolt

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	// "go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestBolt_Read(t *testing.T) {
	type fields struct {
		enabled    bool
		connection string
	}
	type args struct {
		ctx context.Context
		col string
		req *model.ReadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		want1   interface{}
		wantErr bool
	}{
		{
			name: "read single document",
			want: 1,
			want1: map[string]interface{}{
				"_id":           "1",
				"name":          "sharad",
				"team":          "admin",
				"project_count": float64(15),
				"isPrimary":     false,
				"project_details": map[string]interface{}{
					"project_name": "project1",
				},
			},
			fields: fields{
				enabled:    true,
				connection: "read.db",
			},
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"_id": "1",
					},
					Operation: utils.One,
				},
			},
		},
		{
			name: "read multiple document",
			want: 3,
			want1: []interface{}{
				map[string]interface{}{
					"_id":           "2",
					"name":          "jayesh",
					"team":          "admin",
					"project_count": float64(10),
					"isPrimary":     true,
					"project_details": map[string]interface{}{
						"project_name": "project1",
					},
				}, map[string]interface{}{
					"_id":           "3",
					"name":          "noorain",
					"team":          "admin",
					"project_count": float64(52),
					"isPrimary":     true,
					"project_details": map[string]interface{}{
						"project_name": "project1",
					},
				}, map[string]interface{}{
					"_id":           "4",
					"name":          "ali",
					"team":          "admin",
					"project_count": float64(100),
					"isPrimary":     true,
					"project_details": map[string]interface{}{
						"project_name": "project1",
					},
				}},
			fields: fields{
				enabled:    true,
				connection: "read.db",
			},
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"isPrimary": true,
					},
					Operation: utils.All,
				},
			},
		},
	}

	b, err := Init(true, "read.db", "bucketName", nil)
	if err != nil {
		t.Fatal("error initializing database")
	}

	if err := createDatabaseWithTestData(b); err != nil {
		log.Fatal("error test data cannot be created for executing read test", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := b.Read(context.Background(), tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Log(reflect.TypeOf(got1), reflect.TypeOf(tt.want1))
				t.Errorf("Read() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	utils.CloseTheCloser(b)
	if err := os.Remove("read.db"); err != nil {
		t.Error("error removing database file")
	}
}
