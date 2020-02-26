package config

import (
	"fmt"
	"math/rand"
)

// Routes describes the configuration for the routing module
type Routes []*Route

// Route describes the parameters of a single route
type Route struct {
	ID      string        `json:"id" yaml:"id"`
	Source  RouteSource   `json:"source" yaml:"source"`
	Targets []RouteTarget `json:"targets" yaml:"targets"`
}

// SelectTarget returns a target based on the weights assigned
func (r *Route) SelectTarget(weight float64) (RouteTarget, error) {

	// Generate a random float in the range 0 to 1 if provided weight in lesser than zero
	if weight < 0 {
		weight = rand.Float64()
	}

	var cumulativeWeight float64

	// Return the first target which matches the range
	for _, target := range r.Targets {
		cumulativeWeight += target.Weight
		if weight <= cumulativeWeight {
			return target, nil
		}
	}

	// Return error if no targets match
	return RouteTarget{}, fmt.Errorf("no target found for route (%s) - make sure you have defined atleast one target with proper weights", r.Source.URL)
}

// RouteSource is the source of routing
type RouteSource struct {
	Hosts      []string     `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Methods    []string     `json:"methods,omitempty" yaml:"methods,omitempty"`
	URL        string       `json:"url,omitempty" yaml:"url,omitempty"`
	RewriteURL string       `json:"rewrite,omitempty" yaml:"rewrite,omitempty"`
	Type       RouteURLType `json:"type,omitempty" yaml:"type,omitempty"`
}

// RouteTarget is the destination of routing
type RouteTarget struct {
	Host   string  `json:"host,omitempty" yaml:"host,omitempty"`
	Port   string  `json:"port,omitempty" yaml:"port,omitempty"`
	Scheme string  `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Weight float64 `json:"weight,omitempty" yaml:"weight,omitempty"`
}

// RouteURLType describes how the url should be evaluated / matched
type RouteURLType string

const (
	// RoutePrefix is used for prefix matching
	RoutePrefix RouteURLType = "prefix"

	// RouteExact is used for matching the url exactly as it is
	RouteExact RouteURLType = "exact"
)
