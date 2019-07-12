package syncman

import (
	"log"
	"sync"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/deploy"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/projects"
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
	projects      *projects.Projects
	deploy        *deploy.Module
	adminMan      *admin.Manager
}

type node struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

// New creates a new instance of the sync manager
func New(projects *projects.Projects, d *deploy.Module, adminMan *admin.Manager) *SyncManager {
	// Create a SyncManger instance
	s := new(SyncManager)
	s.deploy = d
	s.projects = projects
	return s
}

// Start begins the sync manager operations
func (s *SyncManager) Start(nodeID, configFilePath, gossipPort, raftPort string, seeds []string) error {
	// Save the ports
	s.lock.Lock()
	s.gossipPort = gossipPort
	s.raftPort = raftPort

	s.configFile = configFilePath
	if s.projectConfig.NodeID == "" {
		s.projectConfig.NodeID = nodeID
	}
	// Write the config to file
	config.StoreConfigToFile(s.projectConfig, s.configFile)

	if len(s.projectConfig.Projects) > 0 {
		for _, p := range s.projectConfig.Projects {
			if err := s.projects.StoreProject(p); err != nil {
				log.Println("Load Project Error: ", err)
			}
		}
	}

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

// ClusterSize returns the size of the member list
func (s *SyncManager) ClusterSize() int {
	return s.list.NumMembers()
}
