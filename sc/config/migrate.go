package config

// ServiceMigrate holds the config of service
type ServiceMigrate struct {
	ID        string                      `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`    // eg. http://localhost:8080
	URL       string                      `json:"url,omitempty" yaml:"url,omitempty" mapstructure:"url"` // eg. http://localhost:8080
	Endpoints map[string]*EndpointMigrate `json:"endpoints,omitempty" yaml:"endpoints,omitempty" mapstructure:"endpoints"`
}

// EndpointMigrate holds the config of a endpointMigrate
type EndpointMigrate struct {
	Kind EndpointKind     `json:"kind" yaml:"kind" mapstructure:"kind"`
	Tmpl TemplatingEngine `json:"template,omitempty" yaml:"template,omitempty" mapstructure:"template"`
	// ReqPayloadFormat specifies the payload format
	// depending upon the payload format, the graphQL request that
	// gets converted to http request will use that format as it's payload
	// currently supported formats are application/json,multipart/form-data
	ReqPayloadFormat string   `json:"requestPayloadFormat" yaml:"requestPayloadFormat" mapstructure:"requestPayloadFormat"`
	ReqTmpl          string   `json:"requestTemplate" yaml:"requestTemplate" mapstructure:"requestTemplate"`
	GraphTmpl        string   `json:"graphTemplate" yaml:"graphTemplate" mapstructure:"graphTemplate"`
	ResTmpl          string   `json:"responseTemplate" yaml:"responseTemplate" mapstructure:"responseTemplate"`
	OpFormat         string   `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" mapstructure:"outputFormat"`
	Token            string   `json:"token,omitempty" yaml:"token,omitempty" mapstructure:"token"`
	Claims           string   `json:"claims,omitempty" yaml:"claims,omitempty" mapstructure:"claims"`
	Method           string   `json:"method" yaml:"method" mapstructure:"method"`
	Path             string   `json:"path" yaml:"path" mapstructure:"path"`
	Rule             *Rule    `json:"rule,omitempty" yaml:"rule,omitempty" mapstructure:"rule"`
	Headers          Headers  `json:"headers,omitempty" yaml:"headers,omitempty" mapstructure:"headers"`
	Timeout          int      `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout"` // Timeout is in seconds
	CacheOptions     []string `json:"cacheOptions" yaml:"cacheOptions" mapstructure:"cacheOptions"`
}
