package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/spacecloud-io/space-cloud/modules/tasks"
	"go.uber.org/zap"
)

// AddTask adds a message in a redis stream. It will make a new stream if it doesn't already exist
func (s *Source) AddTask(ctx context.Context, queue string, task *tasks.Task) (taskID string, err error) {
	// First parse the payload as json
	payload, err := json.Marshal(task.Payload)
	if err != nil {
		s.logger.Error("Unable to marshal payload of task to enqueue", zap.String("queue", queue), zap.Error(err))
		return "", err
	}

	// Also parse the metadata. We can ignore its error since the type is just map[string]string
	metadata, _ := json.Marshal(task.Metadata)

	// Prepare object to enqueue
	redisMessage := map[string]any{
		"md": string(metadata),
		"pl": string(payload),
	}

	// Add the message to redis stream
	result := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream:     queue,
		ID:         task.ID,
		MaxLen:     int64(s.Spec.TaskQueueConfig.MaxQueueSize),
		NoMkStream: false,
		Approx:     true,
		Values:     redisMessage,
	})
	if result.Err() != nil {
		s.logger.Error("Unable to enqueue task in redis stream", zap.String("queue", queue), zap.Error(err))
		return "", result.Err()
	}

	return result.Val(), nil
}

// ReceiveTasks returns tasks in the queue. The number of tasks returned is capped via. the count option.
// It first attempts to claim idle tasks only then attempts to read from the stream.
func (s *Source) ReceiveTasks(ctx context.Context, queue, consumerGroup string, opts tasks.ReceiveTasksOptions) ([]tasks.Task, error) {
	// Sanitize the input options first
	if opts.Wait == 0 {
		opts.Wait = s.maxWaitTime
	}
	if opts.Count == 0 {
		opts.Count = 1
	}

	tasksReclaimed, err := s.claimPendingMsgs(ctx, queue, consumerGroup, opts)
	if err != nil {
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf("Claimed %d messages. Needed %d.", len(tasksReclaimed), opts.Count))

	// We can exit if we allready have enough tasks
	if len(tasksReclaimed) == opts.Count {
		s.logger.Debug("We have enough messgaes. Returning...")
		return tasksReclaimed, nil
	}

	// We need to get more tasks from the queue

	// Check how many additional tasks do we require
	opts.Count = opts.Count - len(tasksReclaimed)
	s.logger.Debug(fmt.Sprintf("Retrieving %d more messages.", opts.Count))

	// Dequeue those many messages
	tasksRetrived, err := s.readMsgs(ctx, queue, consumerGroup, opts)
	if err != nil {
		return nil, err
	}

	// Return a combination of reclaimed tasks along with fresh ones
	return append(tasksReclaimed, tasksRetrived...), nil
}

// AckTask acknowledges a pending messgae
func (s *Source) AckTask(ctx context.Context, queue, consumerGroup, taskID string) error {
	_, err := s.client.XAck(ctx, queue, consumerGroup, taskID).Result()
	return err
}

// DeleteTask deletes a message from the queue
func (s *Source) DeleteTask(ctx context.Context, queue, taskID string) error {
	_, err := s.client.XDel(ctx, queue, taskID).Result()
	return err
}

// GetPendingTasks returns the task infos of all the tasks in the queue
func (s *Source) GetPendingTasks(ctx context.Context, queue, consumerGroup string, count int64) ([]tasks.TaskInfo, error) {
	// Set a high value of count as default
	if count == 0 {
		count = 999999999
	}

	// Get the pending tasks
	pendingTasks, err := s.client.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: queue,
		Group:  consumerGroup,
		Start:  "-",
		End:    "+",
		Count:  count,
	}).Result()
	if err != nil {
		return nil, err
	}

	taskInfos := make([]tasks.TaskInfo, len(pendingTasks))
	for i, t := range pendingTasks {
		taskInfos[i] = tasks.TaskInfo{
			ID:          t.ID,
			IdleTime:    uint(t.Idle.Milliseconds()),
			NoOfRetires: uint(t.RetryCount),
			Consumer:    t.Consumer,
		}
	}

	return taskInfos, nil
}

// AddConsumerGroup creates a new consumer group
func (s *Source) AddConsumerGroup(ctx context.Context, queue, consumerGroup string) error {
	_, err := s.client.XGroupCreate(ctx, queue, consumerGroup, "0").Result()
	if err != nil && strings.Contains(err.Error(), "Consumer Group name already exists") {
		return nil
	}

	return err
}
