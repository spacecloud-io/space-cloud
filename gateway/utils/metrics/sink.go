package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const metricsUpdaterInterval = 5 * time.Minute

// NOTE: test not written for below function
func (m *Module) routineFlushMetricsToSink() {
	ticker := time.NewTicker(metricsUpdaterInterval)

	for range ticker.C {
		go m.flushMetrics(m.LoadMetrics())

		go func() {
			// Flush project metrics only if our index is 0
			if index := m.syncMan.GetGatewayIndex(); index == 0 {
				c := m.syncMan.GetGlobalConfig()
				ssl := c.SSL
				for _, project := range c.Projects {
					m.updateSCMetrics(m.generateMetricsRequest(project, ssl))
				}
			}
		}()
	}
}

// NOTE: test not written for below function
func (m *Module) flushMetrics(docs []interface{}) {
	if len(docs) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := m.sink.Insert("operation_metrics").Docs(docs).Apply(ctx)
	if err != nil {
		logrus.Debugln("Unable to push metrics:", err)
		return
	}
	if result.Status != http.StatusOK {
		logrus.Debugln("Unable to push metrics:", result.Error)
	}
}
