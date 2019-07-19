package syncman

import (
	"log"
	"net"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"

	"github.com/spaceuptech/space-cloud/utils"
)

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func (s *SyncManager) handleSerfEvents() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			future := s.raft.VerifyLeader()

			//fmt.Println()
			if err := future.Error(); err != nil {
				//fmt.Println("Node is a follower. Leader:", s.raft.Leader())
			} else {
				//fmt.Println("Node is leader")
			}
			//fmt.Println()

			cfuture := s.raft.GetConfiguration()

			if err := cfuture.Error(); err != nil {
				log.Fatalf("error getting config: %s", err)
			}

			// configuration := cfuture.Configuration()
			// fmt.Println()
			// fmt.Println("Raft configuration:", configuration.Servers)
			// fmt.Println()

		case ev := <-s.serfEvents:

			if memberEvent, ok := ev.(serf.MemberEvent); ok {
				if s.raft == nil {
					break
				}
				leader := s.raft.VerifyLeader()
				for _, member := range memberEvent.Members {

					// Get the address of the member
					addr := member.Addr.String()

					if memberEvent.EventType() == serf.EventMemberJoin {

						if leader.Error() == nil {
							// Add the server as a voter
							f := s.raft.AddVoter(raft.ServerID(addr), raft.ServerAddress(addr+":"+utils.PortRaft), 0, 0)
							if err := f.Error(); err != nil {
								log.Fatalf("Syncman error adding voter: %s", err)
							}
						}

					} else if memberEvent.EventType() == serf.EventMemberLeave || memberEvent.EventType() == serf.EventMemberFailed || memberEvent.EventType() == serf.EventMemberReap {

						if leader.Error() == nil {
							// Remove the server
							f := s.raft.RemoveServer(raft.ServerID(addr), 0, 0)
							if err := f.Error(); err != nil {
								log.Fatalf("error removing server: %s", err)
							}
						}
					}
				}
			} else if ev.EventType() == serf.EventQuery {
				q := ev.(*serf.Query)
				if q.Name == bootstrapEvent {
					s.lock.RLock()
					if err := q.Respond([]byte(s.bootstrap)); err != nil {
						log.Println("Syncman node query response error:", err)
					}
					s.lock.RUnlock()
				}
			} else if userEvent, ok := ev.(serf.UserEvent); ok {
				if userEvent.Name == bootstrapEvent {
					s.lock.Lock()
					s.bootstrap = string(userEvent.Payload)
					s.lock.Unlock()
				}
			}
		}
	}
}

// func (s *SyncManager) getDifferenceNodesInRaft() ([]node, []node) {
// 	list := s.list.Members()
// 	c := s.raft.GetConfiguration()
// 	if err := c.Error(); err != nil {
// 		fmt.Println("Syncman node sync error:", err)
// 		return []node{}, []node{}
// 	}

// 	raftNodes := c.Configuration().Servers

// 	fmt.Println("Syncman exisiting raft nodes:", raftNodes)

// 	missing := []node{}
// 	toBeRemoved := []node{}
// 	for _, m := range list {
// 		found := false
// 		for _, r := range raftNodes {
// 			if string(r.ID) == m.Name {
// 				if strings.Split(string(r.Address), ":")[0] != m.Addr.String() && m.Name != s.myIP {
// 					toBeRemoved = append(toBeRemoved, node{Addr: m.Addr.String() + ":" + utils.PortRaft})
// 					break
// 				}
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			missing = append(missing, node{Addr: m.Addr.String() + ":" + utils.PortRaft})
// 		}
// 	}

// 	for _, r := range raftNodes {
// 		found := false
// 		for _, m := range list {
// 			if string(r.ID) == m.Name {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			toBeRemoved = append(toBeRemoved, node{Addr: string(r.ID)})
// 		}
// 	}

// 	return missing, toBeRemoved
// }
