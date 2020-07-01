package syncman

import (
	"fmt"
	"math"

	"github.com/getlantern/deepcopy"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func (s *Manager) setProjectConfig(conf *config.Project) {
	for i, p := range s.projectConfig.Projects {
		if p.ID == conf.ID {
			s.projectConfig.Projects[i] = conf
			return
		}
	}

	s.projectConfig.Projects = append(s.projectConfig.Projects, conf)
}

func (s *Manager) delete(projectID string) {
	for i, p := range s.projectConfig.Projects {
		if p.ID == projectID {
			remove(s.projectConfig.Projects, i)
			break
		}
	}
}

func remove(s []*config.Project, i int) []*config.Project {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

type scServices []*service

func (a scServices) Len() int           { return len(a) }
func (a scServices) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a scServices) Less(i, j int) bool { return a[i].id < a[j].id }

func calcTokens(n int, tokens int, i int) (start int, end int) {
	tokensPerMember := int(math.Ceil(float64(tokens) / float64(n)))
	start = tokensPerMember * i
	end = start + tokensPerMember - 1
	if end > tokens {
		end = tokens - 1
	}
	return
}

func calcIndex(token, totalTokens, n int) int {
	bucketSize := totalTokens / n
	return token / bucketSize
}

// GetGatewayIndex returns the position of the current gateway instance
func (s *Manager) GetGatewayIndex() int {
	index := 0

	for i, v := range s.services {
		if v.id == s.nodeID {
			index = i
			break
		}
	}
	return index
}

// getConfigWithoutLock returns the config present in the state
func (s *Manager) getConfigWithoutLock(projectID string) (*config.Project, error) {
	// Iterate over all projects stored
	for _, p := range s.projectConfig.Projects {
		if projectID == p.ID {
			proj := new(config.Project)
			if err := deepcopy.Copy(proj, p); err != nil {
				return nil, err
			}

			return proj, nil
		}
	}

	return nil, fmt.Errorf("given project (%s) is not present in state", projectID)
}

// GetNodeID returns node id assigned to sc
func (s *Manager) GetNodeID() string {
	return s.nodeID
}

// GetSpaceCloudURLFromID returns addr for corresponding nodeID
func (s *Manager) GetSpaceCloudURLFromID(nodeID string) (string, error) {
	for _, service := range s.services {
		if nodeID == service.id {
			return service.addr, nil
		}
	}
	return "", fmt.Errorf("service with specified nodeID doesn't exists")
}
