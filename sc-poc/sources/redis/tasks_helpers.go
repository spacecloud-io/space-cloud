package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/spacecloud-io/space-cloud/modules/tasks"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

func (s *Source) claimPendingMsgs(ctx context.Context, queue, consumerGroup string, opts tasks.ReceiveTasksOptions) ([]tasks.Task, error) {
	// Prepare Redis Options
	args := redis.XAutoClaimArgs{
		Stream:  queue,
		Group:   consumerGroup,
		MinIdle: s.maxIdleTime,
		Count:   int64(opts.Count),
		Start:   "0",
	}

	// Attempt to claim pending messages
	msgs, _, err := s.client.XAutoClaim(ctx, &args).Result()
	if err != nil {
		return nil, err
	}

	return getTasksFromRedisMsgs(msgs), nil
}

func (s *Source) readMsgs(ctx context.Context, queue, consumerGroup string, opts tasks.ReceiveTasksOptions) ([]tasks.Task, error) {
	// Configure Redis options
	args := &redis.XReadGroupArgs{
		Streams:  []string{queue, ">"},
		Block:    opts.Wait,
		Count:    int64(opts.Count),
		Group:    consumerGroup,
		Consumer: utils.GetInstanceID(),
		NoAck:    false, // We want to ack the message ourselves.
	}

	// Read messages from the stream
	steams, err := s.client.XReadGroup(ctx, args).Result()
	if err != nil {
		s.logger.Error("Unable to read tasks from redis", zap.String("queue", queue), zap.Error(err))
		return nil, err
	}

	// We should have a single stream only
	return getTasksFromRedisMsgs(steams[0].Messages), nil
}

func getTasksFromRedisMsgs(msgs []redis.XMessage) []tasks.Task {
	arr := make([]tasks.Task, len(msgs))
	for i, msg := range msgs {
		// TODO: We should probably do more stric checks here
		var t tasks.Task
		_ = json.Unmarshal([]byte(msg.Values["pl"].(string)), &t.Payload)
		_ = json.Unmarshal([]byte(msg.Values["md"].(string)), &t.Metadata)
		t.ID = msg.ID

		arr[i] = t
	}

	return arr
}
