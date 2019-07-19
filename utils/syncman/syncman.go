package syncman

import (
	"log"
	"net"
	"sync"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"

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
	list          *serf.Serf
	projects      *projects.Projects
	deploy        *deploy.Module
	myIP          string
	serfEvents    chan serf.Event
	bootstrap     string
	adminMan      *admin.Manager
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
func New(projects *projects.Projects, d *deploy.Module, adminMan *admin.Manager) *SyncManager {
	// Create a SyncManger instance
	return &SyncManager{adminMan: adminMan, myIP: getOutboundIP(), serfEvents: make(chan serf.Event, 16),
		bootstrap: bootstrapPending, deploy: d, projects: projects}
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

	nodes := []node{}
	for _, m := range seeds {
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
