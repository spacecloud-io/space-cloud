package config

// Config holds the entire configuration
type Config struct {
	Projects []*Project `json:"projects" yaml:"projects"` // The key here is the project id
	SSL      *SSL       `json:"ssl" yaml:"ssl"`
	Admin    *Admin     `json:"admin" yaml:"admin"`
}

// ClusterConfig holds the cluster level configuration
type ClusterConfig struct {
	Email         string `json:"email" yaml:"email"`
	EnableMetrics bool   `json:"enableMetrics" yaml:"enableMetrics"`
}

// Project holds the project level configuration
type Project struct {
	Secrets            []*Secret `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	AESKey             string    `json:"aesKey,omitempty" yaml:"aesKey,omitempty"`
	ID                 string    `json:"id,omitempty" yaml:"id,omitempty"`
	Name               string    `json:"name,omitempty" yaml:"name,omitempty"`
	DockerRegistry     string    `json:"dockerRegistry,omitempty" yaml:"dockerRegistry,omitempty"`
	Modules            *Modules  `json:"modules,omitempty" yaml:"modules,omitempty"`
	ContextTimeGraphQL int       `json:"contextTimeGraphQL,omitempty" yaml:"contextTimeGraphQL,omitempty"` // contextTime sets the timeout of query
}

// Secret describes the a secret object
type Secret struct {
	IsPrimary bool   `json:"isPrimary" yaml:"isPrimary"`
	Secret    string `json:"secret" yaml:"secret"`
}

// Admin holds the admin config
type Admin struct {
	ClusterConfig *ClusterConfig `json:"clusterConfig" yaml:"clusterConfig"`
	ClusterID     string         `json:"clusterId" yaml:"clusterId"`
	ClusterKey    string         `json:"clusterKey" yaml:"clusterKey"`
	License       string         `json:"license" yaml:"license"`
}

// AdminUser holds the user credentials and scope
type AdminUser struct {
	User   string `json:"user" yaml:"user"`
	Pass   string `json:"pass" yaml:"pass"`
	Secret string `json:"secret" yaml:"secret"`
}

// SSL holds the certificate and key file locations
type SSL struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Crt     string `json:"crt" yaml:"crt"`
	Key     string `json:"key" yaml:"key"`
}

// Modules holds the config of all the modules of that environment
type Modules struct {
	Crud        Crud            `json:"db" yaml:"db"`
	Auth        Auth            `json:"userMan" yaml:"userMan"`
	Services    *ServicesModule `json:"remoteServices" yaml:"remoteServices"`
	FileStore   *FileStore      `json:"fileStore" yaml:"fileStore"`
	Eventing    Eventing        `json:"eventing,omitempty" yaml:"eventing,omitempty"`
	LetsEncrypt LetsEncrypt     `json:"letsencrypt" yaml:"letsencrypt"`
	Routes      Routes          `json:"ingressRoutes" yaml:"ingressRoutes"`
	Deployments Deployments     `json:"deployments" yaml:"deployments"`
	Secrets     interface{}     `json:"secrets" yaml:"secrets"`
}

// Deployments store all services information for particular project
type Deployments struct {
	Services interface{} `json:"services" yaml:"services"`
}

// Crud holds the mapping of database level configuration
type Crud map[string]*CrudStub // The key here is the alias for database type

// CrudStub holds the config at the database level
type CrudStub struct {
	Type            string                    `json:"type,omitempty" yaml:"type"` // database type
	DBName          string                    `json:"name,omitempty" yaml:"name"` // name of the logical database or schema name according to the database type
	Conn            string                    `json:"conn,omitempty" yaml:"conn"`
	Collections     map[string]*TableRule     `json:"collections,omitempty" yaml:"collections"` // The key here is table name
	PreparedQueries map[string]*PreparedQuery `json:"preparedQueries,omitempty" yaml:"preparedQueries"`
	IsPrimary       bool                      `json:"isPrimary" yaml:"isPrimary"`
	Enabled         bool                      `json:"enabled" yaml:"enabled"`
	BatchTime       int                       `json:"batchTime,omitempty" yaml:"batchTime"`       // time in milli seconds
	BatchRecords    int                       `json:"batchRecords,omitempty" yaml:"batchRecords"` // indicates number of records per batch
}

// PreparedQuery contains the config at the collection level
type PreparedQuery struct {
	ID        string   `json:"id" yaml:"id"`
	SQL       string   `json:"sql" yaml:"sql"`
	Rule      *Rule    `json:"rule" yaml:"rule"`
	Arguments []string `json:"args" yaml:"args"`
}

// TableRule contains the config at the collection level
type TableRule struct {
	IsRealTimeEnabled bool             `json:"isRealtimeEnabled" yaml:"isRealtimeEnabled"`
	Rules             map[string]*Rule `json:"rules" yaml:"rules"` // The key here is query, insert, update or delete
	Schema            string           `json:"schema" yaml:"schema"`
}

// Rule is the authorisation object at the query level
type Rule struct {
	ID      string                 `json:"id,omitempty" yaml:"id,omitempty"`
	Rule    string                 `json:"rule" yaml:"rule"`
	Eval    string                 `json:"eval,omitempty" yaml:"eval,omitempty"`
	Type    string                 `json:"type,omitempty" yaml:"type,omitempty"`
	F1      interface{}            `json:"f1,omitempty" yaml:"f1,omitempty"`
	F2      interface{}            `json:"f2,omitempty" yaml:"f2,omitempty"`
	Clauses []*Rule                `json:"clauses,omitempty" yaml:"clauses,omitempty"`
	DB      string                 `json:"db,omitempty" yaml:"db,omitempty"`
	Col     string                 `json:"col,omitempty" yaml:"col,omitempty"`
	Find    map[string]interface{} `json:"find,omitempty" yaml:"find,omitempty"`
	URL     string                 `json:"url,omitempty" yaml:"url,omitempty"`
	Fields  []string               `json:"fields,omitempty" yaml:"fields,omitempty"`
	Field   string                 `json:"field,omitempty" yaml:"field,omitempty"`
	Value   interface{}            `json:"value,omitempty" yaml:"value,omitempty"`
	Clause  *Rule                  `json:"clause,omitempty" yaml:"clause,omitempty"`
	Name    string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Error   string                 `json:"error,omitempty" yaml:"error,omitempty"`
}

// Auth holds the mapping of the sign in method
type Auth map[string]*AuthStub // The key here is the sign in method

// AuthStub holds the config at a single sign in level
type AuthStub struct {
	ID      string `json:"id" yaml:"id"`
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Secret  string `json:"secret" yaml:"secret"`
}

// ServicesModule holds the config for the service module
type ServicesModule struct {
	Services         Services `json:"externalServices" yaml:"externalServices"`
	InternalServices Services `json:"internalServices" yaml:"internalServices"`
}

// Services holds the config of services
type Services map[string]*Service

// Service holds the config of service
type Service struct {
	ID        string               `json:"id,omitempty" yaml:"id,omitempty"`   // eg. http://localhost:8080
	URL       string               `json:"url,omitempty" yaml:"url,omitempty"` // eg. http://localhost:8080
	Endpoints map[string]*Endpoint `json:"endpoints" yaml:"endpoints"`
}

// Endpoint holds the config of a endpoint
type Endpoint struct {
	Kind      EndpointKind             `json:"kind" yaml:"kind"`
	Tmpl      EndpointTemplatingEngine `json:"template,omitempty" yaml:"template,omitempty"`
	ReqTmpl   string                   `json:"requestTemplate" yaml:"requestTemplate"`
	GraphTmpl string                   `json:"graphTemplate" yaml:"graphTemplate"`
	ResTmpl   string                   `json:"responseTemplate" yaml:"responseTemplate"`
	OpFormat  string                   `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty"`
	Token     string                   `json:"token,omitempty" yaml:"token,omitempty"`
	Method    string                   `json:"method" yaml:"method"`
	Path      string                   `json:"path" yaml:"path"`
	Rule      *Rule                    `json:"rule" yaml:"rule"`
	Headers   []struct {
		Key   string `json:"key" yaml:"key"`
		Value string `json:"value" yaml:"value"`
	} `json:"headers" yaml:"headers"`
}

// EndpointKind describes the type of endpoint. Default value - internal
type EndpointKind string

const (
	// EndpointKindInternal describes a simple or straight forward web-hook call
	EndpointKindInternal EndpointKind = "internal"

	// EndpointKindExternal describes an endpoint on an external server
	EndpointKindExternal EndpointKind = "external"

	// EndpointKindPrepared describes an endpoint on on Space Cloud GraphQL layer
	EndpointKindPrepared EndpointKind = "prepared"
)

// EndpointTemplatingEngine describes the type of endpoint. Default value - go
type EndpointTemplatingEngine string

const (
	// EndpointTemplatingEngineGo describes the go templating engine
	EndpointTemplatingEngineGo EndpointTemplatingEngine = "go"
)

// FileStore holds the config for the file store module
type FileStore struct {
	Enabled   bool        `json:"enabled" yaml:"enabled"`
	StoreType string      `json:"storeType" yaml:"storeType"`
	Conn      string      `json:"conn" yaml:"conn"`
	Endpoint  string      `json:"endpoint" yaml:"endpoint"`
	Bucket    string      `json:"bucket" yaml:"bucket"`
	Secret    string      `json:"secret" yaml:"secret"`
	Rules     []*FileRule `json:"rules,omitempty" yaml:"rules"`
}

// FileRule is the authorization object at the file rule level
type FileRule struct {
	ID     string           `json:"id" yaml:"id"`
	Prefix string           `json:"prefix" yaml:"prefix"`
	Rule   map[string]*Rule `json:"rule" yaml:"rule"` // The key can be create, read, delete
}

// Static holds the config for the static files module
type Static struct {
	Routes         []*StaticRoute `json:"routes" yaml:"routes"`
	InternalRoutes []*StaticRoute `json:"internalRoutes" yaml:"internalRoutes"`
}

// StaticRoute holds the config for each route
type StaticRoute struct {
	ID        string `json:"id,omitempty" yaml:"id,omitempty"`
	Path      string `json:"path" yaml:"path"`
	URLPrefix string `json:"prefix" yaml:"prefix"`
	Host      string `json:"host" yaml:"host"`
	Proxy     string `json:"proxy" yaml:"proxy"`
	Protocol  string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

// Eventing holds the config for the eventing module (task queue)
type Eventing struct {
	Enabled       bool                    `json:"enabled" yaml:"enabled"`
	DBAlias       string                  `json:"dbAlias" yaml:"dbAlias"`
	Rules         map[string]EventingRule `json:"triggers,omitempty" yaml:"triggers"`
	InternalRules map[string]EventingRule `json:"internalTriggers,omitempty" yaml:"internalTriggers,omitempty"`
	SecurityRules map[string]*Rule        `json:"securityRules,omitempty" yaml:"securityRules,omitempty"`
	Schemas       map[string]SchemaObject `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

// EventingRule holds an eventing rule
type EventingRule struct {
	Type    string `json:"type" yaml:"type"`
	Retries int    `json:"retries" yaml:"retries"`
	// Timeout is in milliseconds
	Timeout int               `json:"timeout" yaml:"timeout"`
	ID      string            `json:"id" yaml:"id"`
	URL     string            `json:"url" yaml:"url"`
	Options map[string]string `json:"options" yaml:"options"`
}

// SchemaObject is the body of the request for adding schema
type SchemaObject struct {
	ID     string `json:"id,omitempty" yaml:"id,omitempty"`
	Schema string `json:"schema" yaml:"schema"`
}

// LetsEncrypt describes the configuration for let's encrypt
type LetsEncrypt struct {
	ID                 string   `json:"id,omitempty" yaml:"id,omitempty"`
	WhitelistedDomains []string `json:"domains" yaml:"domains"`
}
