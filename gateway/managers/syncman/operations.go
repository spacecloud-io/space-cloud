package syncman

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetEventSource returns the source id for the space cloud instance
func (s *Manager) GetEventSource() string {
	return fmt.Sprintf("sc-%s", s.nodeID)
}

// GetNodeID returns node id assigned to sc
func (s *Manager) GetNodeID() string {
	return s.nodeID
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

// GetAssignedSpaceCloudURL returns the space cloud id assigned for the provided token
func (s *Manager) GetAssignedSpaceCloudID(ctx context.Context, project string, token int) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	index := calcIndex(token, utils.MaxEventTokens, len(s.services))

	return s.services[index].id, nil
}

// GetSpaceCloudNodeIDs returns the array of space cloud ids
func (s *Manager) GetSpaceCloudNodeIDs(project string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	ids := make([]string, len(s.services))

	for i, svc := range s.services {
		ids[i] = svc.id
	}

	return ids
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

// SetClusterConfig applies the set cluster config
func (s *Manager) SetClusterConfig(ctx context.Context, req *config.ClusterConfig, params model.RequestParams) (int, error) {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()
	s.projectConfig.ClusterConfig = req
	resourceID := config.GenerateResourceID(s.clusterID, "", config.ResourceCluster, "cluster")
	if err := s.store.SetResource(ctx, resourceID, s.projectConfig.ClusterConfig); err != nil {
		return http.StatusInternalServerError, err
	}

	s.globalModules.SetMetricsConfig(s.projectConfig.ClusterConfig.EnableTelemetry)
	s.modules.LetsEncrypt().SetLetsEncryptEmail(req.LetsEncryptEmail)

	return http.StatusOK, nil
}

// GetClusterConfig returns cluster config
func (s *Manager) GetClusterConfig(ctx context.Context, params model.RequestParams) (int, *config.ClusterConfig, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return http.StatusOK, s.projectConfig.ClusterConfig, nil
}

// GetConfig returns the config present in the state
func (s *Manager) GetConfig(projectID string) (*config.ProjectConfig, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over all projects stored
	project, ok := s.projectConfig.Projects[projectID]
	if ok {
		return project.ProjectConfig, nil
	}

	return nil, errors.New("given project is not present in state")
}
