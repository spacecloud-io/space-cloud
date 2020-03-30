package metrics

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (m *Module) routineFlushMetricsToSink() {
	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
		logrus.Println("executing ticker")
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
	logrus.Println("docs of flush", docs)
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
