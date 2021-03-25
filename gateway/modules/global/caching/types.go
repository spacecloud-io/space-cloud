package caching

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// CacheResult store cache result
type CacheResult struct {
	lock           sync.RWMutex
	redisKey       string
	isCacheHit     bool
	isCacheEnabled bool
	result         interface{}
}

// GetResult gets cache result
func (d *CacheResult) GetResult() interface{} {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.result
}

// GetDatabaseResult get cached database result
func (d *CacheResult) GetDatabaseResult() *model.CacheDatabaseResult {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.result.(*model.CacheDatabaseResult)
}

// Key gets cache key
func (d *CacheResult) Key() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.redisKey
}

// IsCacheHit tells if it's a cache hit or miss
func (d *CacheResult) IsCacheHit() bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.isCacheHit
}

// IsCacheEnabled tells if the cache is enabled or not
func (d *CacheResult) IsCacheEnabled() bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.isCacheEnabled
}
