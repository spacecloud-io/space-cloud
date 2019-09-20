package syncman

import (
	"hash/fnv"
	"math"
	"sort"

	"github.com/hashicorp/serf/serf"
	"github.com/spaceuptech/space-cloud/utils"
)

func hash(value string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(value))
	return h.Sum64()
}

type memRange []uint64

func (a memRange) Len() int           { return len(a) }
func (a memRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a memRange) Less(i, j int) bool { return a[i] < a[j] }

// GetAssignedTokens returns the array or tokens assigned to this node
func (s *SyncManager) GetAssignedTokens() (start int, end int) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	myHash := hash(s.list.LocalMember().Name)
	index := 0

	members := memRange{}
	for _, m := range s.list.Members() {
		if m.Status == serf.StatusAlive {
			members = append(members, hash(m.Name))
		}
	}
	sort.Stable(members)

	for i, v := range members {
		if v == myHash {
			index = i
			break
		}
	}

	totalMembers := len(members)
	return calcTokens(totalMembers, utils.MaxEventTokens, index)
}

// GetClusterSize returns the size of the cluster
func (s *SyncManager) GetClusterSize() int {
	return s.list.NumNodes()
}

// GetAliveNodeCount returns the number of alive nodes in the cluster
func (s *SyncManager) GetAliveNodeCount() int {
	count := 0
	for _, member := range s.list.Members() {
		if member.Status == serf.StatusAlive {
			count++
		}
	}

	return count
}

func calcTokens(n int, tokens int, i int) (start int, end int) {
	tokensPerMember := int(math.Ceil(float64(tokens) / float64(n)))
	start = tokensPerMember * i
	end = start + tokensPerMember - 1
	if end > tokens {
		end = tokens - 1
	}
	return
}
