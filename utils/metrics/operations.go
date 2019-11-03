package metrics

import (
	"sync/atomic"
	"time"

	"github.com/spaceuptech/space-cloud/utils"
)

// AddEgress add the bytes to the egress counter of that project
func (m *Module) AddEgress(project string, bytes uint64) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module or bandwidth measurement is disabled
	if !m.config.IsEnabled || m.config.DisableBandwidth {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(project, newMetrics())
	metrics := metricsTemp.(*metrics)

	atomic.AddUint64(&metrics.bw.egressBW, bytes)
}

// AddIngress add the bytes to the ingress counter of that project
func (m *Module) AddIngress(project string, bytes uint64) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module or bandwidth measurement is disabled
	if !m.config.IsEnabled || m.config.DisableBandwidth {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(project, newMetrics())
	metrics := metricsTemp.(*metrics)

	atomic.AddUint64(&metrics.bw.ingressBW, bytes)
}

func (m *Module) AddDBOperation(project, dbType, col string, count int64, op utils.OperationType) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module is disabled
	if !m.config.IsEnabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(project, newMetrics())
	metrics := metricsTemp.(*metrics)

	crudTemp, _ := metrics.crud.LoadOrStore(dbType+":"+col, new(crudMetrics))
	crud := crudTemp.(*crudMetrics)

	switch op {
	case utils.Create:
		atomic.AddUint64(&crud.create, uint64(count))

	case utils.Read:
		atomic.AddUint64(&crud.read, uint64(count))

	case utils.Update:
		atomic.AddUint64(&crud.update, uint64(count))

	case utils.Delete:
		atomic.AddUint64(&crud.delete, uint64(count))

	case utils.Batch:
		atomic.AddUint64(&crud.batch, uint64(count))
	}
}

func (m *Module) LoadMetrics() []interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Create an array of metric docs
	metricDocs := make([]interface{}, 0)

	// Capture the current time
	t := time.Now()

	// Iterate over all projects to generate the metric docs
	m.projects.Range(func(key, value interface{}) bool {

		// Load the project and metrics object
		project := key.(string)
		metrics := value.(*metrics)

		metricDocs = append(metricDocs, m.createBWDocuments(project, &metrics.bw, &t)...)
		metricDocs = append(metricDocs, m.createCrudDocuments(project, &metrics.crud, &t)...)

		// Delete the project
		m.projects.Delete(key)

		return true
	})

	return metricDocs
}
