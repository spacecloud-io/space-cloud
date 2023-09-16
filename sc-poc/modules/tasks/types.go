package tasks

import (
	"context"
	"time"
)

type (
	// Task describes the definition of a task.
	Task struct {
		ID       string            `json:"id,omitempty"`
		Metadata map[string]string `json:"metadata,omitempty"`
		Payload  any               `json:"payload"`
	}

	// TaskInfo describes a task's attributes.
	TaskInfo struct {
		ID          string `json:"id,omitempty"`
		IdleTime    uint   `json:"idleTime"`
		NoOfRetires uint   `json:"noOfRetries"`
		Consumer    string `json:"consumer"`
	}

	// ReceiveTasksOptions describes the options to pass to the dequeue operation
	ReceiveTasksOptions struct {
		// Count specifies the maximum number of tasks to dequeue. The value "0" indicates "1".
		Count int `json:"count"`

		// Wait specifies the maximum amount of time to wait to enqueue "Count" nuumber of tasks.
		// ReadTasks operation may return fewer than "Count" tasks include returning 0 tasks.
		//
		// The default value for wait varies based on implementation of source.
		Wait time.Duration `json:"wait"`
	}

	// TaskQueueSource describes the interface a source much implement to become a taskqueue source
	TaskQueueSource interface {
		AddTask(ctx context.Context, queue string, task *Task) (taskID string, err error)
		ReceiveTasks(ctx context.Context, queue, consumerGroup string, opts ReceiveTasksOptions) (tasks []Task, err error)
		AckTask(ctx context.Context, queue, consumerGroup, taskID string) error
		DeleteTask(ctx context.Context, queue, taskID string) error
		GetPendingTasks(ctx context.Context, queue, consumerGroup string, count int64) (tasks []TaskInfo, err error)
		AddConsumerGroup(ctx context.Context, queue, consumerGroup string) error
	}
)

// Types for the http requests and responses

type (
	// AddTaskRequest describes the request body for the "AddTask" operation.
	AddTaskRequest struct {
		Task
	}

	// AddTaskResponse describes the response body for the "AddTask" operation.
	AddTaskResponse struct {
		TaskID string `json:"taskId"`
	}

	// ReceiveTasksRequest describes the request body for the "ReadTasks" operation.
	ReceiveTasksRequest struct {
		ConsumerGroup string              `json:"consumerGroup"`
		Options       ReceiveTasksOptions `json:"options"`
	}

	// ReceiveTasksResponse describes the response body for the "ReadTasks" operation.
	ReceiveTasksResponse struct {
		Tasks []Task `json:"tasks"`
	}

	// AckTaskRequest describes the request body for the "AckTask" operation.
	AckTaskRequest struct {
		ConsumerGroup string `json:"consumerGroup"`
		TaskID        string `json:"taskId"`
	}

	// DeleteTaskRequest describes the request body for the "DeleteTask" operation.
	DeleteTaskRequest struct {
		TaskID string `json:"taskId"`
	}

	// GetPendingTasksRequest describes the request body for the "GetPendingTasks" operation.
	GetPendingTasksRequest struct {
		ConsumerGroup string `json:"consumerGroup"`
		Count         int64  `json:"count"`
	}

	// GetPendingTasksResponse describes the response body for the "GetPendingTasks" operaiton.
	GetPendingTasksResponse struct {
		Tasks []TaskInfo `json:"tasks"`
	}

	// AddConsumerGroupRequest describes the request body of the "AddConsumerGroup" operation.
	AddConsumerGroupRequest struct {
		ConsumerGroup string `json:"consumerGroup"`
	}
)
