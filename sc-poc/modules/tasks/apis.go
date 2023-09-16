package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"github.com/spacecloud-io/space-cloud/sources/tasks"
	"github.com/spacecloud-io/space-cloud/utils"
)

// GetAPIRoutes returns the apis to be exposed
func (m *Module) GetAPIRoutes() apis.APIs {
	return m.apis
}

func (m *Module) prepareAPIs() {
	for _, tq := range m.taskQueues {
		m.apis = append(m.apis, m.prepareAddTaskRoute(tq))
		m.apis = append(m.apis, m.prepareReceiveTasksRoute(tq))
		m.apis = append(m.apis, m.prepareAckTaskRoute(tq))
		m.apis = append(m.apis, m.prepareDeleteTaskRoute(tq))
		m.apis = append(m.apis, m.prepareGetPendingTasksRoute(tq))
		m.apis = append(m.apis, m.prepareAddConsumerGroupRoute(tq))
	}
}

func (m *Module) prepareAddTaskRoute(t *tasks.Queue) *apis.API {
	// Create an open api path definition
	operation := &openapi3.Operation{
		OperationID: fmt.Sprintf("add-task-%s", t.Name),
		// TODO: Add tags based on the `space-cloud.io/tags` annotations
		Tags: []string{},
	}

	operation = apis.PrepareOpenAPIRequest(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Request payload for adding task to '%s'", t.Name),
		Ptr:         &AddTaskRequest{},
	})
	operation = apis.PrepareOpenAPIResponse(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Response payload for adding task to '%s'", t.Name),
		Ptr:         &AddTaskResponse{},
	})

	var plugins []v1alpha1.HTTPPlugin
	if t.Spec.Operations != nil && t.Spec.Operations.AddTask != nil {
		plugins = t.Spec.Operations.AddTask.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("tq-add-task-%s", t.Name),
		Path:    fmt.Sprintf("/v1/tasks/%s/add", t.Name),
		Plugins: plugins,
		OpenAPI: &apis.OpenAPI{PathDef: &openapi3.PathItem{Post: operation}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prepare context object
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
			defer cancel()

			// Prepare a temporary logger
			logger := m.logger.With(
				zap.String("source", t.Spec.Source),
				zap.String("taskqueue", t.Name),
				zap.String("operation", "add-task"),
			)

			req := new(AddTaskRequest)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				logger.Error("Unable to parse request body", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusBadRequest, err)
				return
			}

			// Get the source for the task
			s := m.taskQueueSources[t.Spec.Source]

			// Perform the "AddTask" operation
			id, err := s.AddTask(ctx, t.Name, &req.Task)
			if err != nil {
				logger.Error("Unable to add task", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			// Send task id in response
			_ = utils.SendResponse(w, http.StatusOK, AddTaskResponse{TaskID: id})
		}),
	}
}

func (m *Module) prepareReceiveTasksRoute(t *tasks.Queue) *apis.API {
	// Create an open api path definition
	operation := &openapi3.Operation{
		OperationID: fmt.Sprintf("receive-tasks-%s", t.Name),
		// TODO: Add tags based on the `space-cloud.io/tags` annotations
		Tags: []string{},
	}

	operation = apis.PrepareOpenAPIRequest(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Request payload for receiving tasks to '%s'", t.Name),
		Ptr:         &ReceiveTasksRequest{},
	})
	operation = apis.PrepareOpenAPIResponse(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Response payload for receiving tasks to '%s'", t.Name),
		Ptr:         &ReceiveTasksResponse{},
	})

	var plugins []v1alpha1.HTTPPlugin
	if t.Spec.Operations != nil && t.Spec.Operations.ReceiveTasks != nil {
		plugins = t.Spec.Operations.ReceiveTasks.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("tq-recieve-tasks-%s", t.Name),
		Path:    fmt.Sprintf("/v1/tasks/%s/receive-tasks", t.Name),
		Plugins: plugins,
		OpenAPI: &apis.OpenAPI{PathDef: &openapi3.PathItem{Post: operation}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prepare context object
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
			defer cancel()

			// Prepare a temporary logger
			logger := m.logger.With(
				zap.String("source", t.Spec.Source),
				zap.String("taskqueue", t.Name),
				zap.String("operation", "receive-tasks"),
			)

			req := new(ReceiveTasksRequest)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				logger.Error("Unable to parse request body", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusBadRequest, err)
				return
			}

			// Get the source for the task
			s := m.taskQueueSources[t.Spec.Source]

			// Perform the "ReceiveTasks" operation
			tasks, err := s.ReceiveTasks(ctx, t.Name, req.ConsumerGroup, req.Options)
			if err != nil {
				logger.Error("Unable to receive tasks", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			// Send task id in response
			_ = utils.SendResponse(w, http.StatusOK, ReceiveTasksResponse{Tasks: tasks})
		}),
	}
}

func (m *Module) prepareAckTaskRoute(t *tasks.Queue) *apis.API {
	// Create an open api path definition
	operation := &openapi3.Operation{
		OperationID: fmt.Sprintf("ack-task-%s", t.Name),
		// TODO: Add tags based on the `space-cloud.io/tags` annotations
		Tags: []string{},
	}

	operation = apis.PrepareOpenAPIRequest(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Request payload for ack task to '%s'", t.Name),
		Ptr:         &AckTaskRequest{},
	})
	operation = apis.PrepareOpenAPIResponse(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Response payload for ack task to '%s'", t.Name),
		Ptr:         nil,
	})

	var plugins []v1alpha1.HTTPPlugin
	if t.Spec.Operations != nil && t.Spec.Operations.AckTask != nil {
		plugins = t.Spec.Operations.AckTask.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("tq-ack-task-%s", t.Name),
		Path:    fmt.Sprintf("/v1/tasks/%s/ack-task", t.Name),
		Plugins: plugins,
		OpenAPI: &apis.OpenAPI{PathDef: &openapi3.PathItem{Post: operation}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prepare context object
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
			defer cancel()

			// Prepare a temporary logger
			logger := m.logger.With(
				zap.String("source", t.Spec.Source),
				zap.String("taskqueue", t.Name),
				zap.String("operation", "ack-task"),
			)

			req := new(AckTaskRequest)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				logger.Error("Unable to parse request body", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusBadRequest, err)
				return
			}

			// Get the source for the task
			s := m.taskQueueSources[t.Spec.Source]

			// Perform the "Ack Task" operation
			if err := s.AckTask(ctx, t.Name, req.ConsumerGroup, req.TaskID); err != nil {
				logger.Error("Unable to ack task", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			// Send task id in response
			w.WriteHeader(http.StatusNoContent)
		}),
	}
}

func (m *Module) prepareDeleteTaskRoute(t *tasks.Queue) *apis.API {
	// Create an open api path definition
	operation := &openapi3.Operation{
		OperationID: fmt.Sprintf("delete-task-%s", t.Name),
		// TODO: Add tags based on the `space-cloud.io/tags` annotations
		Tags: []string{},
	}

	operation = apis.PrepareOpenAPIRequest(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Request payload for delete task from '%s'", t.Name),
		Ptr:         &DeleteTaskRequest{},
	})
	operation = apis.PrepareOpenAPIResponse(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Response payload for delete task from '%s'", t.Name),
		Ptr:         nil,
	})

	var plugins []v1alpha1.HTTPPlugin
	if t.Spec.Operations != nil && t.Spec.Operations.DeleteTask != nil {
		plugins = t.Spec.Operations.DeleteTask.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("tq-delete-task-%s", t.Name),
		Path:    fmt.Sprintf("/v1/tasks/%s/delete-task", t.Name),
		Plugins: plugins,
		OpenAPI: &apis.OpenAPI{PathDef: &openapi3.PathItem{Post: operation}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prepare context object
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
			defer cancel()

			// Prepare a temporary logger
			logger := m.logger.With(
				zap.String("source", t.Spec.Source),
				zap.String("taskqueue", t.Name),
				zap.String("operation", "delete-task"),
			)

			req := new(DeleteTaskRequest)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				logger.Error("Unable to parse request body", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusBadRequest, err)
				return
			}

			// Get the source for the task
			s := m.taskQueueSources[t.Spec.Source]

			// Perform the "Delete Task" operation
			if err := s.DeleteTask(ctx, t.Name, req.TaskID); err != nil {
				logger.Error("Unable to delete task", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			// We got no response for this one
			w.WriteHeader(http.StatusNoContent)
		}),
	}
}

func (m *Module) prepareGetPendingTasksRoute(t *tasks.Queue) *apis.API {
	// Create an open api path definition
	operation := &openapi3.Operation{
		OperationID: fmt.Sprintf("get-pending-tasks-%s", t.Name),
		// TODO: Add tags based on the `space-cloud.io/tags` annotations
		Tags: []string{},
	}

	operation = apis.PrepareOpenAPIRequest(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Request payload for get pending tasks to '%s'", t.Name),
		Ptr:         &ReceiveTasksRequest{},
	})
	operation = apis.PrepareOpenAPIResponse(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Response payload for get pending tasks to '%s'", t.Name),
		Ptr:         &ReceiveTasksResponse{},
	})

	var plugins []v1alpha1.HTTPPlugin
	if t.Spec.Operations != nil && t.Spec.Operations.GetPendingTasks != nil {
		plugins = t.Spec.Operations.GetPendingTasks.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("tq-get-pending-tasks-%s", t.Name),
		Path:    fmt.Sprintf("/v1/tasks/%s/get-pending-tasks", t.Name),
		Plugins: plugins,
		OpenAPI: &apis.OpenAPI{PathDef: &openapi3.PathItem{Post: operation}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prepare context object
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
			defer cancel()

			// Prepare a temporary logger
			logger := m.logger.With(
				zap.String("source", t.Spec.Source),
				zap.String("taskqueue", t.Name),
				zap.String("operation", "get-pending-tasks"),
			)

			req := new(GetPendingTasksRequest)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				logger.Error("Unable to parse request body", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusBadRequest, err)
				return
			}

			// Get the source for the task
			s := m.taskQueueSources[t.Spec.Source]

			// Perform the "GetPendingTasks" operation
			tasks, err := s.GetPendingTasks(ctx, t.Name, req.ConsumerGroup, req.Count)
			if err != nil {
				logger.Error("Unable to get pending tasks", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			// Send task id in response
			_ = utils.SendResponse(w, http.StatusOK, GetPendingTasksResponse{Tasks: tasks})
		}),
	}
}

func (m *Module) prepareAddConsumerGroupRoute(t *tasks.Queue) *apis.API {
	// Create an open api path definition
	operation := &openapi3.Operation{
		OperationID: fmt.Sprintf("add-consumer-group-%s", t.Name),
		// TODO: Add tags based on the `space-cloud.io/tags` annotations
		Tags: []string{},
	}

	operation = apis.PrepareOpenAPIRequest(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Request payload for get pending tasks to '%s'", t.Name),
		Ptr:         &ReceiveTasksRequest{},
	})
	operation = apis.PrepareOpenAPIResponse(operation, apis.OpenAPIPayloadModifier{
		Description: fmt.Sprintf("Response payload for get pending tasks to '%s'", t.Name),
		Ptr:         &ReceiveTasksResponse{},
	})

	var plugins []v1alpha1.HTTPPlugin
	if t.Spec.Operations != nil && t.Spec.Operations.GetPendingTasks != nil {
		plugins = t.Spec.Operations.GetPendingTasks.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("tq-add-consumer-group-%s", t.Name),
		Path:    fmt.Sprintf("/v1/tasks/%s/add-consumer-group", t.Name),
		Plugins: plugins,
		OpenAPI: &apis.OpenAPI{PathDef: &openapi3.PathItem{Post: operation}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prepare context object
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
			defer cancel()

			// Prepare a temporary logger
			logger := m.logger.With(
				zap.String("source", t.Spec.Source),
				zap.String("taskqueue", t.Name),
				zap.String("operation", "get-pending-tasks"),
			)

			req := new(AddConsumerGroupRequest)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				logger.Error("Unable to parse request body", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusBadRequest, err)
				return
			}

			// Get the source for the task
			s := m.taskQueueSources[t.Spec.Source]

			// Perform the "AddConsumerGroup" operation
			if err := s.AddConsumerGroup(ctx, t.Name, req.ConsumerGroup); err != nil {
				logger.Error("Unable to create consumer group", zap.Error(err))
				utils.SendErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			// We got no response for this one
			w.WriteHeader(http.StatusNoContent)
		}),
	}
}
