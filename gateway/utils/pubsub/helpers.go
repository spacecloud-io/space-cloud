package pubsub

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type subscription struct {
	ch     <-chan *redis.Message
	pubsub *redis.PubSub
}

func (m *Module) getTopicName(topic string) string {
	return fmt.Sprintf("%s-%s", m.projectID, topic)
}
