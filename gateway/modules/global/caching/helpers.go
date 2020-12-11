package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

const (
	databaseJoinTypeResult = "result"
	databaseJoinTypeJoin   = "join"
	databaseJoinTypeAlways = "always"

	keyTypeTTL        = "ttl"
	keyTypeInvalidate = "invalidate"
)

func removeTablePrefixInWhereClauseFields(rootTableName string, whereClause map[string]interface{}) bool {
	for key, value := range whereClause {
		if arr := strings.Split(key, "."); len(arr) == 2 {
			// Return true if a where clause has been put in nested table
			if arr[0] != rootTableName {
				return true
			}

			// Remove og key with column key
			delete(whereClause, key)
			whereClause[arr[1]] = value
		}

		newObj, ok := value.(map[string]interface{})
		if ok {
			if removeTablePrefixInWhereClauseFields(rootTableName, newObj) {
				return true
			}
		}
	}

	return false
}

/*
Key Formats:
	OgKey/ResultKey => clusterID::projectID::db-schema-resource::DbAlias::Col::keyType::operationType::whereclause::readoptions::aggregate,postprocess,matchWhere => len of key (11)
	HalfJoinKey => clusterID::projectID::db-schema-resource::DbAlias:Col::operationType::columnName => len of key (7)
	FullJoinKey => HalfJoinKey:::OgKey/ResultKey
*/

// Database keys
func (c *Cache) generateDatabaseResultKey(projectID, dbAlias, tableName, keyType string, req *model.ReadRequest) string {
	whereClause, _ := json.Marshal(req.Find)
	readOptions, _ := json.Marshal(req.Options)
	aggregate, _ := json.Marshal(req.Aggregate)
	matchWhere, _ := json.Marshal(req.MatchWhere)
	return fmt.Sprintf("%s::%s::%s::%s::%s::%s::%s", c.generateDatabaseTablePrefixKey(projectID, dbAlias, tableName), keyType, databaseJoinTypeResult, whereClause, readOptions, aggregate, matchWhere)
}

// the prefix parameter should include tableName::operationType::columnName
func (c *Cache) generateFullDatabaseJoinKey(projectID, dbAlias, prefix, keyType, ogKey string) string {
	return c.generateHalfDatabaseJoinKey(projectID, dbAlias, prefix, keyType) + ":::" + ogKey
}

// the prefix parameter should include tableName::operationType::columnName
func (c *Cache) generateHalfDatabaseJoinKey(projectID, dbAlias, prefix, keyType string) string {
	arr := strings.Split(prefix, "::")
	temp1 := []string{arr[0], keyType}
	temp1 = append(temp1, arr[1:]...)
	return fmt.Sprintf("%s::%s", c.generateDatabaseAliasPrefixKey(projectID, dbAlias), strings.Join(temp1, "::"))
}

func (c *Cache) generateDatabaseTablePrefixKey(projectID, dbAlias, tableName string) string {
	return fmt.Sprintf("%s::%s", c.generateDatabaseAliasPrefixKey(projectID, dbAlias), tableName)
}

func (c *Cache) generateDatabaseAliasPrefixKey(projectID, dbAlias string) string {
	return fmt.Sprintf("%s::%s", c.generateDatabaseResourcePrefixKey(projectID), dbAlias)
}

func (c *Cache) generateDatabaseResourcePrefixKey(projectID string) string {
	return fmt.Sprintf("%s::%s::%s", c.clusterID, projectID, config.ResourceDatabaseSchema)
}

func (c *Cache) splitDatabaseOGKey(ctx context.Context, ogKey string) (clusterID, projectID, resourceType, dbAlias, col, keyType, joinOpType string, whereClause map[string]interface{}, readOptions *model.ReadOptions, err error) {
	realDbKeyArr := strings.Split(ogKey, "::")
	if len(realDbKeyArr) < 8 {
		return "", "", "", "", "", "", "", nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid redis result database key provided", nil, map[string]interface{}{"key": ogKey})
	}

	if err := json.Unmarshal([]byte(realDbKeyArr[7]), &whereClause); err != nil {
		return "", "", "", "", "", "", "", nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to un marshal where clause from redis result database key", err, map[string]interface{}{"key": ogKey})
	}
	if err := json.Unmarshal([]byte(realDbKeyArr[8]), &readOptions); err != nil {
		return "", "", "", "", "", "", "", nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to un marshal read options from redis result database key", err, map[string]interface{}{"key": ogKey})
	}
	return realDbKeyArr[0], realDbKeyArr[1], realDbKeyArr[2], realDbKeyArr[3], realDbKeyArr[4], realDbKeyArr[5], realDbKeyArr[6], whereClause, readOptions, nil
}

func (c *Cache) splitHalfDatabaseJoinKey(ctx context.Context, halfJoinKey string) (clusterID, projectID, resourceType, dbAlias, col, keyType, joinOpType, columnName string, err error) {
	realDbKeyArr := strings.Split(halfJoinKey, "::")
	if len(realDbKeyArr) < 7 {
		return "", "", "", "", "", "", "", "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid redis half database join key key provided", nil, map[string]interface{}{"key": halfJoinKey})
	}

	return realDbKeyArr[0], realDbKeyArr[1], realDbKeyArr[2], realDbKeyArr[3], realDbKeyArr[4], realDbKeyArr[5], realDbKeyArr[6], realDbKeyArr[7], nil
}

func (c *Cache) getOgKeyFromFullJoinKey(fullJoinKey string) string {
	arr := strings.Split(fullJoinKey, ":::")
	if len(arr) == 2 {
		return arr[1]
	}
	return ""
}

func (c *Cache) splitFullDatabaseKey(ctx context.Context, redisKey string) (isJoinKey bool, clusterID, projectID, dbAlias, col, keyType, joinOpType, columnName string, whereClause map[string]interface{}, readOptions *model.ReadOptions, err error) {
	arr := strings.Split(redisKey, ":::")
	if len(arr) == 2 { // FullJoinKey
		_, _, _, _, col, keyType, joinOpType, columnName, err = c.splitHalfDatabaseJoinKey(ctx, arr[0])
		if err != nil {
			return false, "", "", "", "", "", "", "", nil, nil, err
		}

		clusterID, projectID, _, dbAlias, _, _, _, whereClause, _, err = c.splitDatabaseOGKey(ctx, arr[1])
		if err != nil {
			return false, "", "", "", "", "", "", "", nil, nil, err
		}
		return true, clusterID, projectID, dbAlias, col, keyType, joinOpType, columnName, whereClause, readOptions, nil
	} else if len(arr) == 1 { // ogKey
		clusterID, projectID, _, dbAlias, col, keyType, joinOpType, whereClause, readOptions, err = c.splitDatabaseOGKey(ctx, arr[0])
		if err != nil {
			return false, "", "", "", "", "", "", "", nil, nil, err
		}
		return false, clusterID, projectID, dbAlias, col, keyType, joinOpType, "none", whereClause, readOptions, nil
	}
	return false, "", "", "", "", "", "", "", nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid redis database key provided", nil, map[string]interface{}{"key": redisKey})
}

// Remote service keys
func (c *Cache) generateRemoteServiceKey(projectID, serviceID, endpoint string, cacheOptions []interface{}) string {
	data, _ := json.Marshal(cacheOptions)
	return fmt.Sprintf("%s::%s", c.generateRemoteServiceEndpointPrefixKey(projectID, serviceID, endpoint), string(data))
}

func (c *Cache) generateRemoteServiceEndpointPrefixKey(projectID, serviceID, endpoint string) string {
	return fmt.Sprintf("%s::%s", c.generateRemoteServicePrefixKey(projectID, serviceID), endpoint)
}

func (c *Cache) generateRemoteServicePrefixKey(projectID, serviceID string) string {
	return fmt.Sprintf("%s::%s", c.generateRemoteServiceResourcePrefixKey(projectID), serviceID)
}

func (c *Cache) generateRemoteServiceResourcePrefixKey(projectID string) string {
	return fmt.Sprintf("%s::%s::%s", c.clusterID, projectID, config.ResourceRemoteService)
}

// Ingress keys
func (c *Cache) generateIngressRoutingKey(routeID string, cacheOptions []interface{}) string {
	data, _ := json.Marshal(cacheOptions)
	return fmt.Sprintf("%s::%s", c.generateIngressRoutingPrefixWithRouteID(routeID), string(data))
}

func (c *Cache) generateIngressRoutingPrefixWithRouteID(routeID string) string {
	return fmt.Sprintf("%s::%s", c.generateIngressRoutingResourcePrefixKey(), routeID)
}

func (c *Cache) generateIngressRoutingResourcePrefixKey() string {
	return fmt.Sprintf("%s::%s", c.clusterID, config.ResourceIngressRoute)
}

// helpers

func (c *Cache) isCachingEnabledForTable(ctx context.Context, projectID, dbAlias, col string) bool {
	rules, ok := c.dbRules[projectID]
	if !ok {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Unknown project (%s) provided, caching not enabled for table (%s)", projectID, col), map[string]interface{}{"dbAlias": dbAlias, "col": col})
		return false
	}

	rule, ok := rules[config.GenerateResourceID(c.clusterID, projectID, config.ResourceDatabaseRule, dbAlias, col, "rule")]
	if !ok {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Caching is disabled for table (%s) with database alias (%s) in project (%s))", col, dbAlias, projectID), map[string]interface{}{"dbAlias": dbAlias, "col": col})
		return false
	}
	return rule.EnableCacheInvalidation
}

func (c *Cache) instantInvalidationDelete(ctx context.Context, projectID, dbAlias, ogKey string) error {
	_, _, _, _, _, _, _, _, readOptions, err := c.splitDatabaseOGKey(ctx, ogKey)
	if err != nil {
		return err
	}

	if readOptions != nil {
		joinKeysMapping := make(map[string]map[string]string)
		utils.ExtractJoinInfoForInstantInvalidate(readOptions.Join, joinKeysMapping)

		// delete all the join keys
		for intermediateJoinKey := range joinKeysMapping {
			fullJoinKey := c.generateFullDatabaseJoinKey(projectID, dbAlias, intermediateJoinKey, keyTypeInvalidate, ogKey)
			if err := c.redisClient.Del(ctx, fullJoinKey).Err(); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete redis full join database key (%s) for instant invalidation", fullJoinKey), err, map[string]interface{}{"dbAlias": dbAlias, "projectId": projectID, "resultKey": ogKey})
			}
		}
	}

	// delete the result key
	if err := c.redisClient.Del(ctx, ogKey).Err(); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete redis database result key (%s) for instant invalidation", ogKey), err, map[string]interface{}{"dbAlias": dbAlias, "projectId": projectID})
	}

	return nil
}

func (c *Cache) startRedisScanner(ctx context.Context, receiver chan []string, closer chan struct{}, pattern string) {
	var nextCursor uint64 = 0
	var keysArr []string
	var err error

	for {
		scan := c.redisClient.Scan(ctx, nextCursor, pattern, 20)
		keysArr, nextCursor, err = scan.Result()
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list redis keys with prefix (%s)", pattern), err, map[string]interface{}{})
			closer <- struct{}{}
			break
		}
		receiver <- keysArr
		if nextCursor == 0 {
			time.Sleep(5 * time.Second)
			closer <- struct{}{}
			break
		}
	}
}
