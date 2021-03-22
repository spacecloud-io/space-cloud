package pubsub

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Module deals with pub sub related activities
type Module struct {
	lock sync.Mutex

	// Redis client
	client *redis.Client

	// Internal variables
	projectID string
	mapping   map[string]*subscription
}

// New creates a new instance of the client
func New(projectID, conn string) (*Module, error) {
	// Set a default connection string if not provided
	if conn == "" {
		conn = "localhost:6379"
	}

	c := redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "",
		DB:       0,
	})

	// Create a temporary context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := c.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Module{client: c, projectID: projectID, mapping: map[string]*subscription{}}, nil
}

// Close closes the redis client along with the active subscriptions on it
func (m *Module) Close() {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Close all active subscriptions first
	for _, sub := range m.mapping {
		_ = sub.pubsub.Close()
	}
	m.mapping = map[string]*subscription{}

	// Close the redis client
	_ = m.client.Close()
}
