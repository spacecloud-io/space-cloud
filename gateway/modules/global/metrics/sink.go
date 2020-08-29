package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"
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
				helpers.Logger.LogDebug(helpers.GetRequestID(nil), "Unable to get internal access token", map[string]interface{}{"error": err.Error()})
				continue
			}
			result := struct {
				Error  string        `json:"error"`
				Result []interface{} `json:"result"`
			}{}
			url := fmt.Sprintf("http://%s/v1/runner/metrics", m.syncMan.GetRunnerAddr())
			if err := m.syncMan.MakeHTTPRequest(context.Background(), http.MethodGet, url, token, "", map[string]interface{}{}, &result); err != nil {
				helpers.Logger.LogDebug(helpers.GetRequestID(nil), "Unable to fetch metrics from runner", map[string]interface{}{"error": err.Error()})
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
		_ = helpers.Logger.LogError(helpers.GetRequestID(nil), "Unable to push metrics", err, nil)
		return
	}
	if result.Status != http.StatusOK {
		_ = helpers.Logger.LogError(helpers.GetRequestID(nil), "Unable to push metrics", err, map[string]interface{}{"statusCode": result.Status})
	}
}
