package syncman

import (
	"io/ioutil"
	"strconv"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
)

func (s *SyncManager) initMembership(nodeID string, seeds []string) error {

	// Create a membership config and assign name
	c := memberlist.DefaultLocalConfig()
	c.Name = nodeID
	c.LogOutput = ioutil.Discard

	// Assign a custom event delegate
	c.Events = (*eventDelegate)(s)

	// Assign the port
	portInt, err := strconv.Atoi(s.gossipPort)
	if err != nil {
		return err
	}
	c.BindPort = portInt
	c.BindAddr = "0.0.0.0"
	c.AdvertisePort = portInt

	// Create a membership list
	list, err := memberlist.Create(c)
	if err != nil {
		return err
	}

	array := make([]string, len(seeds))
	for i, seed := range seeds {
		array[i] = seed + ":" + s.gossipPort
	}

	// Join an existing cluster by specifying at least one known member.
	_, err = list.Join(seeds)
	if err != nil {
		return err
	}

	// Save member list
	s.list = list
	return nil
}

type eventDelegate SyncManager

func (d *eventDelegate) NotifyJoin(n *memberlist.Node) {
	if d.raft != nil {
		d.raft.AddVoter(raft.ServerID(n.Name), raft.ServerAddress(n.Addr.String()+":"+d.raftPort), 0, 0).Error()
	}
}

func (d *eventDelegate) NotifyLeave(n *memberlist.Node) {
	d.raft.RemoveServer(raft.ServerID(n.Name), 0, 0).Error()
}

func (d *eventDelegate) NotifyUpdate(n *memberlist.Node) {}
