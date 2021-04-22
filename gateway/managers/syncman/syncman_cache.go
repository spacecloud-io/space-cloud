package syncman

import (
	"context"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SetCacheConfig sets the caching config
func (s *Manager) SetCacheConfig(ctx context.Context, cacheConfig *config.CacheConfig, params model.RequestParams) (int, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), err
		}

		// Gracefully return
		return hookResponse.Status(), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	if cacheConfig.DefaultTTL == 0 {
		cacheConfig.DefaultTTL = utils.DefaultCacheTTLTimeout
	}

	if err := s.modules.Caching().SetCachingConfig(ctx, cacheConfig); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set caching module", err, nil)
	}

	s.projectConfig.CacheConfig = cacheConfig

	resourceID := config.GenerateResourceID(s.clusterID, "noProject", config.ResourceCacheConfig, "cache")
	if err := s.store.SetResource(ctx, resourceID, cacheConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetCacheConfig returns the cache config stored
func (s *Manager) GetCacheConfig(ctx context.Context, params model.RequestParams) (int, interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), nil, nil
	}

	// Acquire a lock
	s.lock.RLock()
	defer s.lock.RUnlock()
	return http.StatusOK, s.projectConfig.CacheConfig, nil
}
