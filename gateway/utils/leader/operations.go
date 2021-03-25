package leader

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// GetLeaderNodeID gets leader node id
func (s *Module) GetLeaderNodeID(ctx context.Context) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	value, err := s.pubsubClient.GetKey(ctx, leaderElectionRedisKey)
	if err != nil {
		return "", err
	}

	return value, nil
}

// IsLeader checks if the provided node id is the leader gateway or not
func (s *Module) IsLeader(ctx context.Context, nodeID string) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	value, err := s.pubsubClient.GetKey(ctx, leaderElectionRedisKey)
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return nodeID == value, nil
}
