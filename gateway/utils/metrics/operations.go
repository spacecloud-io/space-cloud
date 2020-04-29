package metrics

import (
	"sync/atomic"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// AddEventingType counts the number of time a particular event type is called
func (m *Module) AddEventingType(project, eventingType string) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Return if the metrics module is disabled
	if m.isMetricDisabled {
		return
	}
	value, _ := m.projects.LoadOrStore(generateEventingKey(project, eventingType), newMetrics())
	metrics := value.(*metrics)
	atomic.AddUint64(&metrics.eventing, uint64(1))
}

// AddFunctionOperation counts the number of time a particular function gets invoked
func (m *Module) AddFunctionOperation(project, service, function string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module is disabled
	if m.isMetricDisabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(generateFunctionKey(project, service, function), newMetrics())
	metrics := metricsTemp.(*metrics)
	atomic.AddUint64(&metrics.function, uint64(1))
}

// AddDBOperation adds a operation to the database
func (m *Module) AddDBOperation(project, dbType, col string, count int64, op utils.OperationType) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Return if the metrics module is disabled
	if m.isMetricDisabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(generateDatabaseKey(project, dbType, col), newMetrics())
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
	}
}

// AddFileOperation adds a operation to the database
func (m *Module) AddFileOperation(project, storeType string, op utils.OperationType) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module is disabled
	if m.isMetricDisabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(generateFileKey(project, storeType), newMetrics())
	metrics := metricsTemp.(*metrics)

	switch op {
	case utils.Create:
		atomic.AddUint64(&metrics.fileStore.create, uint64(1))

	case utils.Read:
		atomic.AddUint64(&metrics.fileStore.read, uint64(1))

	case utils.Delete:
		atomic.AddUint64(&metrics.fileStore.delete, uint64(1))

	case utils.List:
		atomic.AddUint64(&metrics.fileStore.list, uint64(1))
	}
}

// LoadMetrics loads the metrics
// NOTE: test not written for below function
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
		metrics := value.(*metrics)
		switch getModuleName(key.(string)) {
		case eventingModule:
			metricDocs = append(metricDocs, m.createEventDocument(key.(string), metrics.eventing, t)...)
		case fileModule:
			metricDocs = append(metricDocs, m.createFileDocuments(key.(string), &metrics.fileStore, t)...)
		case databaseModule:
			metricDocs = append(metricDocs, m.createCrudDocuments(key.(string), &metrics.crud, t)...)
		case remoteServiceModule:
			metricDocs = append(metricDocs, m.createFunctionDocument(key.(string), metrics.function, t)...)
		}
		// Delete the project
		m.projects.Delete(key)

		return true
	})

	return metricDocs
}
