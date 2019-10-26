package syncman

import (
	"github.com/spaceuptech/space-cloud/config"
)

func (s *Manager) SetFileStore(projectConfig *config.Project, value *config.FileStore) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig.Modules.FileStore.Enabled = value.Enabled
	projectConfig.Modules.FileStore.StoreType = value.StoreType
	projectConfig.Modules.FileStore.Conn = value.Conn
	projectConfig.Modules.FileStore.Endpoint = value.Endpoint
	projectConfig.Modules.FileStore.Bucket = value.Bucket

	return s.setProject(projectConfig)
}

func (s *Manager) SetFileRule(projectConfig *config.Project, value *config.FileRule) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	projectConfig.Modules.FileStore.Rules = append(projectConfig.Modules.FileStore.Rules, value)

	return s.setProject(projectConfig)
}

func (s *Manager) SetDeleteFileRule(projectConfig *config.Project, filename string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	temp := projectConfig.Modules.FileStore.Rules
	for i, v := range projectConfig.Modules.FileStore.Rules {
		if v.Name == filename {
			temp = append(temp[:i], temp[i+1:]...)
		}
	}
	projectConfig.Modules.FileStore.Rules = temp
	return s.setProject(projectConfig)
}
