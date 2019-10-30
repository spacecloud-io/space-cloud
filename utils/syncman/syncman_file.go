package syncman

import (
	"errors"
	"github.com/spaceuptech/space-cloud/config"
)

func (s *Manager) SetFileStore(project string, value *config.FileStore) error {
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

	return s.setProject(projectConfig)
}

func (s *Manager) SetFileRule(project string, value *config.FileRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig, err := s.getConfigWithoutLock(project)
	if err != nil {
		return err
	}
	for _, val := range projectConfig.Modules.FileStore.Rules {
		if val.Name == value.Name {
			return errors.New("rule with name " + value.Name + " already exists")
		}
	}
	projectConfig.Modules.FileStore.Rules = append(projectConfig.Modules.FileStore.Rules, value)

	return s.setProject(projectConfig)
}

func (s *Manager) SetDeleteFileRule(project, filename string) error {
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
	return s.setProject(projectConfig)
}
