package global

import (
	"github.com/spaceuptech/space-cloud/gateway/managers"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/caching"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/metrics"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

// Global holds global modules
type Global struct {
	letsencrypt *letsencrypt.LetsEncrypt
	metrics     *metrics.Module
	routing     *routing.Routing
	caching     *caching.Cache
}

// New creates a new global object
func New(clusterID, nodeID string, isDev bool, managers *managers.Managers) (*Global, error) {
	// when ever gateway starts it will send metrics
	m, err := metrics.New(clusterID, nodeID, false, managers.Admin(), managers.Sync(), !isDev)
	if err != nil {
		return nil, err
	}

	// Initialise a lets encrypt client
	le, err := letsencrypt.New()
	if err != nil {
		return nil, err
	}

	// Initialise the routing module
	r := routing.New()

	// Initialise the caching module
	c := caching.Init(clusterID, nodeID)
	c.SetAdminModule(managers.Admin())
	r.SetCachingModule(c)

	return &Global{letsencrypt: le, metrics: m, routing: r, caching: c}, nil
}

// LetsEncrypt returns the letsencrypt module
func (g *Global) LetsEncrypt() *letsencrypt.LetsEncrypt {
	return g.letsencrypt
}

// Metrics returns the metrics module
func (g *Global) Metrics() *metrics.Module {
	return g.metrics
}

// Routing returns the routing module
func (g *Global) Routing() *routing.Routing {
	return g.routing
}

// Caching returns the caching module
func (g *Global) Caching() *caching.Cache {
	return g.caching
}
