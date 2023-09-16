package redis

import (
	"time"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/tasks"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var gvr = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "redissources"}

func init() {
	source.RegisterSource(Source{}, gvr)
}

// Source describes the compiled redis source
type Source struct {
	v1alpha1.RedisSource

	// Internal stuff
	logger *zap.Logger `json:"-"`
	client *client     `json:"-"`

	// Storing some config here as optimization
	maxWaitTime time.Duration `json:"-"`
	maxIdleTime time.Duration `json:"-"`
}

// CaddyModule returns the Caddy module information.
func (Source) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(gvr)),
		New: func() caddy.Module { return new(Source) },
	}
}

// Provision provisions the source
func (s *Source) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger(s).With(zap.String("source", s.Name), zap.String("kind", "Redis"))

	// Get an instance of the redis client
	c, _, err := redisPool.LoadOrNew(createKey(s.Spec), createNewClient(ctx.Context, s.Spec))
	if err != nil {
		s.logger.Error("Unable to initialize redis client", zap.Error(err))
	}
	s.client = c.(*client)

	// Prepare default task config if not already provided
	if s.Spec.TaskQueueConfig == nil {
		s.Spec.TaskQueueConfig = DefaultTaskConfig
	}

	if s.Spec.TaskQueueConfig.MaxQueueSize == 0 {
		s.Spec.TaskQueueConfig.MaxQueueSize = DefaultTaskConfig.MaxQueueSize
	}

	if s.Spec.TaskQueueConfig.CommonTaskQueueConfig.DefaultWaitTime == "" {
		s.Spec.TaskQueueConfig.CommonTaskQueueConfig.DefaultWaitTime = "1m"
	}

	if s.Spec.TaskQueueConfig.CommonTaskQueueConfig.DefaultIdleTime == "" {
		s.Spec.TaskQueueConfig.CommonTaskQueueConfig.DefaultIdleTime = "3m"
	}

	// Let's parse the max wait time
	s.maxWaitTime, err = time.ParseDuration(s.Spec.TaskQueueConfig.DefaultWaitTime)
	if err != nil {
		s.logger.Error("Unable to parse wait time specified in taskqueue configuration", zap.Error(err))
		return err
	}
	s.maxIdleTime, err = time.ParseDuration(s.Spec.TaskQueueConfig.DefaultIdleTime)
	if err != nil {
		s.logger.Error("Unable to parse idle time specified in taskqueue configuration", zap.Error(err))
		return err
	}

	return nil
}

// Cleanup wraps up the source
func (s *Source) Cleanup() error {
	// Simply return if the client wasn't initialized
	if s.client == nil {
		return nil
	}

	// Remove the redis client from the pool
	_, err := redisPool.Delete(createKey(s.Spec))
	return err
}

// GetPriority returns the priority of the source.
func (s *Source) GetPriority() int {
	return 100
}

// GetProviders returns the providers this source is applicable for
func (s *Source) GetProviders() []string {
	return []string{"tasks"}
}

// Interface guards
var (
	_ caddy.Provisioner     = (*Source)(nil)
	_ caddy.CleanerUpper    = (*Source)(nil)
	_ source.Source         = (*Source)(nil)
	_ tasks.TaskQueueSource = (*Source)(nil)
)
