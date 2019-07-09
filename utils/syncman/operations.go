package syncman

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hashicorp/raft"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// SetGlobalConfig sets the global config. This must be called before the Start command.
func (s *SyncManager) SetGlobalConfig(c *config.Config) {
	s.lock.Lock()
	s.projectConfig = c
	s.lock.Unlock()
}

// GetGlobalConfig gets the global config
func (s *SyncManager) GetGlobalConfig() *config.Config {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.projectConfig
}

func makeRequest(method, token, url string, data *bytes.Buffer) error {

	// Create the http request
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return err
	}

	// Add token header
	req.Header.Add("Authorization", "Bearer "+token)

	// Create a http client and fire the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	m := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&m)

	if resp.StatusCode != http.StatusOK {
		return errors.New(m["error"].(string))
	}

	return nil
}

// SetConfig applies the config to the raft log
func (s *SyncManager) SetConfig(token string, project *config.Project) error {
	if s.raft.State() != raft.Leader {
		// Marshal json into byte array
		data, _ := json.Marshal(project)

		// Get the raft leader addr
		addr := s.raft.Leader()

		// Make the http request
		return makeRequest("POST", token, "http://"+string(addr)+":8080/v1/api/config", bytes.NewBuffer(data))
	}

	if !s.validateConfigOp(project) {
		return errors.New("Please upgrade your instance")
	}

	// Create a raft command
	c := &model.RaftCommand{Kind: utils.RaftCommandSet, Project: project, ID: project.ID}
	data, _ := json.Marshal(c)

	// Apply the command to the raft log
	return s.raft.Apply(data, 0).Error()
}

// SetDeployConfig applies the config to the raft log
func (s *SyncManager) SetDeployConfig(token string, deploy *config.Deploy) error {
	if s.raft.State() != raft.Leader {
		// Marshal json into byte array
		data, _ := json.Marshal(deploy)

		// Get the raft leader addr
		addr := s.raft.Leader()

		// Make the http request
		return makeRequest("POST", token, "http://"+string(addr)+":8080/v1/api/config/deploy", bytes.NewBuffer(data))
	}

	// Create a raft command
	c := &model.RaftCommand{Kind: utils.RaftCommandSetDeploy, Deploy: deploy}
	data, _ := json.Marshal(c)

	// Apply the command to the raft log
	return s.raft.Apply(data, 0).Error()
}

// DeleteConfig applies the config to the raft log
func (s *SyncManager) DeleteConfig(token, projectID string) error {
	if s.raft.State() != raft.Leader {

		// Get the raft leader addr
		addr := s.raft.Leader()

		// Make the http request
		return makeRequest("DELETE", token, "http://"+string(addr)+"/v1/api/config/"+projectID, nil)
	}

	// Create a raft command
	c := &model.RaftCommand{Kind: utils.RaftCommandDelete, ID: projectID}
	data, _ := json.Marshal(c)

	// Apply the command to the raft log
	return s.raft.Apply(data, 0).Error()
}

// GetConfig returns the config present in the state
func (s *SyncManager) GetConfig(projectID string) (*config.Project, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over all projects stored
	for _, p := range s.projectConfig.Projects {
		if projectID == p.ID {
			return p, nil
		}
	}

	return nil, errors.New("Given project is not present in state")
}
