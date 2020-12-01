package caching

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-test/deep"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestCache_generateDatabaseAliasPrefixKey(t *testing.T) {
	type args struct {
		projectID string
		dbAlias   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal Call",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
			},
			want: fmt.Sprintf("chicago::myProject::%s::db", config.ResourceDatabaseSchema),
		},
	}
	c := Init("chicago", "auto-8e55f484-288d-11eb-adc1-0242ac120002")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.generateDatabaseAliasPrefixKey(tt.args.projectID, tt.args.dbAlias)
			if res := deep.Equal(got, tt.want); res != nil {
				t.Errorf("generateDatabaseAliasPrefixKey() differences = %v", got)
			}
		})
	}
}

func TestCache_generateDatabaseResourcePrefixKey(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal Call",
			args: args{
				projectID: "myProject",
			},
			want: fmt.Sprintf("chicago::myProject::%s", config.ResourceDatabaseSchema),
		},
	}
	c := Init("chicago", "auto-8e55f484-288d-11eb-adc1-0242ac120002")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.generateDatabaseResourcePrefixKey(tt.args.projectID); got != tt.want {
				t.Errorf("generateDatabaseResourcePrefixKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_generateDatabaseResultKey(t *testing.T) {
	type args struct {
		projectID string
		dbAlias   string
		tableName string
		keyType   string
		req       *model.ReadRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Invalidate key",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
				tableName: "posts",
				keyType:   keyTypeInvalidate,
				req: &model.ReadRequest{
					GroupBy:     nil,
					Aggregate:   nil,
					Find:        nil,
					Operation:   "",
					Options:     nil,
					IsBatch:     false,
					Extras:      nil,
					PostProcess: nil,
					MatchWhere:  nil,
					Cache:       nil,
				},
			},
			want: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::null::null::null::null", config.ResourceDatabaseSchema, keyTypeInvalidate, databaseJoinTypeResult),
		},
		{
			name: "TTL key",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
				tableName: "posts",
				keyType:   keyTypeTTL,
				req: &model.ReadRequest{
					GroupBy:     nil,
					Aggregate:   nil,
					Find:        nil,
					Operation:   "",
					Options:     nil,
					IsBatch:     false,
					Extras:      nil,
					PostProcess: nil,
					MatchWhere:  nil,
					Cache:       nil,
				},
			},
			want: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::null::null::null::null", config.ResourceDatabaseSchema, keyTypeTTL, databaseJoinTypeResult),
		},
	}
	c := Init("chicago", "auto-8e55f484-288d-11eb-adc1-0242ac120002")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.generateDatabaseResultKey(tt.args.projectID, tt.args.dbAlias, tt.args.tableName, tt.args.keyType, tt.args.req); got != tt.want {
				t.Errorf("generateDatabaseResultKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_generateDatabaseTablePrefixKey(t *testing.T) {
	type args struct {
		projectID string
		dbAlias   string
		tableName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal Call",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
				tableName: "posts",
			},
			want: fmt.Sprintf("chicago::myProject::%s::db::posts", config.ResourceDatabaseSchema),
		},
	}
	c := Init("chicago", "auto-8e55f484-288d-11eb-adc1-0242ac120002")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.generateDatabaseTablePrefixKey(tt.args.projectID, tt.args.dbAlias, tt.args.tableName); got != tt.want {
				t.Errorf("generateDatabaseTablePrefixKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_generateFullDatabaseJoinKey(t *testing.T) {
	type args struct {
		projectID string
		dbAlias   string
		prefix    string
		ogKey     string
		keyType   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Invalidate key",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
				prefix:    fmt.Sprintf("posts::%s::user_id", databaseJoinTypeJoin),
				ogKey:     fmt.Sprintf("chicago::myProject::%s::db::posts::%s::null::null::null::null::null", config.ResourceDatabaseSchema, databaseJoinTypeResult),
				keyType:   keyTypeInvalidate,
			},
			want: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::user_id:::%s", config.ResourceDatabaseSchema, keyTypeInvalidate, databaseJoinTypeJoin, fmt.Sprintf("chicago::myProject::%s::db::posts::%s::null::null::null::null::null", config.ResourceDatabaseSchema, databaseJoinTypeResult)),
		},
		{
			name: "TTL key",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
				prefix:    fmt.Sprintf("posts::%s::user_id", databaseJoinTypeJoin),
				ogKey:     fmt.Sprintf("chicago::myProject::%s::db::posts::%s::null::null::null::null::null", config.ResourceDatabaseSchema, databaseJoinTypeResult),
				keyType:   keyTypeTTL,
			},
			want: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::user_id:::%s", config.ResourceDatabaseSchema, keyTypeTTL, databaseJoinTypeJoin, fmt.Sprintf("chicago::myProject::%s::db::posts::%s::null::null::null::null::null", config.ResourceDatabaseSchema, databaseJoinTypeResult)),
		},
	}
	c := Init("chicago", "auto-8e55f484-288d-11eb-adc1-0242ac120002")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.generateFullDatabaseJoinKey(tt.args.projectID, tt.args.dbAlias, tt.args.prefix, tt.args.keyType, tt.args.ogKey); got != tt.want {
				t.Errorf("generateFullDatabaseJoinKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_generateHalfDatabaseJoinKey(t *testing.T) {
	type args struct {
		projectID string
		dbAlias   string
		prefix    string
		keyType   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TTL key",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
				prefix:    fmt.Sprintf("posts::%s::user_id", databaseJoinTypeJoin),
				keyType:   keyTypeTTL,
			},
			want: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::user_id", config.ResourceDatabaseSchema, keyTypeTTL, databaseJoinTypeJoin),
		},
		{
			name: "Invalidate key",
			args: args{
				projectID: "myProject",
				dbAlias:   "db",
				prefix:    fmt.Sprintf("posts::%s::user_id", databaseJoinTypeJoin),
				keyType:   keyTypeInvalidate,
			},
			want: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::user_id", config.ResourceDatabaseSchema, keyTypeInvalidate, databaseJoinTypeJoin),
		},
	}
	c := Init("chicago", "auto-8e55f484-288d-11eb-adc1-0242ac120002")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.generateHalfDatabaseJoinKey(tt.args.projectID, tt.args.dbAlias, tt.args.prefix, tt.args.keyType); got != tt.want {
				t.Errorf("generateHalfDatabaseJoinKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_splitHalfDatabaseJoinKey(t *testing.T) {
	type args struct {
		ctx         context.Context
		halfJoinKey string
	}
	tests := []struct {
		name             string
		args             args
		wantClusterID    string
		wantProjectID    string
		wantResourceType string
		wantDbAlias      string
		wantCol          string
		wantJoinOpType   string
		wantColumnName   string
		wantKeyType      string
		wantErr          bool
	}{
		{
			name: "Invalidate key",
			args: args{
				ctx:         context.Background(),
				halfJoinKey: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::user_id", config.ResourceDatabaseSchema, keyTypeInvalidate, databaseJoinTypeJoin),
			},
			wantClusterID:    "chicago",
			wantProjectID:    "myProject",
			wantDbAlias:      "db",
			wantCol:          "posts",
			wantJoinOpType:   databaseJoinTypeJoin,
			wantColumnName:   "user_id",
			wantResourceType: string(config.ResourceDatabaseSchema),
			wantKeyType:      keyTypeInvalidate,
			wantErr:          false,
		},
		{
			name: "TTL key",
			args: args{
				ctx:         context.Background(),
				halfJoinKey: fmt.Sprintf("chicago::myProject::%s::db::posts::%s::%s::user_id", config.ResourceDatabaseSchema, keyTypeTTL, databaseJoinTypeJoin),
			},
			wantClusterID:    "chicago",
			wantProjectID:    "myProject",
			wantDbAlias:      "db",
			wantCol:          "posts",
			wantJoinOpType:   databaseJoinTypeJoin,
			wantColumnName:   "user_id",
			wantResourceType: string(config.ResourceDatabaseSchema),
			wantKeyType:      keyTypeTTL,
			wantErr:          false,
		},
	}
	c := Init("chicago", "auto-8e55f484-288d-11eb-adc1-0242ac120002")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClusterID, gotProjectID, gotResourceType, gotDbAlias, gotCol, gotKeyType, gotJoinOpType, gotColumnName, err := c.splitHalfDatabaseJoinKey(tt.args.ctx, tt.args.halfJoinKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitHalfDatabaseJoinKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKeyType != tt.wantKeyType {
				t.Errorf("splitHalfDatabaseJoinKey() gotClusterID = %v, want %v", gotClusterID, tt.wantClusterID)
			}
			if gotClusterID != tt.wantClusterID {
				t.Errorf("splitHalfDatabaseJoinKey() gotClusterID = %v, want %v", gotClusterID, tt.wantClusterID)
			}
			if gotProjectID != tt.wantProjectID {
				t.Errorf("splitHalfDatabaseJoinKey() gotProjectID = %v, want %v", gotProjectID, tt.wantProjectID)
			}
			if gotResourceType != tt.wantResourceType {
				t.Errorf("splitHalfDatabaseJoinKey() gotResourceType = %v, want %v", gotResourceType, tt.wantResourceType)
			}
			if gotDbAlias != tt.wantDbAlias {
				t.Errorf("splitHalfDatabaseJoinKey() gotDbAlias = %v, want %v", gotDbAlias, tt.wantDbAlias)
			}
			if gotCol != tt.wantCol {
				t.Errorf("splitHalfDatabaseJoinKey() gotCol = %v, want %v", gotCol, tt.wantCol)
			}
			if gotJoinOpType != tt.wantJoinOpType {
				t.Errorf("splitHalfDatabaseJoinKey() gotJoinOpType = %v, want %v", gotJoinOpType, tt.wantJoinOpType)
			}
			if gotColumnName != tt.wantColumnName {
				t.Errorf("splitHalfDatabaseJoinKey() gotColumnName = %v, want %v", gotColumnName, tt.wantColumnName)
			}
		})
	}
}
