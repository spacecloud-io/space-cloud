package caching

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Cache holds module config
type Cache struct {
	lock sync.RWMutex

	clusterID string
	nodeID    string

	admin *admin.Manager

	config      *config.CacheConfig
	dbRules     map[string]config.DatabaseRules // key is the project id
	redisClient *redis.Client
}

// Init creates a new instance of the cache module
func Init(clusterID, nodeID string) *Cache {
	return &Cache{clusterID: clusterID, nodeID: nodeID, config: new(config.CacheConfig), dbRules: map[string]config.DatabaseRules{}}
}

// SetCachingConfig sets caching config
func (c *Cache) SetCachingConfig(ctx context.Context, cacheConfig *config.CacheConfig) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if cacheConfig == nil {
		return nil
	}

	if err := c.admin.CheckIfCachingCanBeEnabled(ctx); err != nil {
		return err
	}

	if cacheConfig.DefaultTTL == 0 {
		cacheConfig.DefaultTTL = utils.DefaultCacheTTLTimeout
	}

	if cacheConfig.Enabled {
		// Create a new redis client
		redisClient := redis.NewClient(&redis.Options{
			Addr:     cacheConfig.Conn,
			Password: "", // no password set
			DB:       0,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot connect to redis cache", err, map[string]interface{}{"conn": cacheConfig.Conn, "ttl": cacheConfig.DefaultTTL, "isEnable": cacheConfig.Enabled})
		}

		// Store the client for future use
		c.redisClient = redisClient
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Successfully connected to redis databases", map[string]interface{}{"conn": cacheConfig.Conn, "ttl": cacheConfig.DefaultTTL, "isEnable": cacheConfig.Enabled})
	} else if c.redisClient != nil {
		// Close the client if it is present already
		_ = c.redisClient.Close()
		c.redisClient = nil
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Successfully closed redis database connection", map[string]interface{}{"conn": cacheConfig.Conn, "ttl": cacheConfig.DefaultTTL, "isEnable": cacheConfig.Enabled})
	}

	c.config = cacheConfig
	return nil
}

func (c *Cache) AddDBRules(projectID string, dbRules config.DatabaseRules) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dbRules[projectID] = dbRules
}

func (c *Cache) SetAdminModule(admin *admin.Manager) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.admin = admin
}
