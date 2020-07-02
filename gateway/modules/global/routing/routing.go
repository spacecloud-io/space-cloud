package routing

import (
	"sync"
	"text/template"
)

// Routing manages the routing functionality of space cloud
type Routing struct {
	lock sync.RWMutex

	routes      routeMapping
	goTemplates map[string]*template.Template
}

// New creates a new instance of the routing module
func New() *Routing {
	return &Routing{routes: routeMapping{}, goTemplates: map[string]*template.Template{}}
}
