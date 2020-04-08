package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const metricsUpdaterInterval = 30 * time.Second

func (m *Module) routineFlushMetricsToSink() {
	ticker := time.NewTicker(metricsUpdaterInterval)

	for range ticker.C {
		go m.flushMetrics(m.LoadMetrics())

		find, set, min, isSkip := m.generateMetricsRequest()
		if isSkip {
			continue
		}
		m.updateSCMetrics(find, set, min)
	}
}

func (m *Module) flushMetrics(docs []interface{}) {
	if len(docs) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	result, err := m.sink.Insert("operation_metrics").Docs(docs).Apply(ctx)
	if err != nil {
		logrus.Errorf("error querying database got error")
	}
	if result.Status != http.StatusOK {
		logrus.Errorf("error querying database got status (%d) (%s)", result.Status, result.Error)
	}
}
