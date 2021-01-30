package model

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/spaceuptech/helpers"
)

const (
	// DefaultRequestRetries specifies the default values of service request retries
	DefaultRequestRetries int32 = 3
	// DefaultRequestTimeout specifies the default values of service request timeouts
	DefaultRequestTimeout int64 = 180 // Time in seconds
)

// Routes describes the configuration for the routing module
type Routes []*Route

// Route describes the parameters of a single route
type Route struct {
	ID             string        `json:"id" yaml:"id"`
	RequestRetries int32         `json:"requestRetries" yaml:"requestRetries"`
	RequestTimeout int64         `json:"requestTimeout" yaml:"requestTimeout"`
	Matchers       []*Matcher    `json:"matchers" yaml:"matchers"`
	Source         RouteSource   `json:"source" yaml:"source"`
	Targets        []RouteTarget `json:"targets" yaml:"targets"`
}

// SelectTarget returns a target based on the weights assigned
func (r *Route) SelectTarget(ctx context.Context, weight int32) (RouteTarget, error) {

	// Generate a random float in the range 0 to 1 if provided weight in lesser than zero
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

// Matcher store the rules, which are used for traffic splitting between service
type Matcher struct {
	URL     *HTTPMatcher   `json:"url,omitempty" yaml:"url,omitempty"`
	Headers []*HTTPMatcher `json:"headers,omitempty" yaml:"headers,omitempty"`
}

// HTTPMatcher is matcher type
type HTTPMatcher struct {
	Key          string             `json:"key,omitempty" yaml:"key,omitempty"`
	Value        string             `json:"value,omitempty" yaml:"value,omitempty"`
	Type         RouteHTTPMatchType `json:"type,omitempty" yaml:"type,omitempty"`
	IsIgnoreCase bool               `json:"isIgnoreCase,omitempty" yaml:"isIgnoreCase,omitempty"`
}

// RouteSource is the source of routing
type RouteSource struct {
	Protocol   Protocol     `json:"protocol" yaml:"protocol"`
	Hosts      []string     `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Methods    []string     `json:"methods,omitempty" yaml:"methods,omitempty"`
	URL        string       `json:"url,omitempty" yaml:"url,omitempty"`
	RewriteURL string       `json:"rewrite,omitempty" yaml:"rewrite,omitempty"`
	Type       RouteURLType `json:"type,omitempty" yaml:"type,omitempty"`
	Port       int32        `json:"port,omitempty" yaml:"port,omitempty"`
}

// RouteTarget is the destination of routing
type RouteTarget struct {
	Host    string          `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int32           `json:"port,omitempty" yaml:"port,omitempty"`
	Weight  int32           `json:"weight,omitempty" yaml:"weight,omitempty"`
	Version string          `json:"version,omitempty" yaml:"version,omitempty"`
	Type    RouteTargetType `json:"type,omitempty" yaml:"type,omitempty"`
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

// RouteHTTPMatchType defines http match type
type RouteHTTPMatchType string

const (
	// RouteHTTPMatchTypeExact is used for exact match
	RouteHTTPMatchTypeExact RouteHTTPMatchType = "exact"
	// RouteHTTPMatchTypeRegex is used for regex match
	RouteHTTPMatchTypeRegex RouteHTTPMatchType = "regex"
	// RouteHTTPMatchTypePrefix is used for prefix match
	RouteHTTPMatchTypePrefix RouteHTTPMatchType = "prefix"
	// RouteHTTPMatchTypeCheckPresence is used for only checking the presence of header in the http request
	RouteHTTPMatchTypeCheckPresence RouteHTTPMatchType = "check-presence"
)
