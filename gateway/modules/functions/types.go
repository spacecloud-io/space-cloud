package functions

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
	SetRemoteServiceKey(ctx context.Context, redisKey string, remoteServiceCacheOptions *caching.CacheResult, cache *config.ReadCacheOptions, result interface{}) error
	GetRemoteService(ctx context.Context, projectID, serviceID, endpoint string, cache *config.ReadCacheOptions, cacheOptions []interface{}) (*caching.CacheResult, error)
}
