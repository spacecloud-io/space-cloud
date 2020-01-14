package bolt

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"

	"go.etcd.io/bbolt"
)

func TestBolt_Delete(t *testing.T) {
	type fields struct {
		enabled    bool
		connection string
		client     *bbolt.DB
	}
	type args struct {
		ctx     context.Context
		project string
		col     string
		req     *model.DeleteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "delete single document",
			want: 1,
			fields: fields{
				enabled:    true,
				connection: "embedded.db",
			},
			args: args{
				ctx:     context.Background(),
				project: "gateway",
				col:     "project",
				req: &model.DeleteRequest{
					Find: map[string]interface{}{
						"_id": "1",
					},
					Operation: utils.One,
				},
			},
		},
		{
			name: "delete multiple document",
			want: 2,
			fields: fields{
				enabled:    true,
				connection: "embedded.db",
			},
			args: args{
				ctx:     context.Background(),
				project: "gateway",
				col:     "project",
				req: &model.DeleteRequest{
					Find: map[string]interface{}{
						"isPrimary": true,
					},
					Operation: utils.All,
				},
			},
		},
	}

	b, err := Init(true, "embedded.db")
	if err != nil {
		t.Fatal("error initializing database")
	}

	if err := createDatabaseWithTestData(b); err != nil {
		log.Fatal("error test data cannot be created for executing read test")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.Delete(tt.args.ctx, tt.args.project, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
	if err := os.Remove("embedded.db"); err != nil {
		t.Log("error removing database file")
	}
}
