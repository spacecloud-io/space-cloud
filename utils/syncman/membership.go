package syncman

import (
	"io/ioutil"
	"strconv"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
)

func (s *SyncManager) initMembership(seeds []string) error {

	// Create a membership config and assign name
	c := memberlist.DefaultLocalConfig()
	c.Name = s.projectConfig.NodeID
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
		if d.raft.State() == raft.Leader {
			toBeAdded, toBeRemoved := d.getDifferenceNodesInRaft()
			for _, n := range toBeAdded {
				d.raft.AddVoter(raft.ServerID(n.ID), raft.ServerAddress(n.Addr), 0, 0).Error()
			}

			for _, n := range toBeRemoved {
				d.raft.RemoveServer(raft.ServerID(n.ID), 0, 0).Error()
			}
		}
	}
}

func (d *eventDelegate) NotifyLeave(n *memberlist.Node) {
	if d.raft != nil {
		if d.raft.State() == raft.Leader {
			toBeAdded, toBeRemoved := d.getDifferenceNodesInRaft()
			for _, n := range toBeAdded {
				d.raft.AddVoter(raft.ServerID(n.ID), raft.ServerAddress(n.Addr), 0, 0).Error()
			}

			for _, n := range toBeRemoved {
				d.raft.RemoveServer(raft.ServerID(n.ID), 0, 0).Error()
			}
		}
	}
}

func (d *eventDelegate) NotifyUpdate(n *memberlist.Node) {}
