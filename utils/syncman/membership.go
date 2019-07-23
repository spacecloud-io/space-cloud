package syncman

import (
	"io/ioutil"
	"strconv"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
)

func (s *SyncManager) initMembership(seeds []node) error {

	// Create a membership config and assign name
	c := memberlist.DefaultLocalConfig()
	c.Name = s.myIP
	c.LogOutput = ioutil.Discard
	// Assign the port
	portInt, err := strconv.Atoi(s.gossipPort)
	if err != nil {
		return err
	}
	c.BindPort = portInt
	c.BindAddr = "0.0.0.0"
	c.AdvertisePort = portInt
	c.AdvertiseAddr = s.myIP

	serfConfig := serf.DefaultConfig()
	serfConfig.NodeName = s.myIP
	serfConfig.EventCh = s.serfEvents
	serfConfig.MemberlistConfig = c
	serfConfig.LogOutput = ioutil.Discard

	// Create a serf list
	list, err := serf.Create(serfConfig)
	if err != nil {
		return err
	}

	array := make([]string, len(seeds))
	for i, seed := range seeds {
		array[i] = seed.Addr + ":" + s.gossipPort
	}

	// fmt.Println("Membership created:", len(array))

	// Join an existing cluster by specifying at least one known member.
	_, err = list.Join(array, false)
	if err != nil {
		return err
	}

	// fmt.Println("Membership done:", list.NumNodes())

	res, err := list.Query(bootstrapEvent, []byte{}, nil)
	if err != nil {
		return err
	}

	s.lock.Lock()
	for r := range res.ResponseCh() {
		if s.bootstrap != bootstrapDone {
			s.bootstrap = string(r.Payload)
			// fmt.Println("Syncman: Got query response", s.bootstrap)
		}
	}
	s.lock.Unlock()

	// Save member list
	s.list = list
	return nil
}

// func (s *SyncManager) handleMemberChange() {
// 	s.membersLock.Lock()
// 	defer s.membersLock.Unlock()

// 	fmt.Println("Membership: I'm the leader")
// 	time.Sleep(10 * time.Second)
// 	toBeAdded, toBeRemoved := s.getDifferenceNodesInRaft()

// 	for _, n := range toBeRemoved {
// 		fmt.Println("Membership: Removing:", n.Addr)
// 		if err := s.raft.RemoveServer(raft.ServerID(n.Addr), 0, 0).Error(); err != nil {
// 			fmt.Println("Membership: Error in adding Removing:", n.Addr, err)
// 		} else {
// 			fmt.Println("Membership: Removed:", n.Addr)
// 		}
// 	}
// 	for _, n := range toBeAdded {
// 		if err := s.raft.AddVoter(raft.ServerID(n.Addr), raft.ServerAddress(n.Addr), 0, 0).Error(); err != nil {
// 			fmt.Println("Membership: Error in adding Adding:", n.Addr, err)
// 		} else {
// 			fmt.Println("Membership: Added:", n.Addr)
// 		}
// 	}

// 	c := s.raft.GetConfiguration()
// 	if err := c.Error(); err != nil {
// 		fmt.Println("Syncman node sync error:", err)
// 		return
// 	}

// 	fmt.Println("Syncman after update raft nodes:", c.Configuration().Servers)
// }

// // NotifyJoin is called when a member leaves the gossip list
// func (s *SyncManager) NotifyJoin(n *memberlist.Node) {
// 	fmt.Println("Membership: New member in:", n.Name)
// 	if s.raft != nil {
// 		if s.raft.VerifyLeader().Error() != nil {
// 			go s.handleMemberChange()
// 		} else {
// 			fmt.Println("Membership: Not the leader")
// 		}
// 	} else {
// 		fmt.Println("Membership: raft not initialised")
// 	}
// }

// // NotifyLeave is called when a member leaves the gossip list
// func (s *SyncManager) NotifyLeave(n *memberlist.Node) {
// 	fmt.Println("Membership: Member out:", n.Name)
// 	if s.raft != nil {
// 		if s.raft.VerifyLeader().Error() != nil {
// 			go s.handleMemberChange()
// 		} else {
// 			fmt.Println("Membership: Not the leader")
// 		}
// 	} else {
// 		fmt.Println("Membership: raft not initialised")
// 	}
// }

// // NotifyUpdate is called when a member modifies something
// func (s *SyncManager) NotifyUpdate(n *memberlist.Node) {}
