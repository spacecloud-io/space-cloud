package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

const metricsUpdaterInterval = 5 * time.Minute

// NOTE: test not written for below function
func (m *Module) routineFlushMetricsToSink() {
	ticker := time.NewTicker(metricsUpdaterInterval)

	for range ticker.C {
		if m.isMetricDisabled {
			continue
		}
		go m.flushMetrics(m.LoadMetrics())

		if m.syncMan.GetRunnerAddr() != "" {
			token, err := m.adminMan.GetInternalAccessToken()
			if err != nil {
				utils.LogDebug("Unable to get internal access token", "metrics", "routine-flush-metrics-to-sink", map[string]interface{}{"error": err})
				continue
			}
			result := struct {
				Error  string        `json:"error"`
				Result []interface{} `json:"result"`
			}{}
			url := fmt.Sprintf("http://%s/v1/runner/metrics", m.syncMan.GetRunnerAddr())
			if err := m.syncMan.MakeHTTPRequest(context.Background(), http.MethodGet, url, token, "", map[string]interface{}{}, &result); err != nil {
				utils.LogDebug("Unable to fetch metrics from runner", "metrics", "routine-flush-metrics-to-sink", map[string]interface{}{"error": err})
				continue
			}
			go m.flushMetrics(result.Result)
		}

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
