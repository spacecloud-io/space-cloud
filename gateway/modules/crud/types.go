package crud

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/caching"
)

type integrationManagerInterface interface {
	InvokeHook(ctx context.Context, params model.RequestParams) config.IntegrationAuthResponse
}

type cachingInterface interface {
	SetDatabaseKey(ctx context.Context, projectID, dbAlias, col string, result *model.CacheDatabaseResult, dbCacheOptions *caching.CacheResult, cache *config.ReadCacheOptions, cacheJoinInfo map[string]map[string]string) error
	GetDatabaseKey(ctx context.Context, projectID, dbAlias, tableName string, req *model.ReadRequest) (*caching.CacheResult, error)
}
