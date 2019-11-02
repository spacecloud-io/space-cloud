package metrics

import (
	"sync"

	"github.com/spaceuptech/space-cloud/modules/crud"
)

type Module struct {
	lock sync.RWMutex

	// Variables for metric state
	nodeID   string
	projects sync.Map // key -> project; value -> *metrics

	// Variables to store the configuration
	config Config

	// Variables to interact with the sink
	sink *crud.Module
}

// The configuration required by the metrics module
type Config struct {
	IsEnabled        bool
	DisableBandwidth bool
	SinkType         string
	SinkConn         string
	Scope            string
}

func New(nodeID string, config *Config) (*Module, error) {

	// Return an empty object if the module isn't enabled
	if !config.IsEnabled {
		return new(Module), nil
	}

	// Initialise the sink
	c, err := initialiseSink(config)
	if err != nil {
		return nil, err
	}

	// Create a new metrics module
	m := &Module{nodeID: nodeID, sink: c, config: *config}

	// Start routine to flush metrics to the sink
	go m.routineFlushMetricsToSink()

	return m, nil
}

type metrics struct {
	bw   bwMetrics
	crud sync.Map // key -> dbType:col; value -> *dbMetrics
}

type bwMetrics struct {
	ingressBW uint64
	egressBW  uint64
}

type crudMetrics struct {
	create uint64
	read   uint64
	update uint64
	delete uint64
	batch  uint64
}
