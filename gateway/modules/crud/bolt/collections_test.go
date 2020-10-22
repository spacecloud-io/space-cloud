package bolt

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestBolt_GetCollections(t *testing.T) {

	b, err := Init(true, "bolt.db", "bucketName", nil)
	if err != nil {
		t.Fatal("error initializing database")
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		b       *Bolt
		args    args
		want    []utils.DatabaseCollections
		wantErr bool
	}{
		{
			name: "get collection occurs",
			b:    b,
			args: args{ctx: context.Background()},
			want: []utils.DatabaseCollections{{TableName: "project_details"}},
		},
	}

	if err := createDatabaseWithTestData(b); err != nil {
		log.Fatal("error test data cannot be created for executing collections test", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.GetCollections(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Bolt.GetCollections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bolt.GetCollections() = %v, want %v", got, tt.want)
			}
		})
	}
	utils.CloseTheCloser(b)
	if err := os.Remove("bolt.db"); err != nil {
		t.Error("error removing database file")
	}
}

func TestBolt_DeleteCollection(t *testing.T) {

	b, err := Init(true, "delete.db", "bucketName", nil)
	if err != nil {
		t.Fatal("error initializing database")
	}

	type args struct {
		ctx context.Context
		col string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []utils.DatabaseCollections
	}{
		{
			name: "delete collection doesn't take place",
			args: args{ctx: context.Background(), col: "invalid"},
			want: []utils.DatabaseCollections{{TableName: "project_details"}},
		},
		{
			name: "delete collection takes place",
			args: args{ctx: context.Background(), col: "project_details"},
			want: []utils.DatabaseCollections{},
		},
	}

	if err := createDatabaseWithTestData(b); err != nil {
		log.Fatalf("error test data cannot be created for executing delete collection test -%v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := b.DeleteCollection(context.Background(), tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("Bolt.DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := b.GetCollections(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Bolt.DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bolt.DeleteCollection() got = %v, want %v", got, tt.want)
			}
		})
	}

	utils.CloseTheCloser(b.client)
	if err := os.Remove("delete.db"); err != nil {
		t.Error("error removing database file")
	}
}
