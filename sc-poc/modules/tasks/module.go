package tasks

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/provider"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/sources/tasks"
)

func init() {
	caddy.RegisterModule(Module{})
	provider.Register("tasks", 0)
}

// Module describes the state of the tasks module
type Module struct {
	Workspace string `json:"workspace"`

	// For internal use
	logger           *zap.Logger
	taskQueueSources map[string]TaskQueueSource
	taskQueues       []*tasks.Queue
	apis             apis.APIs
}

// CaddyModule returns the Caddy module information.
func (Module) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "provider.tasks",
		New: func() caddy.Module { return new(Module) },
	}
}

// Provision sets up the auth module.
func (m *Module) Provision(ctx caddy.Context) error {
	// Get the logger
	m.logger = ctx.Logger(m).With(zap.String("workspace", m.Workspace))

	// Initialise internal variables
	m.taskQueueSources = make(map[string]TaskQueueSource)

	// Get all the dependencies
	sourceManT, _ := ctx.App("source")
	sourceMan := sourceManT.(*source.App)

	// First get all the task queue sources
	for _, s := range sourceMan.GetSources(m.Workspace, "tasks") {
		v, ok := s.(TaskQueueSource)
		if !ok {
			continue
		}
		// Store the task queue for future reference
		m.taskQueueSources[s.GetName()] = v
	}

	// Then we get the task queues
	for _, s := range sourceMan.GetSources(m.Workspace, "tasks") {
		v, ok := s.(*tasks.Queue)
		if !ok {
			continue
		}

		// Check if the task queue has a corresponding source.
		// TODO: Have the sources they target automatically create the queue
		// For eg. Create a queue object in SQS.
		if _, p := m.taskQueueSources[v.Spec.Source]; !p {
			m.logger.Error("No source found for taskqueue",
				zap.String("source", v.Spec.Source),
				zap.String("taskqueue", v.GetName()))
			continue
		}

		// Add the task queue object for future reference
		m.taskQueues = append(m.taskQueues, v)
	}

	// Add apis for each task queue
	m.prepareAPIs()
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*Module)(nil)
	_ apis.App          = (*Module)(nil)
)
