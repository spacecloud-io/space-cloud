package caching

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (c *Cache) GetDatabaseKey(ctx context.Context, projectID, dbAlias, tableName string, req *model.ReadRequest) (*CacheResult, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	cacheResult := new(CacheResult)

	// Return nil if there is no cache request
	if req.Cache == nil {
		return cacheResult, nil
	}

	// Throw error if user is trying to use caching without enabling it
	if !c.config.Enabled && req.Cache != nil {
		return cacheResult, errors.New("caching module is not enabled")
	}

	// Throw error if user is trying to use instant invalidate on disabled table
	if req.Cache.InstantInvalidate && !c.isCachingEnabledForTable(ctx, projectID, dbAlias, tableName) {
		return cacheResult, fmt.Errorf("enable instant invalidation for table - %s", tableName)
	}

	// Prepare a unique key for operation
	var redisKey string
	if !req.Cache.InstantInvalidate {
		redisKey = c.generateDatabaseResultKey(projectID, dbAlias, tableName, keyTypeTTL, req)
	} else {
		redisKey = c.generateDatabaseResultKey(projectID, dbAlias, tableName, keyTypeInvalidate, req)
	}

	// Check if result is present in cache
	key, isCacheHit, result, err := c.get(ctx, redisKey)
	if err != nil {
		return nil, err
	}

	// Update the cache result object
	cacheResult.redisKey = key
	cacheResult.isCacheHit = isCacheHit
	cacheResult.isCacheEnabled = true

	if !isCacheHit {
		return cacheResult, nil
	}

	// Unmarshal the result and return on cache hit
	v := new(model.CacheDatabaseResult)
	if err := json.Unmarshal(result, v); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to json unmarshal result of database key", err, map[string]interface{}{"key": key})
	}
	cacheResult.result = v

	return cacheResult, nil
}

func (c *Cache) SetDatabaseKey(ctx context.Context, projectID, dbAlias, col string, result *model.CacheDatabaseResult, dbCacheOptions *CacheResult, cache *config.ReadCacheOptions, cacheJoinInfo map[string]map[string]string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if !dbCacheOptions.IsCacheEnabled() || cache == nil {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Set database cache, Either caching is not enabled on table (%s) or user hasn't specified to cache the request", col), map[string]interface{}{"cache": cache})
		return nil
	}

	if cache.TTL <= 0 {
		cache.TTL = int64(c.config.DefaultTTL)
	}

	// Need to make a key for each joint table.
	for prefix, obj := range cacheJoinInfo {
		// We need to make an array here
		arr := make([]interface{}, 0)
		for key, value := range obj {
			arr = append(arr, key, value)
		}

		// Make the key for the joint table
		var fullJoinKey string
		if cache.InstantInvalidate {
			fullJoinKey = c.generateFullDatabaseJoinKey(projectID, dbAlias, prefix, keyTypeInvalidate, dbCacheOptions.redisKey)
		} else {
			fullJoinKey = c.generateFullDatabaseJoinKey(projectID, dbAlias, prefix, keyTypeTTL, dbCacheOptions.redisKey)
		}

		// Set the key and value in the cache
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting new full database join key in cache", map[string]interface{}{"ttl": cache.TTL, "isInstantInvalidate": cache.InstantInvalidate, "key": fullJoinKey})
		if err := c.redisClient.HSet(ctx, fullJoinKey, arr...).Err(); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to h set join info in redis", err, map[string]interface{}{"key": fullJoinKey})
		}

		_ = c.redisClient.Expire(ctx, fullJoinKey, time.Duration(cache.TTL)*time.Second).Err()
	}

	// Marshal the result and store it in the cache
	data, _ := json.Marshal(result)
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting new key in cache", map[string]interface{}{"ttl": cache.TTL, "isInstantInvalidate": cache.InstantInvalidate, "key": dbCacheOptions.redisKey})
	if err := c.redisClient.Set(ctx, dbCacheOptions.redisKey, string(data), time.Duration(cache.TTL)*time.Second).Err(); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set database result in redis", err, map[string]interface{}{"key": dbCacheOptions.redisKey})
	}
	return nil
}

func (c *Cache) InvalidateDatabaseCache(ctx context.Context, projectID, dbAlias, rootTable, opType string, doc map[string]interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Return if caching is not enabled for that table
	if !c.isCachingEnabledForTable(ctx, projectID, dbAlias, rootTable) {
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Instant invalidation not enabled for table (%s). Skipping invalidation event.", rootTable), nil)
		return nil
	}

	// Generate pattern to iterate over cache
	pattern := c.generateDatabaseTablePrefixKey(projectID, dbAlias, rootTable) + "::" + keyTypeInvalidate + "*"

	var nextCursor uint64 = 0
	var keysArr []string
	var err error
	for {
		scan := c.redisClient.Scan(ctx, nextCursor, pattern, 20)
		keysArr, nextCursor, err = scan.Result()
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list redis keys with prefix (%s)", pattern), err, map[string]interface{}{})
			break
		}
		for _, redisKey := range keysArr {
			_, _, _, _, _, _, joinType, _, _, _, err := c.splitFullDatabaseKey(ctx, redisKey)
			if err != nil {
				return err
			}

			switch joinType {
			case databaseJoinTypeAlways:
				fullJoinKey := redisKey
				if err := c.instantInvalidationDelete(ctx, projectID, dbAlias, c.getOgKeyFromFullJoinKey(fullJoinKey)); err != nil {
					return err
				}

			case databaseJoinTypeJoin:
				fullJoinKey := redisKey
				if opType == utils.EventDBDelete {
					if err := c.instantInvalidationDelete(ctx, projectID, dbAlias, c.getOgKeyFromFullJoinKey(fullJoinKey)); err != nil {
						return err
					}
					continue
				}

				_, _, _, _, _, _, _, columnName, _, _, err := c.splitFullDatabaseKey(ctx, fullJoinKey)
				columnValue, ok := doc[columnName]
				if !ok {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Column name (%s) not found in doc object", columnName), nil, map[string]interface{}{"fullJoinKey": fullJoinKey})
				}

				res := c.redisClient.HExists(ctx, fullJoinKey, fmt.Sprintf("%v", columnValue))
				doesExists, err := res.Result()
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to check existence of hset key (%s) having map key (%s)", fullJoinKey, columnValue), err, nil)
				}

				if doesExists {
					if err := c.instantInvalidationDelete(ctx, projectID, dbAlias, c.getOgKeyFromFullJoinKey(fullJoinKey)); err != nil {
						return err
					}
				}

			case databaseJoinTypeResult:
				ogKey := redisKey

				_, _, _, _, _, _, _, whereClause, _, err := c.splitDatabaseOGKey(ctx, ogKey)
				if err != nil {
					return err
				}

				if opType == utils.EventDBDelete {
					if err := c.instantInvalidationDelete(ctx, projectID, dbAlias, ogKey); err != nil {
						return err
					}
					continue
				}

				if removeTablePrefixInWhereClauseFields(rootTable, whereClause) || utils.Validate(whereClause, doc) {
					if err := c.instantInvalidationDelete(ctx, projectID, dbAlias, ogKey); err != nil {
						return err
					}
					continue
				}
			default:
				return err
			}
		}

		if nextCursor == 0 {
			break
		}
	}

	return nil
}
