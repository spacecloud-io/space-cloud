package metrics

import (
	"strings"
	"sync"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func newMetrics() *metrics {
	return &metrics{}
}

func (m *Module) createFileDocuments(project string, opMetrics *sync.Map, t *time.Time) []interface{} {
	docs := make([]interface{}, 0)

	opMetrics.Range(func(key, value interface{}) bool {
		parts := key.(string)
		metrics := value.(*metricOperations)
		module := "file"
		if metrics.create > 0 {
			docs = append(docs, m.createDocument(project, parts, "na", module, utils.Create, metrics.create, t))
		}

		if metrics.read > 0 {
			docs = append(docs, m.createDocument(project, parts, "na", module, utils.Read, metrics.read, t))
		}

		if metrics.update > 0 {
			docs = append(docs, m.createDocument(project, parts, "na", module, utils.Update, metrics.update, t))
		}

		if metrics.delete > 0 {
			docs = append(docs, m.createDocument(project, parts, "na", module, utils.Delete, metrics.delete, t))
		}

		// Delete the key from the map
		opMetrics.Delete(key)

		return true
	})

	return docs
}

func (m *Module) createCrudDocuments(project string, opMetrics *sync.Map, t *time.Time) []interface{} {
	docs := make([]interface{}, 0)

	opMetrics.Range(func(key, value interface{}) bool {
		parts := strings.Split(key.(string), ":")
		metrics := value.(*metricOperations)
		module := "db"
		if metrics.create > 0 {
			docs = append(docs, m.createDocument(project, parts[0], parts[1], module, utils.Create, metrics.create, t))
		}

		if metrics.read > 0 {
			docs = append(docs, m.createDocument(project, parts[0], parts[1], module, utils.Read, metrics.read, t))
		}

		if metrics.update > 0 {
			docs = append(docs, m.createDocument(project, parts[0], parts[1], module, utils.Update, metrics.update, t))
		}

		if metrics.delete > 0 {
			docs = append(docs, m.createDocument(project, parts[0], parts[1], module, utils.Delete, metrics.delete, t))
		}

		if metrics.batch > 0 {
			docs = append(docs, m.createDocument(project, parts[0], parts[1], module, utils.Batch, metrics.batch, t))
		}

		// Delete the key from the map
		opMetrics.Delete(key)

		return true
	})

	return docs
}

func (m *Module) createDocument(project, dbType, col, module string, op utils.OperationType, count uint64, t *time.Time) interface{} {
	return map[string]interface{}{
		"id":         ksuid.New().String(),
		"project_id": project,
		"module":     module,
		"type":       op,
		"sub_type":   col,
		"ts":         t.String(),
		"count":      count,
		"driver":     dbType,
		"node_id":    m.nodeID,
		"cluster_id": m.clusterID,
	}
}
