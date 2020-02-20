package bolt

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestBolt_DeleteCollection(t *testing.T) {

	b, err := Init(true, "delete.db")
	if err != nil {
		t.Fatal("error initializing database")
	}

	type args struct {
		ctx     context.Context
		project string
		col     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []utils.DatabaseCollections
	}{
		{
			name: "invalid project",
			args: args{ctx: context.Background(), project: "not-gateway", col: "project"},
		},
		{
			name: "delete collection doesn't take place",
			args: args{ctx: context.Background(), project: "gateway", col: "invalid"},
			want: []utils.DatabaseCollections{{TableName: "project"}},
		},
		{
			name: "delete collection takes place",
			args: args{ctx: context.Background(), project: "gateway", col: "project"},
			want: []utils.DatabaseCollections{},
		},
	}

	if err := createDatabaseWithTestData(b); err != nil {
		log.Fatalf("error test data cannot be created for executing delete collection test -%v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := b.DeleteCollection(tt.args.ctx, tt.args.project, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("Bolt.DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	utils.CloseTheCloser(b.client)
	if err := os.Remove("delete.db"); err != nil {
		t.Error("error removing database file")
	}
}
