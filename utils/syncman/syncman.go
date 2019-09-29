package syncman

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
)

// SyncManager syncs the project config between folders
type SyncManager struct {
	lock          sync.RWMutex
	membersLock   sync.Mutex
	raft          *raft.Raft
	projectConfig *config.Config
	configFile    string
	gossipPort    string
	raftPort      string
	list          *serf.Serf
	cb            func(*config.Config) error
	adminMan      *admin.Manager
	myIP          string
	serfEvents    chan serf.Event
	bootstrap     string
}

const (
	bootstrapPending string = "pending"
	bootstrapDone    string = "done"
	bootstrapEvent   string = "bootstrap"
)

type node struct {
	Addr string `json:"addr"`
}

// New creates a new instance of the sync manager
func New(adminMan *admin.Manager) *SyncManager {
	// Create a SyncManger instance
	return &SyncManager{adminMan: adminMan, myIP: getOutboundIP(), serfEvents: make(chan serf.Event, 16), bootstrap: bootstrapPending}
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

	if len(s.projectConfig.Projects) > 0 {
		cb(s.projectConfig)
	}

	s.lock.Unlock()

	nodes := []node{}
	for _, m := range seeds {
		if m == "127.0.0.1" {
			m = s.myIP
		}
		addrs, err := net.LookupHost(m)
		if err != nil {
			log.Printf("Syncman: Cant look up host %s error %v", m, err)
			continue
		}
		nodes = append(nodes, node{Addr: addrs[0]})
	}

	go s.handleSerfEvents()

	// Start the membership protocol
	if err := s.initMembership(nodes); err != nil {
		return err
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

// SetProjectConfig applies the set project config command to the raft log
func (s *SyncManager) SetProjectConfig(token string, project *config.Project) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.raft.VerifyLeader().Error() != nil {
		// Marshal json into byte array
		data, _ := json.Marshal(project)

		// Get the raft leader addr
		addr := strings.Split(string(s.raft.Leader()), ":")[0]

		// Create the http request
		req, err := http.NewRequest("POST", "http://"+string(addr)+":4122/v1/api/config/projects", bytes.NewBuffer(data))
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
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("Syncman Error:", m, resp.StatusCode)
			return errors.New("Operation failed")
		}

		return nil
	}

	// Validate the operation
	if !s.adminMan.ValidateSyncOperation(s.projectConfig, project) {
		return errors.New("Please upgrade your instance")
	}

	// Create a raft command
	c := &model.RaftCommand{Kind: utils.RaftCommandSet, Project: project, ID: project.ID}
	data, _ := json.Marshal(c)

	// Apply the command to the raft log
	return s.raft.Apply(data, 0).Error()
}

// SetStaticConfig applies the set project config command to the raft log
func (s *SyncManager) SetStaticConfig(token string, static *config.Static) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.raft.VerifyLeader().Error() != nil {
		// Marshal json into byte array
		data, _ := json.Marshal(static)

		// Get the raft leader addr
		addr := strings.Split(string(s.raft.Leader()), ":")[0]

		// Create the http request
		req, err := http.NewRequest("POST", "http://"+string(addr)+":4122/v1/api/config/static", bytes.NewBuffer(data))
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
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("Syncman Error:", m, resp.StatusCode)
			return errors.New("Operation failed")
		}

		return nil
	}

	// Create a raft command
	c := &model.RaftCommand{Kind: utils.RaftCommandSetStatic, Static: static}
	data, _ := json.Marshal(c)

	// Apply the command to the raft log
	return s.raft.Apply(data, 0).Error()
}

// DeleteProjectConfig applies delete project config command to the raft log
func (s *SyncManager) DeleteProjectConfig(token, projectID string) error {
	// Acquire a lock
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.raft.VerifyLeader().Error() != nil {

		// Get the raft leader addr
		addr := strings.Split(string(s.raft.Leader()), ":")[0]

		// Create the http request
		req, err := http.NewRequest("DELETE", "http://"+string(addr)+":4122/v1/api/config/"+projectID, nil)
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
