package config

import (
	"fmt"
	"math/rand"
	"strings"
)

// Routes describes the configuration for the routing module
type Routes []*Route

// Len return length of routes
func (a Routes) Len() int { return len(a) }

// Swap swaps two element of the array
func (a Routes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less compares two elements of the array
func (a Routes) Less(i, j int) bool {
	arrayI := strings.Split(a[i].Source.URL, "/")
	arrayJ := strings.Split(a[j].Source.URL, "/")

	lenI := len(arrayI)
	lenJ := len(arrayJ)

	if arrayI[lenI-1] == "" {
		lenI--
	}
	if arrayJ[lenJ-1] == "" {
		lenJ--
	}

	return lenI > lenJ
}

// Route describes the parameters of a single route
type Route struct {
	ID      string        `json:"id" yaml:"id"`
	Project string        `json:"project" yaml:"project"`
	Source  RouteSource   `json:"source" yaml:"source"`
	Targets []RouteTarget `json:"targets" yaml:"targets"`
	Rule    *Rule         `json:"rule" yaml:"rule"`
	Modify  struct {
		Tmpl     EndpointTemplatingEngine `json:"template,omitempty" yaml:"template,omitempty"`
		ReqTmpl  string                   `json:"requestTemplate" yaml:"requestTemplate"`
		ResTmpl  string                   `json:"responseTemplate" yaml:"responseTemplate"`
		OpFormat string                   `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty"`
		Headers  []struct {
			Key   string `json:"key" yaml:"key"`
			Value string `json:"value" yaml:"value"`
		} `json:"headers" yaml:"headers"`
	} `json:"modify" yaml:"modify"`
}

// SelectTarget returns a target based on the weights assigned
func (r *Route) SelectTarget(weight int32) (RouteTarget, error) {

	// Generate a random float in the range 0 to 100 if provided weight in lesser than zero
	if weight < 0 {
		weight = rand.Int31n(100)
	}

	var cumulativeWeight int32

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
	Hosts      []string     `json:"hosts" yaml:"hosts"`
	Methods    []string     `json:"methods" yaml:"methods"`
	URL        string       `json:"url" yaml:"url"`
	RewriteURL string       `json:"rewrite" yaml:"rewrite"`
	Type       RouteURLType `json:"type" yaml:"type"`
	Port       int32        `json:"port" yaml:"port"`
}

// RouteTarget is the destination of routing
type RouteTarget struct {
	Host    string          `json:"host" yaml:"host"`
	Port    int32           `json:"port" yaml:"port"`
	Scheme  string          `json:"scheme" yaml:"scheme"`
	Weight  int32           `json:"weight" yaml:"weight"`
	Version string          `json:"version" yaml:"version"`
	Type    RouteTargetType `json:"type" yaml:"type"`
}

// RouteURLType describes how the url should be evaluated / matched
type RouteURLType string

const (
	// RoutePrefix is used for prefix matching
	RoutePrefix RouteURLType = "prefix"

	// RouteExact is used for matching the url exactly as it is
	RouteExact RouteURLType = "exact"
)

// RouteTargetType describes how the target should be selected
type RouteTargetType string

const (
	// RouteTargetVersion is used to route to versions of the same service
	RouteTargetVersion RouteTargetType = "version"

	// RouteTargetExternal is used to route to external services
	RouteTargetExternal RouteTargetType = "external"
)
