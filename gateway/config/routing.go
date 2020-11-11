package config

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/spaceuptech/helpers"
)

// GlobalRoutesConfig describes the project level config for ingress routing
type GlobalRoutesConfig struct {
	RequestHeaders  Headers `json:"headers" yaml:"headers" mapstructure:"headers"`
	ResponseHeaders Headers `json:"resHeaders" yaml:"resHeaders" mapstructure:"resHeaders"`
}

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
	ID      string        `json:"id" yaml:"id" mapstructure:"id"`
	Project string        `json:"project" yaml:"project" mapstructure:"project"`
	Source  RouteSource   `json:"source" yaml:"source" mapstructure:"source"`
	Targets []RouteTarget `json:"targets" yaml:"targets" mapstructure:"targets"`
	Rule    *Rule         `json:"rule" yaml:"rule" mapstructure:"rule"`
	Modify  struct {
		Tmpl            TemplatingEngine `json:"template,omitempty" yaml:"template,omitempty" mapstructure:"template"`
		ReqTmpl         string           `json:"requestTemplate" yaml:"requestTemplate" mapstructure:"requestTemplate"`
		ResTmpl         string           `json:"responseTemplate" yaml:"responseTemplate" mapstructure:"responseTemplate"`
		OpFormat        string           `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" mapstructure:"outputFormat"`
		RequestHeaders  Headers          `json:"headers" yaml:"headers" mapstructure:"headers"`
		ResponseHeaders Headers          `json:"resHeaders" yaml:"resHeaders" mapstructure:"resHeaders"`
	} `json:"modify" yaml:"modify" mapstructure:"modify"`
}

// SelectTarget returns a target based on the weights assigned
func (r *Route) SelectTarget(ctx context.Context, weight int32) (RouteTarget, error) {

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
	return RouteTarget{}, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("No target found for route (%s) - make sure you have defined atleast one target with proper weights", r.Source.URL), nil, nil)
}

// RouteSource is the source of routing
type RouteSource struct {
	Hosts      []string     `json:"hosts" yaml:"hosts" mapstructure:"hosts"`
	Methods    []string     `json:"methods" yaml:"methods" mapstructure:"methods"`
	URL        string       `json:"url" yaml:"url" mapstructure:"url"`
	RewriteURL string       `json:"rewrite" yaml:"rewrite" mapstructure:"rewrite"`
	Type       RouteURLType `json:"type" yaml:"type" mapstructure:"type"`
	Port       int32        `json:"port" yaml:"port" mapstructure:"port"`
}

// RouteTarget is the destination of routing
type RouteTarget struct {
	Host    string          `json:"host" yaml:"host" mapstructure:"host"`
	Port    int32           `json:"port" yaml:"port" mapstructure:"port"`
	Scheme  string          `json:"scheme" yaml:"scheme" mapstructure:"scheme"`
	Weight  int32           `json:"weight" yaml:"weight" mapstructure:"weight"`
	Version string          `json:"version" yaml:"version" mapstructure:"version"`
	Type    RouteTargetType `json:"type" yaml:"type" mapstructure:"type"`
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
