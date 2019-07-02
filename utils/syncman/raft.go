package syncman

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (s *SyncManager) initRaft(nodeID string, seeds []*node) error {

	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(utils.RaftSnapshotDirectory, "raft.db"))
	if err != nil {
		return fmt.Errorf("new bolt store: %s", err)
	}
	logStore := boltDB
	stableStore := boltDB

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(utils.RaftSnapshotDirectory, 3, os.Stderr)
	if err != nil {
		return err
	}

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", ":"+s.raftPort)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(":"+s.raftPort, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)
	config.LogOutput = ioutil.Discard

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
			ID:      raft.ServerID(seed.ID),
			Address: raft.ServerAddress(seed.Addr),
		})
	}

	err = r.BootstrapCluster(raft.Configuration{Servers: servers}).Error()
	return err
}

// Apply applies a Raft log entry to the key-value store
func (s *SyncManager) Apply(l *raft.Log) interface{} {
	var c model.RaftCommand
	json.Unmarshal(l.Data, &c)

	switch c.Kind {
	case utils.RaftCommandSet:
		s.lock.Lock()
		s.projectConfig = c.Project
		s.projectID = c.ID
		s.lock.Unlock()

		// Write the config to file
		config.StoreConfigToFile(c.Project, s.configFile)

		s.cb(s.projectConfig)

	case utils.RaftCommandDelete:

	}
	return nil
}

// Restore stores the key-value store to a previous state.
func (s *SyncManager) Restore(rc io.ReadCloser) error {
	project := new(config.Project)
	if err := json.NewDecoder(rc).Decode(project); err != nil {
		return err
	}

	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	s.projectConfig = project
	s.projectID = project.ID
	return nil
}

// Snapshot returns a snapshot of the key-value store.
func (s *SyncManager) Snapshot() (raft.FSMSnapshot, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return &fsmSnapshot{store: s.projectConfig}, nil
}

type fsmSnapshot struct {
	store *config.Project
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
