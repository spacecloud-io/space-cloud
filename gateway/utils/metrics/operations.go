package metrics

import (
	"fmt"
	"strings"
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
	value, _ := m.projects.LoadOrStore(fmt.Sprintf("%s:%s:%s", "eventing", project, eventingType), newMetrics())
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

	metricsTemp, _ := m.projects.LoadOrStore(fmt.Sprintf("%s:%s:%s:%s", "function", project, service, function), newMetrics())
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

	metricsTemp, _ := m.projects.LoadOrStore(fmt.Sprintf("%s:%s:%s:%s", "db", project, dbType, col), newMetrics())
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

	metricsTemp, _ := m.projects.LoadOrStore(fmt.Sprintf("%s:%s:%s", "file", project, storeType), newMetrics())
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
		v := strings.Split(key.(string), ":")
		metrics := value.(*metrics)
		switch v[0] {
		case "eventing":
			metricDocs = append(metricDocs, m.createEventDocument(strings.Join(v[1:], ":"), metrics.eventing, t)...)
		case "file":
			metricDocs = append(metricDocs, m.createFileDocuments(strings.Join(v[1:], ":"), &metrics.fileStore, t)...)
		case "db":
			metricDocs = append(metricDocs, m.createCrudDocuments(strings.Join(v[1:], ":"), &metrics.crud, t)...)
		case "function":
			metricDocs = append(metricDocs, m.createFunctionDocument(strings.Join(v[1:], ":"), metrics.function, t)...)
		}
		// Delete the project
		m.projects.Delete(key)

		return true
	})

	return metricDocs
}
