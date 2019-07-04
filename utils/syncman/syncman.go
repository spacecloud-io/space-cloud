package syncman

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// SyncManager syncs the project config between folders
type SyncManager struct {
	lock          sync.RWMutex
	raft          *raft.Raft
	projectConfig *config.Config
	configFile    string
	gossipPort    string
	raftPort      string
	list          *memberlist.Memberlist
	cb            func(*config.Config) error
}

type node struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

// New creates a new instance of the sync manager
func New() *SyncManager {
	// Create a SyncManger instance
	return new(SyncManager)
}

// Start begins the sync manager operations
func (s *SyncManager) Start(nodeID, configFilePath, gossipPort, raftPort string, seeds []string, cb func(*config.Config) error) error {
	// Save the ports
	s.lock.Lock()
	s.gossipPort = gossipPort
	s.raftPort = raftPort

	// Set the callback
	s.cb = cb

	s.configFile = configFilePath
	if s.projectConfig.NodeID == "" {
		s.projectConfig.NodeID = nodeID
	}
	// Write the config to file
	config.StoreConfigToFile(s.projectConfig, s.configFile)

	s.lock.Unlock()

	// Start the membership protocol
	if err := s.initMembership(seeds); err != nil {
		return err
	}

	nodes := []*node{}
	for _, m := range s.list.Members() {
		nodes = append(nodes, &node{ID: m.Name, Addr: m.Addr.String() + ":" + raftPort})
	}

	if err := s.initRaft(nodes); err != nil {
		return err
	}

	return nil
}

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

// SetConfig applies the config to the raft log
func (s *SyncManager) SetConfig(token string, project *config.Project) error {
	if s.raft.State() != raft.Leader {
		// Marshal json into byte array
		data, _ := json.Marshal(project)

		// Get the raft leader addr
		addr := s.raft.Leader()

		// Create the http request
		req, err := http.NewRequest("POST", "http://"+string(addr)+"/v1/api/config", bytes.NewBuffer(data))
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

	if !s.validateConfigOp(project) {
		return errors.New("Please upgrade your instance")
	}

	// Create a raft command
	c := &model.RaftCommand{Kind: utils.RaftCommandSet, Project: project, ID: project.ID}
	data, _ := json.Marshal(c)

	// Apply the command to the raft log
	return s.raft.Apply(data, 0).Error()
}

// DeleteConfig applies the config to the raft log
func (s *SyncManager) DeleteConfig(token, projectID string) error {
	if s.raft.State() != raft.Leader {

		// Get the raft leader addr
		addr := s.raft.Leader()

		// Create the http request
		req, err := http.NewRequest("DELETE", "http://"+string(addr)+"/v1/api/config/"+projectID, nil)
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

// ClusterSize returns the size of the member list
func (s *SyncManager) ClusterSize() int {
	return s.list.NumMembers()
}

func (s *SyncManager) validateConfigOp(project *config.Project) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, p := range s.projectConfig.Projects {
		if p.ID == project.ID {
			return true
		}
	}
	if len(s.projectConfig.Projects) == 0 {
		return true
	}

	return false
}
