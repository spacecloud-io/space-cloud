package pubsub

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// CheckAndSet sets the key value pair and returns if the key already existed or not
func (m *Module) CheckAndSet(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	// Flag to check if the value already existed
	exists := true

	// See if the key exists
	key = m.getTopicName(key)
	err := m.client.Get(ctx, key).Err()
	if err != nil && err != redis.Nil {
		return exists, err
	}

	// Mark key as exists
	if err == redis.Nil {
		exists = false
	}

	// Set the key value pair
	if err := m.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return exists, err
	}

	return exists, nil
}

// CheckIfKeyExists checks if the key exists
func (m *Module) CheckIfKeyExists(ctx context.Context, key string) (bool, error) {
	// See if the key exists
	key = m.getTopicName(key)
	err := m.client.Get(ctx, key).Err()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// Mark key as exists
	if err == redis.Nil {
		return false, nil
	}

	return true, nil
}
