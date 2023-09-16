package tasks

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvr = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "taskqueues"}

func init() {
	source.RegisterSource(Queue{}, gvr)
}

// Queue describes the compiled redis source
type Queue struct {
	v1alpha1.TaskQueue
}

// CaddyModule returns the Caddy module information.
func (Queue) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(gvr)),
		New: func() caddy.Module { return new(Queue) },
	}
}

// GetPriority returns the priority of the source.
func (s *Queue) GetPriority() int {
	return 0
}

// GetProviders returns the providers this source is applicable for
func (s *Queue) GetProviders() []string {
	return []string{"tasks"}
}

// Interface guards
var (
	_ source.Source = (*Queue)(nil)
)
