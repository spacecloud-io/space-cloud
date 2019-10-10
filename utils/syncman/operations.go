package syncman

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// ResolveURL returns an url for the provided service config
func (s *Manager) ResolveURL(kind, url, scheme string) string {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	// TODO: implement integration with consul
	return fmt.Sprintf("%s://%s", scheme, url)
}

// GetAssignedSpaceCloudURL returns the space cloud url assigned for the provided token
func (s *Manager) GetAssignedSpaceCloudURL(project string, token int) string {
	// TODO: implement integration with consul to get correct SC url
	return fmt.Sprintf("http://localhost:4122/v1/api/%s/eventing/process", project)
}

// GetSpaceCloudNodeURLs returns the array of space cloud urls
func (s *Manager) GetSpaceCloudNodeURLs(project string) []string {
	// TODO: implement integration with consul to get SC urls of current cluster
	return []string{fmt.Sprintf("http://localhost:4122/v1/api/%s/realtime/process", project)}
}

// GetAssignedTokens returns the array or tokens assigned to this node
func (s *Manager) GetAssignedTokens() (start int, end int) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// myHash := hash(s.list.LocalMember().Name)
	// index := 0

	// members := memRange{}
	// for _, m := range s.list.Members() {
	// 	if m.Status == serf.StatusAlive {
	// 		members = append(members, hash(m.Name))
	// 	}
	// }
	// sort.Stable(members)

	// for i, v := range members {
	// 	if v == myHash {
	// 		index = i
	// 		break
	// 	}
	// }

	// totalMembers := len(members)
	totalMembers := 1
	index := 0
	return calcTokens(totalMembers, utils.MaxEventTokens, index)
}

// GetClusterSize returns the size of the cluster
func (s *Manager) GetClusterSize() int {
	// TODO implement this function
	return 1
}

// SetProjectConfig applies the set project config command to the raft log
func (s *Manager) SetProjectConfig(project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	s.setProjectConfig(project)
	if err := s.cb(s.projectConfig); err != nil {
		return err
	}

	return config.StoreConfigToFile(s.projectConfig, s.configFile)
}

// SetStaticConfig applies the set project config command to the raft log
func (s *Manager) SetStaticConfig(static *config.Static) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	s.projectConfig.Static = static
	if err := s.cb(s.projectConfig); err != nil {
		return err
	}

	return config.StoreConfigToFile(s.projectConfig, s.configFile)
}

// DeleteProjectConfig applies delete project config command to the raft log
func (s *Manager) DeleteProjectConfig(projectID string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	s.delete(projectID)
	if err := s.cb(s.projectConfig); err != nil {
		return err
	}

	return config.StoreConfigToFile(s.projectConfig, s.configFile)
}

// GetConfig returns the config present in the state
func (s *Manager) GetConfig(projectID string) (*config.Project, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over all projects stored
	for _, p := range s.projectConfig.Projects {
		if projectID == p.ID {
			return p, nil
		}
	}

	return nil, errors.New("given project is not present in state")
}

func calcTokens(n int, tokens int, i int) (start int, end int) {
	tokensPerMember := int(math.Ceil(float64(tokens) / float64(n)))
	start = tokensPerMember * i
	end = start + tokensPerMember - 1
	if end > tokens {
		end = tokens - 1
	}
	return
}
