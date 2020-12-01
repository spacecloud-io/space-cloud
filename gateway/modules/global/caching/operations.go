package caching

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (c *Cache) SetRemoteServiceKey(ctx context.Context, redisKey string, remoteServiceCacheOptions *CacheResult, cache *config.ReadCacheOptions, result interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !remoteServiceCacheOptions.IsCacheEnabled() || cache == nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Set remote service cache, user hasn't specified to cache the request", nil)
		return nil
	}

	if cache.TTL <= 0 {
		cache.TTL = int64(c.config.DefaultTTL)
	}

	data, _ := json.Marshal(result)
	return c.set(ctx, redisKey, cache, string(data))
}

func (c *Cache) GetRemoteService(ctx context.Context, projectID, serviceID, endpoint string, cache *config.ReadCacheOptions, cacheOptions []interface{}) (*CacheResult, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cacheResult := new(CacheResult)
	if !c.config.Enabled || cache == nil {
		return cacheResult, nil
	}

	key, isCacheHit, result, err := c.get(ctx, c.generateRemoteServiceKey(projectID, serviceID, endpoint, cacheOptions))
	if err != nil {
		return cacheResult, err
	}

	cacheResult.redisKey = key
	cacheResult.isCacheHit = isCacheHit
	cacheResult.isCacheEnabled = true

	if !isCacheHit {
		return cacheResult, err
	}

	var v interface{}
	if err := json.Unmarshal(result, &v); err != nil {
		return cacheResult, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to json unmarshal result of remote service key", err, map[string]interface{}{"key": key})
	}
	cacheResult.result = v

	return cacheResult, nil
}

func (c *Cache) GetIngressRoute(ctx context.Context, routeID string, cacheOptions []interface{}) (string, bool, *model.CacheIngressRoute, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.config.Enabled {
		return "", false, nil, nil
	}

	key, isCacheHit, result, err := c.get(ctx, c.generateIngressRoutingKey(routeID, cacheOptions))
	if err != nil {
		return "", false, nil, err
	}

	if !isCacheHit {
		return key, isCacheHit, new(model.CacheIngressRoute), nil
	}

	obj := new(model.CacheIngressRoute)
	if err := json.Unmarshal(result, obj); err != nil {
		return "", false, nil, err
	}

	return key, isCacheHit, obj, err
}

func (c *Cache) SetIngressRouteKey(ctx context.Context, redisKey string, cache *config.ReadCacheOptions, result *model.CacheIngressRoute) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if cache.TTL <= 0 {
		cache.TTL = int64(c.config.DefaultTTL)
	}

	data, _ := json.Marshal(result)
	return c.set(ctx, redisKey, cache, string(data))
}

func (c *Cache) ConnectionState(ctx context.Context) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.redisClient.Ping(ctx).Err() == nil
}

func (c *Cache) PurgeCache(ctx context.Context, projectID string, req *model.CachePurgeRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.config.Enabled {
		return nil
	}

	if projectID == "*" {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot purge all the projects at once", nil, map[string]interface{}{"projectID": projectID, "requestObj": req})
	}

	prefixKey := ""
	switch req.Resource {
	case config.ResourceRemoteService:
		if req.ServiceId != "*" && req.ID != "*" {
			prefixKey = c.generateRemoteServiceEndpointPrefixKey(projectID, req.DbAlias, req.ID)
		} else if req.ServiceId != "*" {
			prefixKey = c.generateRemoteServicePrefixKey(projectID, req.ServiceId)
		} else {
			prefixKey = c.generateRemoteServiceResourcePrefixKey(projectID)
		}
	case config.ResourceDatabaseSchema:
		if req.DbAlias != "*" && req.ID != "*" {
			prefixKey = c.generateDatabaseTablePrefixKey(projectID, req.DbAlias, req.ID)
		} else if req.DbAlias != "*" {
			prefixKey = c.generateDatabaseAliasPrefixKey(projectID, req.DbAlias)
		} else {
			prefixKey = c.generateDatabaseResourcePrefixKey(projectID)
		}
	case config.ResourceIngressRoute:
		if req.ID != "*" {
			prefixKey = c.generateIngressRoutingPrefixWithRouteID(req.ID)
		} else {
			prefixKey = c.generateIngressRoutingResourcePrefixKey()
		}
	case "*":
		prefixKey = c.clusterID + "::"
	default:
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid resource type provided (%s) for purging cache", req.Resource), nil, map[string]interface{}{"projectID": projectID, "requestObj": req})
	}
	// list & delete
	if prefixKey != "" {
		prefixKey += "*" // pattern for matching prefix in redis
		keys := c.redisClient.Keys(ctx, prefixKey)
		keysArr, err := keys.Result()
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list redis keys with prefix (%s)", prefixKey), err, map[string]interface{}{"projectID": projectID, "requestObj": req})
		}
		for _, key := range keysArr {
			if err := c.redisClient.Del(ctx, key).Err(); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to purge redis key (%s)", key), err, map[string]interface{}{"projectID": projectID, "requestObj": req})
			}
		}
	}
	return nil
}
