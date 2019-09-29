package syncman

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (s *SyncManager) initRaft(seeds []node) error {
	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(utils.RaftSnapshotDirectory, 3, os.Stderr)
	if err != nil {
		return err
	}

	// Create a boltdb store
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(utils.RaftSnapshotDirectory, "raft.db"))
	if err != nil {
		return err
	}
	logStore := boltDB
	stableStore := boltDB

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", s.myIP+":"+s.raftPort)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(s.myIP+":"+s.raftPort, addr, 3, 10*time.Second, ioutil.Discard)
	if err != nil {
		return err
	}

	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(s.myIP)
	config.LogOutput = ioutil.Discard
	config.SnapshotThreshold = 2

	// Instantiate the Raft systems.

	r, err := raft.NewRaft(config, s, logStore, stableStore, snapshots, transport)
	if err != nil {
		return err
	}

	// Store the raft object
	s.raft = r

	servers := []raft.Server{}
	for _, seed := range seeds {
		servers = append(servers, raft.Server{
			ID:      raft.ServerID(seed.Addr),
			Address: raft.ServerAddress(seed.Addr + ":" + utils.PortRaft),
		})
	}

	// fmt.Println("Bootstrapping cluster:", servers)
	if err := r.BootstrapCluster(raft.Configuration{Servers: servers}).Error(); err != nil {
		//r.AddVoter()
		// fmt.Println("Syncman bootsrapping error:", err)
	} else {
		if err := s.list.UserEvent(bootstrapEvent, []byte(bootstrapDone), false); err != nil {
			log.Fatal(err)
		}
	}

	c := s.raft.GetConfiguration()
	if err := c.Error(); err != nil {
		log.Println("Syncman node sync error:", err)
	}
	return nil
}

// Apply applies a Raft log entry to the key-value store
func (s *SyncManager) Apply(l *raft.Log) interface{} {

	var c model.RaftCommand
	json.Unmarshal(l.Data, &c)

	switch c.Kind {
	case utils.RaftCommandSet:

		found := false
		for i, p := range s.projectConfig.Projects {
			if p.ID == c.Project.ID {
				s.projectConfig.Projects[i] = c.Project
				found = true
			}
		}
		if !found && len(s.projectConfig.Projects) == 0 {
			s.projectConfig.Projects = append(s.projectConfig.Projects, c.Project)
		}
		// Write the config to file
		config.StoreConfigToFile(s.projectConfig, s.configFile)

		go s.cb(s.projectConfig)

	case utils.RaftCommandDelete:
		for i, p := range s.projectConfig.Projects {
			if p.ID == c.ID {
				remove(s.projectConfig.Projects, i)
				break
			}
		}
		config.StoreConfigToFile(s.projectConfig, s.configFile)

	case utils.RaftCommandSetStatic:
		if s.projectConfig.Static == nil {
			s.projectConfig.Static = &config.Static{}
		}

		s.projectConfig.Static.Routes = c.Static.Routes

		go s.cb(s.projectConfig)

		// Write the config to file
		config.StoreConfigToFile(s.projectConfig, s.configFile)
	}
	return nil
}

// Restore stores the key-value store to a previous state.
func (s *SyncManager) Restore(rc io.ReadCloser) error {
	project := new(config.Config)
	if err := json.NewDecoder(rc).Decode(project); err != nil {
		return err
	}

	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	s.projectConfig = project
	return nil
}

// Snapshot returns a snapshot of the key-value store.
func (s *SyncManager) Snapshot() (raft.FSMSnapshot, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return &fsmSnapshot{store: s.projectConfig}, nil
}

type fsmSnapshot struct {
	store *config.Config
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}

func remove(s []*config.Project, i int) []*config.Project {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
