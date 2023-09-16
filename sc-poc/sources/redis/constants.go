package redis

import (
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var (
	DefaultTaskConfig = &v1alpha1.RedisTaskQueueConfig{
		MaxQueueSize: 1000,
		CommonTaskQueueConfig: v1alpha1.CommonTaskQueueConfig{
			DefaultWaitTime: "10m",
		},
	}
)
