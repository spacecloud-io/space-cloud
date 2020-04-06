package metrics

import (
	"strings"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func newMetrics() *metrics {
	return &metrics{}
}

func (m *Module) createFileDocuments(key string, metrics *metricOperations, t string) []interface{} {
	docs := make([]interface{}, 0)

	arr := strings.Split(key, ":")
	module := "file"
	if metrics.create > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], "na", module, utils.Create, metrics.create, t))
	}

	if metrics.read > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], "na", module, utils.Read, metrics.read, t))
	}

	if metrics.delete > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], "na", module, utils.Delete, metrics.delete, t))
	}

	if metrics.list > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], "na", module, utils.List, metrics.list, t))
	}

	return docs
}

func (m *Module) createCrudDocuments(key string, value *metricOperations, t string) []interface{} {
	docs := make([]interface{}, 0)

	arr := strings.Split(key, ":")
	module := "db"
	if value.create > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], arr[2], module, utils.Create, value.create, t))
	}

	if value.read > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], arr[2], module, utils.Read, value.read, t))
	}

	if value.update > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], arr[2], module, utils.Update, value.update, t))
	}

	if value.delete > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], arr[2], module, utils.Delete, value.delete, t))
	}

	if value.batch > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], arr[2], module, utils.Batch, value.batch, t))
	}

	return docs
}

func (m *Module) createEventDocument(key string, count uint64, t string) []interface{} {
	arr := strings.Split(key, ":")
	docs := make([]interface{}, 0)
	if count > 0 {
		docs = append(docs, m.createDocument(arr[0], "na", "na", "eventing", utils.OperationType(arr[1]), count, t))
	}
	return docs
}

func (m *Module) createFunctionDocument(key string, count uint64, t string) []interface{} {
	arr := strings.Split(key, ":")
	docs := make([]interface{}, 0)
	if count > 0 {
		docs = append(docs, m.createDocument(arr[0], arr[1], arr[2], "function", "calls", count, t))
	}
	return docs
}

func (m *Module) createDocument(project, driver, subType, module string, op utils.OperationType, count uint64, t string) interface{} {
	return map[string]interface{}{
		"id":         ksuid.New().String(),
		"project_id": project,
		"module":     module,
		"type":       op,
		"sub_type":   subType,
		"ts":         t,
		"count":      count,
		"driver":     driver,
		"node_id":    m.nodeID,
		"cluster_id": m.clusterID,
	}
}
