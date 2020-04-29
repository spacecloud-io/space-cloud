package metrics

import (
	"sync"

	api "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/db"

	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// Module struct for metrics
type Module struct {
	lock   sync.RWMutex
	isProd bool

	// Variables for metric state
	clusterID string
	nodeID    string
	projects  sync.Map // key -> project; value -> *metrics
	// Variables to store the configuration
	isMetricDisabled bool
	// Variables to interact with the sink
	sink *db.DB

	// Global modules
	adminMan *admin.Manager
	syncMan  *syncman.Manager
}

// Config is the configuration required by the metrics module
type Config struct {
	IsDisabled bool
}

// New creates a new instance of the metrics module
func New(clusterID, nodeID string, isMetricDisabled bool, adminMan *admin.Manager, syncMan *syncman.Manager, isProd bool) (*Module, error) {

	// Return an empty object if the module isn't enabled
	if isMetricDisabled {
		return new(Module), nil
	}

	// Initialise the sink
	conn := api.New("spacecloud", "api.spaceuptech.com", true).DB("db")

	// Create a new metrics module
	m := &Module{nodeID: nodeID, clusterID: clusterID, sink: conn, isMetricDisabled: isMetricDisabled, adminMan: adminMan, syncMan: syncMan, isProd: isProd}
	// Start routine to flush metrics to the sink
	go m.routineFlushMetricsToSink()

	return m, nil
}

type metrics struct {
	crud      metricOperations // key -> dbType:col; value -> *metricOperations
	fileStore metricOperations // key -> storeType value -> *metricOperations
	eventing  uint64
	function  uint64
}

type metricOperations struct {
	create uint64
	read   uint64
	update uint64
	delete uint64
	list   uint64
}
