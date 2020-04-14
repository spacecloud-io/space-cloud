package metrics

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	api "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/db"

	"github.com/spaceuptech/space-cloud/runner/model"
)

const notApplicable = "na"

// Module holds config of metrics
type Module struct {
	lock              sync.RWMutex
	isMetricsDisabled bool

	// Variables for metric state
	clusterID  string
	nodeID     string
	driverType string
	projects   sync.Map // key -> project; value -> *metrics

	// Variables to interact with the sink
	sink *db.DB
}

type metrics struct {
	serviceCall uint64
}

// New creates a instance of metrics package
func New(isMetricDisabled bool, driverType model.DriverType) *Module {
	m := &Module{
		isMetricsDisabled: isMetricDisabled,
		clusterID:         os.Getenv("CLUSTER_ID"),
		nodeID:            ksuid.New().String(),
		sink:              api.New("spacecloud", "api.spaceuptech.com", true).DB("db"),
		driverType:        string(driverType),
	}
	go m.routineMetrics()

	return m
}

const metricsUpdaterInterval = 5 * time.Minute

// RoutineMetrics routinely sends anonymous metrics
func (m *Module) routineMetrics() {
	ticker := time.NewTicker(metricsUpdaterInterval)
	defer ticker.Stop()

	for range ticker.C {
		_ = m.updateRunnerMetrics(m.loadMetrics())
	}
}

func (m *Module) updateRunnerMetrics(doc []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := m.sink.Insert("operation_metrics").Docs(doc).Apply(ctx)
	if err != nil {
		logrus.Errorf("error querying database got error")
		return err
	}
	if result.Status != http.StatusOK {
		logrus.Errorf("error querying database got status (%d) (%s)", result.Status, result.Error)
		return fmt.Errorf("unknown status code (%d) (%s)", result.Status, result.Error)
	}

	return nil
}

func newMetrics() *metrics {
	return &metrics{}
}

func (m *Module) createDocument(project, driver, subType, module string, op string, count uint64, t string) map[string]interface{} {
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

// AddServiceCall counts how many times service apply gets called
func (m *Module) AddServiceCall(projectID string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the metrics module is disabled
	if m.isMetricsDisabled {
		return
	}

	metricsTemp, _ := m.projects.LoadOrStore(projectID, newMetrics())
	metrics := metricsTemp.(*metrics)

	atomic.AddUint64(&metrics.serviceCall, uint64(1))
}

func (m *Module) loadMetrics() []interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()

	docs := make([]interface{}, 0)

	m.projects.Range(func(key, value interface{}) bool {
		metrics := value.(*metrics)
		if metrics.serviceCall > 0 {
			docs = append(docs, m.createDocument(key.(string), m.driverType, notApplicable, "service", "apply", metrics.serviceCall, time.Now().String()))
		}
		// Delete the key from the map
		m.projects.Delete(key)
		return true
	})
	return docs
}
