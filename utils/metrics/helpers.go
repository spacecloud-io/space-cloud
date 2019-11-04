package metrics

import (
	"strings"
	"sync"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/utils"
)

func newMetrics() *metrics {
	return &metrics{bw: bwMetrics{}}
}

func (m *Module) createCrudDocuments(project string, dbMetrics *sync.Map, t *time.Time) []interface{} {
	docs := make([]interface{}, 0)

	dbMetrics.Range(func(key, value interface{}) bool {
		parts := strings.Split(key.(string), ":")
		metrics := value.(*crudMetrics)

		if metrics.create > 0 {
			docs = append(docs, m.createCrudDocument(project, parts[0], parts[1], utils.Create, metrics.create, t))
		}

		if metrics.read > 0 {
			docs = append(docs, m.createCrudDocument(project, parts[0], parts[1], utils.Read, metrics.read, t))
		}

		if metrics.update > 0 {
			docs = append(docs, m.createCrudDocument(project, parts[0], parts[1], utils.Update, metrics.update, t))
		}

		if metrics.delete > 0 {
			docs = append(docs, m.createCrudDocument(project, parts[0], parts[1], utils.Delete, metrics.delete, t))
		}

		if metrics.batch > 0 {
			docs = append(docs, m.createCrudDocument(project, parts[0], parts[1], utils.Batch, metrics.batch, t))
		}

		// Delete the key from the map
		dbMetrics.Delete(key)

		return true
	})

	return docs
}

func (m *Module) createCrudDocument(project, dbType, col string, op utils.OperationType, count uint64, t *time.Time) interface{} {
	return map[string]interface{}{
		"id":         ksuid.New().String(),
		"project_id": project,
		"module":     "db",
		"type":       op,
		"sub_type":   col,
		"ts":         *t,
		"count":      count,
		"driver":     dbType,
		"node_id":    "sc-" + m.nodeID,
	}
}

func (m *Module) createFileDocument(project string, storeType utils.FileStoreType, op utils.FileOpType, count uint64, t time.Time) interface{} {
	return map[string]interface{}{
		"id":         ksuid.New().String(),
		"project_id": project,
		"module":     "file",
		"type":       op,
		"sub_type":   "na",
		"ts":         t,
		"count":      count,
		"driver":     storeType,
		"node_id":    "sc-" + m.nodeID,
	}
}

func (m *Module) createBWDocuments(project string, bw *bwMetrics, t *time.Time) []interface{} {
	docs := make([]interface{}, 0)

	if bw.egressBW > 0 {
		docs = append(docs, m.createBWDocument(project, "egress", bw.egressBW, t))
	}

	if bw.ingressBW > 0 {
		docs = append(docs, m.createBWDocument(project, "ingress", bw.ingressBW, t))
	}

	return docs
}

func (m *Module) createBWDocument(project, op string, count uint64, t *time.Time) interface{} {
	return map[string]interface{}{
		"id":         ksuid.New().String(),
		"project_id": project,
		"module":     "bw",
		"type":       op,
		"sub_type":   "na",
		"ts":         *t,
		"count":      count,
		"driver":     "na",
		"node_id":    "sc-" + m.nodeID,
	}
}
