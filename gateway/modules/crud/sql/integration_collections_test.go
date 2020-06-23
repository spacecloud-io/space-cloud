// +build integration

package sql

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
				{TableName: "event_logs"},
				{TableName: "invocation_logs"},
				{TableName: "orders"},
				{TableName: "raw_batch"},
				{TableName: "raw_query"},
			},
			wantErr: false,
		},
	}

	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("GetCollections() Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetCollections(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCollections() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCollections() got = %v, want %v", got, tt.want)
			}
		})
	}
}
