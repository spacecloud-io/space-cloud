package metrics

import (
	"github.com/sirupsen/logrus"
	api "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/db"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
	"sync"
)

// Module struct for metrics
type Module struct {
	lock    sync.RWMutex
	syncMan *syncman.Manager
	isProd  bool
	ssl     *config.SSL

	// Variables for metric state
	clusterID string
	nodeID    string
	projects  sync.Map // key -> project; value -> *metrics
	eventing  map[string]int
	// Variables to store the configuration
	config Config

	// Variables to interact with the sink
	sink *db.DB
}

// Config is the configuration required by the metrics module
type Config struct {
	IsDisabled       bool
	DisableBandwidth bool
	SinkType         string
	SinkConn         string
	Scope            string
}

// New creates a new instance of the metrics module
func New(clusterID, nodeID string, config *Config, syncMan *syncman.Manager, isProd bool) (*Module, error) {

	// Return an empty object if the module isn't enabled
	if config.IsDisabled {
		return new(Module), nil
	}

	// Initialise the sink
	conn := api.New("projectmetrics", "localhost:4122", false).DB("postgres")

	// Create a new metrics module
	m := &Module{nodeID: nodeID, clusterID: clusterID, sink: conn, config: *config, syncMan: syncMan, isProd: isProd}
	logrus.Println("metrics started")
	// Start routine to flush metrics to the sink
	go m.routineFlushMetricsToSink()

	return m, nil
}

// SetSSL sets ssl field
func (m *Module) SetSSL(ssl *config.SSL) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.ssl = ssl
}

type metrics struct {
	crud      sync.Map // key -> dbType:col; value -> *metricOperations
	fileStore sync.Map // key -> storeType value -> *metricOperations
	function  uint64
}

type metricOperations struct {
	create uint64
	read   uint64
	update uint64
	delete uint64
	batch  uint64
}
