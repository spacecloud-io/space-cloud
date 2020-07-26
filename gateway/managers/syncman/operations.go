package syncman

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetEventSource returns the source id for the space cloud instance
func (s *Manager) GetEventSource() string {
	return fmt.Sprintf("sc-%s", s.nodeID)
}

// GetClusterID get cluster id
func (s *Manager) GetClusterID() string {
	return s.clusterID
}

// GetNodesInCluster get total number of gateways
func (s *Manager) GetNodesInCluster() int {
	if len(s.services) == 0 {
		return 1
	}
	return len(s.services)
}

// GetAssignedSpaceCloudURL returns the space cloud url assigned for the provided token
func (s *Manager) GetAssignedSpaceCloudURL(ctx context.Context, project string, token int) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	index := calcIndex(token, utils.MaxEventTokens, len(s.services))

	return fmt.Sprintf("http://%s/v1/api/%s/eventing/process", s.services[index].addr, project), nil
}

// GetSpaceCloudNodeURLs returns the array of space cloud urls
func (s *Manager) GetSpaceCloudNodeURLs(project string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	urls := make([]string, len(s.services))

	for i, svc := range s.services {
		urls[i] = fmt.Sprintf("http://%s/v1/api/%s/realtime/process", svc.addr, project)
	}

	return urls
}

// GetRealtimeURL get the url of realtime
func (s *Manager) GetRealtimeURL(project string) string {
	return fmt.Sprintf("http://localhost:%d/v1/api/%s/realtime/handle", s.port, project)
}

// GetAssignedTokens returns the array or tokens assigned to this node
func (s *Manager) GetAssignedTokens() (start, end int) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	index := s.GetGatewayIndex()

	return calcTokens(len(s.services), utils.MaxEventTokens, index)
}

func (s *Manager) setProject(ctx context.Context, project *config.Project) error {
	s.setProjectConfig(project)
	return s.store.SetProject(ctx, project)
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
