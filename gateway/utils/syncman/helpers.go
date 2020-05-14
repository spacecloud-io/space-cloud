package syncman

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/getlantern/deepcopy"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
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

			s.projectConfig.Projects = remove(s.projectConfig.Projects, i)
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

// GetGatewayIndex returns the position of th current gateway instance
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
func (s *Manager) checkIfLeaderGateway(nodeID string) bool {
	return strings.HasSuffix(nodeID, "-0")
}

func (s *Manager) getLeaderGateway() (*service, error) {
	for _, service := range s.services {
		if s.checkIfLeaderGateway(service.id) {
			return service, nil
		}
	}
	return nil, errors.New("leader gateway not found")
}
func (s *Manager) PingLeader() error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for i := 0; i <= 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		service, err := s.getLeaderGateway()
		if err != nil {
			_ = utils.LogError("Unable to ping server", err)

			// Sleep for 5 seconds before trying again
			time.Sleep(5 * time.Second)
			continue
		}

		if err := s.MakeHTTPRequest(ctx, "GET", fmt.Sprintf("http://%s/v1/config/env", service.addr), "", "", struct{}{}, &map[string]interface{}{}); err != nil {
			_ = utils.LogError("Unable to ping server", err)

			// Sleep for 5 seconds before trying again
			time.Sleep(5 * time.Second)
			continue
		}

		return nil
	}

	return errors.New("leader unavailable")
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
