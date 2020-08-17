// +build integration

package mgo

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestSQL_GetCollections(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []utils.DatabaseCollections
		wantErr bool
	}{
		{
			name: "Get collections",
			args: args{
				ctx: context.Background(),
			},
			want: []utils.DatabaseCollections{
				{TableName: "companies"},
				{TableName: "customers"},
				{TableName: "orders"},
			},
			wantErr: false,
		},
	}

	db, err := Init(true, *connection, "myproject")
	if err != nil {
		t.Fatal("GetCollections() Couldn't establishing connection with database", dbType)
		return
	}

	// create collections
	if _, err := db.client.Database("myproject").Collection("companies").InsertOne(context.Background(), map[string]interface{}{"id": "1"}); err != nil {
		t.Log("GetCollections couldn insert data companies collection", err)
		return
	}
	if _, err := db.client.Database("myproject").Collection("orders").InsertOne(context.Background(), map[string]interface{}{"id": "1"}); err != nil {
		t.Log("GetCollections couldn insert data orders collection", err)
		return
	}
	if _, err := db.client.Database("myproject").Collection("customers").InsertOne(context.Background(), map[string]interface{}{"id": "1"}); err != nil {
		t.Log("GetCollections couldn insert data customers collection", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetCollections(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCollections() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotCount := 0
			for _, wantValue := range tt.want {
				isFound := false
				for _, gotValue := range got {
					if reflect.DeepEqual(gotValue, wantValue) {
						isFound = true
					}
				}
				if isFound {
					gotCount++
				}
			}
			if gotCount != len(tt.want) {
				t.Errorf("GetCollections() got = %v, want %v", got, tt.want)
			}
		})
	}
}
