package metrics

import (
	"context"
	"log"
	"time"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) routineFlushMetricsToSink() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		go m.flushMetrics(m.LoadMetrics())
	}
}

// Right now we return a crud block since we only suppport databases as a sink.
// In the future this would return an interface to abstract the sinks
func initialiseSink(c *Config) (*crud.Module, error) {

	// Create a new crud module
	sink := crud.Init(true)

	// Configure the crud module
	if err := sink.SetConfig(c.Scope, config.Crud{c.SinkType: &config.CrudStub{Enabled: true, Conn: c.SinkConn}}); err != nil {
		return nil, err
	}

	return sink, nil
}

func (m *Module) flushMetrics(docs []interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := m.sink.Create(ctx, m.config.SinkType, m.config.Scope, "metrics",
		&model.CreateRequest{Document: docs, Operation: utils.All}); err != nil {
		log.Println("Metrics module: Couldn't flush metrics to disk -", err)
	}
}
