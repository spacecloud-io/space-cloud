package leader

import (
	"context"
	"sync"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/utils/pubsub"
)

var leaderElectionRedisKey = "license-manager"

const leaderTime = 10 * time.Second // time for which a gateway can be a leader in a cluster

type Module struct {
	lock sync.RWMutex

	pubsubClient *pubsub.Module
	nodeID       string
	cbs          map[string]func()
}

func New(nodeID string, module *pubsub.Module) *Module {
	m := &Module{pubsubClient: module, nodeID: nodeID, cbs: map[string]func(){}}

	// Start the background routines
	go m.applyForLeaderPosition()
	go m.renewYourLeaderPosition()

	return m
}

func (s *Module) AddCallBack(id string, cb func()) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cbs[id] = cb
}

func (s *Module) RemoveCallBack(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.cbs, id)
}

func (s *Module) applyForLeaderPosition() {
	// Create a new ticker
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	s.applyForLeader()
	for {
		select {
		case <-ticker.C:
			s.applyForLeader()
		}
	}
}

func (s *Module) applyForLeader() {
	isKeySet, err := s.pubsubClient.SetKeyIfNotExists(context.Background(), leaderElectionRedisKey, s.nodeID, leaderTime)
	if err != nil {
		helpers.Logger.LogDebug("applyForLeaderPosition", "Unable to participate in leader election", map[string]interface{}{"key": leaderElectionRedisKey, "nodeId": s.nodeID})
		return
	}

	// If key has been set in redis, that means you have become the leader
	if isKeySet {
		helpers.Logger.LogDebug("applyForLeaderPosition", "Selected as leader", map[string]interface{}{"nodeId": s.nodeID})
		// execute all the callbacks
		s.lock.RLock()
		for _, cb := range s.cbs {
			go cb()
		}
		s.lock.RUnlock()
	}
}

func (s *Module) renewYourLeaderPosition() {
	// Create a new ticker
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.pubsubClient.RenewKeyTTLOnMatch(context.Background(), leaderElectionRedisKey, s.nodeID, leaderTime); err != nil {
				helpers.Logger.LogDebug("renewYourLeaderPosition", "Unable to renew leader position", map[string]interface{}{"key": leaderElectionRedisKey, "nodeId": s.nodeID})
				continue
			}
		}
	}
}
