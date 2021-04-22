package caching

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func (c *Cache) get(ctx context.Context, redisKey string) (string, bool, []byte, error) {
	result, err := c.redisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil { // key not present
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Key not present in redis, it's a cache miss", map[string]interface{}{"key": redisKey})
		return redisKey, false, nil, nil
	} else if err != nil {
		return "", false, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to get key from redis", err, map[string]interface{}{"key": redisKey})
	} else {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "It's a cache hit", map[string]interface{}{"key": redisKey})
		return redisKey, true, []byte(result), nil
	}
}

func (c *Cache) set(ctx context.Context, redisKey string, cache *config.ReadCacheOptions, result interface{}) error {
	if !c.config.Enabled || cache == nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Set Cache, Caching module is disabled or user hasn't specified to cache the request", map[string]interface{}{"cache": cache})
		return nil
	}

	if cache.TTL <= 0 {
		cache.TTL = int64(c.config.DefaultTTL)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting new key in cache", map[string]interface{}{"ttl": cache.TTL, "isInstantInvalidate": cache.InstantInvalidate, "key": redisKey})
	if err := c.redisClient.Set(ctx, redisKey, result, time.Duration(cache.TTL)*time.Second).Err(); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set result in redis", err, map[string]interface{}{"key": redisKey})
	}
	return nil
}
