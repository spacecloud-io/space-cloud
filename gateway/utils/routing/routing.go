package routing

import "sync"

// Routing manages the routing functionality of space cloud
type Routing struct {
	lock sync.RWMutex

	routes routeMapping
}

// New creates a new instance of the routing module
func New() *Routing {
	return &Routing{routes: routeMapping{}}
}
