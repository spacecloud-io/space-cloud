package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SetFileStore sets the file store module
func (s *Manager) SetFileStore(ctx context.Context, project string, value *config.FileStoreConfig, params model.RequestParams) (int, error) {
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

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	projectConfig.FileStoreConfig = value

	if err := s.modules.SetFileStoreConfig(ctx, project, projectConfig.FileStoreConfig); err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error setting file store config", err, nil)
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceFileStoreConfig, "filestore")
	if err := s.store.SetResource(ctx, resourceID, value); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetFileRule sets the rule for file store
func (s *Manager) SetFileRule(ctx context.Context, project, id string, value *config.FileRule, params model.RequestParams) (int, error) {
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

	value.ID = id
	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceFileStoreRule, id)
	if projectConfig.FileStoreRules == nil {
		projectConfig.FileStoreRules = config.FileStoreRules{resourceID: value}
	} else {
		projectConfig.FileStoreRules[resourceID] = value
	}

	if err := s.modules.SetFileStoreSecurityRuleConfig(ctx, project, projectConfig.FileStoreRules); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.SetResource(ctx, resourceID, value); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// SetDeleteFileRule deletes a rule from file store
func (s *Manager) SetDeleteFileRule(ctx context.Context, project, filename string, params model.RequestParams) (int, error) {
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

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceFileStoreRule, filename)
	delete(projectConfig.FileStoreRules, resourceID)

	if err := s.modules.SetFileStoreSecurityRuleConfig(ctx, project, projectConfig.FileStoreRules); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := s.store.DeleteResource(ctx, resourceID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// GetFileStoreConfig gets file store config
func (s *Manager) GetFileStoreConfig(ctx context.Context, project string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	// Acquire a lock
	s.lock.RLock()
	defer s.lock.RUnlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, []interface{}{projectConfig.FileStoreConfig}, nil
}

// GetFileStoreRules gets file store rules from config
func (s *Manager) GetFileStoreRules(ctx context.Context, project, ruleID string, params model.RequestParams) (int, []interface{}, error) {
	// Check if the request has been hijacked
	hookResponse := s.integrationMan.InvokeHook(ctx, params)
	if hookResponse.CheckResponse() {
		// Check if an error occurred
		if err := hookResponse.Error(); err != nil {
			return hookResponse.Status(), nil, err
		}

		// Gracefully return
		return hookResponse.Status(), hookResponse.Result().([]interface{}), nil
	}

	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(ctx, project)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if ruleID != "*" {
		resourceID := config.GenerateResourceID(s.clusterID, project, config.ResourceFileStoreRule, ruleID)
		fileRule, ok := projectConfig.FileStoreRules[resourceID]
		if ok {
			return http.StatusOK, []interface{}{fileRule}, nil
		}
		return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Security rule (%s) of file store does not exists", ruleID), fmt.Errorf("security rule not found in config"), nil)
	}

	fileRules := []interface{}{}
	for _, value := range projectConfig.FileStoreRules {
		fileRules = append(fileRules, value)
	}
	return http.StatusOK, fileRules, nil
}
