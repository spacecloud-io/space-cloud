package caching

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

type CacheResult struct {
	lock           sync.RWMutex
	redisKey       string
	isCacheHit     bool
	isCacheEnabled bool
	result         interface{}
}

func (d *CacheResult) GetResult() interface{} {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.result
}

func (d *CacheResult) GetDatabaseResult() *model.CacheDatabaseResult {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.result.(*model.CacheDatabaseResult)
}

func (d *CacheResult) Key() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.redisKey
}

func (d *CacheResult) IsCacheHit() bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.isCacheHit
}

func (d *CacheResult) IsCacheEnabled() bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.isCacheEnabled
}
