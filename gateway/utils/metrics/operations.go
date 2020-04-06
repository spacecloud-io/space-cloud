package metrics

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// AddEventingType counts the number of time a particular event type is called
func (m *Module) AddEventingType(project, eventingType string) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Return if the metrics module is disabled
	if m.config.IsDisabled {
		return
	}

	value, ok := m.eventing.LoadOrStore(fmt.Sprintf("%s:%s", project, eventingType), uint64(1))
	if ok {
		v := value.(uint64)
		atomic.AddUint64(&v, uint64(1))
		return
	}
}

// AddFunctionOperation counts the number of time a particular function gets invoked
func (m *Module) AddFunctionOperation(project, service, function string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module is disabled
	if m.config.IsDisabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(fmt.Sprintf("%s:%s:%s", project, service, function), newMetrics())
	metrics := metricsTemp.(*metrics)

	atomic.AddUint64(&metrics.function, uint64(1))
}

// AddDBOperation adds a operation to the database
func (m *Module) AddDBOperation(project, dbType, col string, count int64, op utils.OperationType) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module is disabled
	if m.config.IsDisabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(fmt.Sprintf("%s:%s:%s", project, dbType, col), newMetrics())
	metrics := metricsTemp.(*metrics)

	switch op {
	case utils.Create:
		atomic.AddUint64(&metrics.crud.create, uint64(count))

	case utils.Read:
		atomic.AddUint64(&metrics.crud.read, uint64(count))

	case utils.Update:
		atomic.AddUint64(&metrics.crud.update, uint64(count))

	case utils.Delete:
		atomic.AddUint64(&metrics.crud.delete, uint64(count))

	case utils.Batch:
		atomic.AddUint64(&metrics.crud.batch, uint64(count))
	}
}

// AddFileOperation adds a operation to the database
func (m *Module) AddFileOperation(project, storeType string, op utils.OperationType) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module is disabled
	if m.config.IsDisabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(fmt.Sprintf("%s:%s", project, storeType), newMetrics())
	metrics := metricsTemp.(*metricOperations)

	switch op {
	case utils.Create:
		atomic.AddUint64(&metrics.create, uint64(1))

	case utils.Read:
		atomic.AddUint64(&metrics.read, uint64(1))

	case utils.Update:
		atomic.AddUint64(&metrics.update, uint64(1))

	case utils.Delete:
		atomic.AddUint64(&metrics.delete, uint64(1))

	case utils.List:
		atomic.AddUint64(&metrics.list, uint64(1))
	}
}

// LoadMetrics loads the metrics
func (m *Module) LoadMetrics() []interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Create an array of metric docs)
	metricDocs := make([]interface{}, 0)

	// Capture the current time
	t := time.Now().Format(time.RFC3339)

	// Iterate over all projects to generate the metric docs
	m.projects.Range(func(key, value interface{}) bool {

		// Load the project and metrics object
		project := key.(string)
		metrics := value.(*metrics)

		metricDocs = append(metricDocs, m.createCrudDocuments(project, &metrics.crud, t)...)
		metricDocs = append(metricDocs, m.createFileDocuments(project, &metrics.fileStore, t)...)
		metricDocs = append(metricDocs, m.createFunctionDocument(project, metrics.function, t)...)
		metricDocs = append(metricDocs, m.createEventDocument(project, metrics.eventing, t)...)
		// Delete the project
		m.projects.Delete(key)

		return true
	})

	return metricDocs
}
