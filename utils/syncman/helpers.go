package syncman

import (
	"log"

	"github.com/spaceuptech/space-cloud/utils"
)

func (d *eventDelegate) getDifferenceNodesInRaft() ([]node, []node) {
	list := d.list.Members()
	c := d.raft.GetConfiguration()
	if err := c.Error(); err != nil {
		log.Println("Syncman node sync error:", err)
		return []node{}, []node{}
	}
	raftNodes := c.Configuration().Servers

	missing := []node{}
	for _, m := range list {
		found := false
		for _, r := range raftNodes {
			if string(r.ID) == m.Name {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, node{ID: m.Name, Addr: m.Addr.String() + ":" + utils.PortRaft})
		}
	}

	toBeRemoved := []node{}
	for _, r := range raftNodes {
		found := false
		for _, m := range list {
			if string(r.ID) == m.Name {
				found = true
				break
			}
		}
		if !found {
			toBeRemoved = append(toBeRemoved, node{ID: string(r.ID)})
		}
	}

	return missing, toBeRemoved
}
