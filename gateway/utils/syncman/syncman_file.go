package syncman

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetFileStore sets the file store module
func (s *Manager) SetFileStore(ctx context.Context, project string, value *config.FileStore) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	projectConfig.Modules.FileStore.Enabled = value.Enabled
	projectConfig.Modules.FileStore.StoreType = value.StoreType
	projectConfig.Modules.FileStore.Conn = value.Conn
	projectConfig.Modules.FileStore.Endpoint = value.Endpoint
	projectConfig.Modules.FileStore.Bucket = value.Bucket

	if err := s.modules.SetFileStoreConfig(project, projectConfig.Modules.FileStore); err != nil {
		logrus.Errorf("error setting file store config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetFileRule sets the rule for file store
func (s *Manager) SetFileRule(ctx context.Context, project string, value *config.FileRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	var doesExist bool
	for index, val := range projectConfig.Modules.FileStore.Rules {
		if val.Name == value.Name {
			projectConfig.Modules.FileStore.Rules[index] = value
			doesExist = true
		}
	}

	if !doesExist {
		projectConfig.Modules.FileStore.Rules = append(projectConfig.Modules.FileStore.Rules, value)
	}

	if err := s.modules.SetFileStoreConfig(project, projectConfig.Modules.FileStore); err != nil {
		logrus.Errorf("error setting file store config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}

// SetDeleteFileRule deletes a rule from file store
func (s *Manager) SetDeleteFileRule(ctx context.Context, project, filename string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}

	temp := projectConfig.Modules.FileStore.Rules
	for i, v := range projectConfig.Modules.FileStore.Rules {
		if v.Name == filename {
			temp = append(temp[:i], temp[i+1:]...)
			break
		}
	}
	projectConfig.Modules.FileStore.Rules = temp

	if err := s.modules.SetFileStoreConfig(project, projectConfig.Modules.FileStore); err != nil {
		logrus.Errorf("error setting file store config - %s", err.Error())
		return err
	}

	return s.setProject(ctx, projectConfig)
}
