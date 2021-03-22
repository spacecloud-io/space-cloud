package scaler

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/common/model"
)

func (s *Scaler) queryPrometheus(ctx context.Context, project, service, version, scalingMode string) (int64, error) {
	// Extract the prometheus metric name
	var metricName string
	switch scalingMode {
	case "requests-per-second":
		metricName = "sc:total_requests:rate"
	case "active-requests":
		metricName = "sc:active_requests:avg"
	default:
		return 0, fmt.Errorf("invalid scalingMode (%s) provided", scalingMode)
	}

	query := preparePrometheusQuery(project, service, version, metricName, "30s")
	result, _, err := s.prometheusClient.Query(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}
	vector := result.(model.Vector)
	if len(vector) == 0 {
		return 0, nil
	}

	return int64(vector[0].Value), nil
}

func preparePrometheusQuery(project, service, version, metric, duration string) string {
	return fmt.Sprintf("ceil(%s%s{kubernetes_namespace=\"%s\", app=\"%s\", version=\"%s\"})", metric, duration, project, service, version)
}
