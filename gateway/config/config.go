package config

import (
	"net/http"
)

// Config holds the entire configuration
type Config struct {
	Projects         Projects         `json:"projects" yaml:"projects" mapstructure:"projects"` // The key here is the project id
	SSL              *SSL             `json:"ssl" yaml:"ssl" mapstructure:"ssl"`
	ClusterConfig    *ClusterConfig   `json:"clusterConfig" yaml:"clusterConfig" mapstructure:"clusterConfig"`
	Integrations     Integrations     `json:"integrations" yaml:"integrations" mapstructure:"integrations"`
	IntegrationHooks IntegrationHooks `json:"integrationsHooks" yaml:"integrationsHooks" mapstructure:"integrationsHooks"`
	License          *License         `json:"license" yaml:"license" mapstructure:"license"`
}

// ClusterConfig holds the cluster level configuration
type ClusterConfig struct {
	LetsEncryptEmail string `json:"letsencryptEmail" yaml:"letsencryptEmail" mapstructure:"letsencryptEmail"`
	EnableTelemetry  bool   `json:"enableTelemetry" yaml:"enableTelemetry" mapstructure:"enableTelemetry"`
}

// License holds license information
type License struct {
	LicenseKey   string `json:"licenseKey" yaml:"licenseKey" mapstructure:"licenseKey"`
	LicenseValue string `json:"licenseValue" yaml:"licenseValue" mapstructure:"licenseValue"`
	License      string `json:"license" yaml:"license" mapstructure:"license"`
}

// Projects is a map which stores config information of all project in a cluster
type Projects map[string]*Project // Key here is project id

// DatabaseConfigs is a map which stores database config information
type DatabaseConfigs map[string]*DatabaseConfig // Key here is resource id --> clusterId--projectId--resourceType--dbAlias

// DatabaseSchemas is a map which stores database schema information
type DatabaseSchemas map[string]*DatabaseSchema // Key here is resource id --> clusterId--projectId--resourceType--dbAlias-tableName

// DatabaseRules is a map which stores database rules information
type DatabaseRules map[string]*DatabaseRule // Key here is resource id --> clusterId--projectId--resourceType--dbAlias-tableName-rule

// DatabasePreparedQueries is a map which stores database prepared query information
type DatabasePreparedQueries map[string]*DatbasePreparedQuery // Key here is resource id --> clusterId--projectId--resourceType--dbAlias-prepareQueryId

// EventingSchemas is a map which stores eventing schema information
type EventingSchemas map[string]*EventingSchema // Key here is resource id --> clusterId--projectId--resourceType--schemaId

// EventingRules is a map which stores database config information
type EventingRules map[string]*Rule // Key here is resource id --> clusterId--projectId--resourceType--ruleId

// EventingTriggers is a map which stores database config information
type EventingTriggers map[string]*EventingTrigger // Key here is resource id --> clusterId--projectId--resourceType--triggerId

// FileStoreRules is a map which stores database config information
type FileStoreRules map[string]*FileRule // Key here is resource id --> clusterId--projectId--resourceType--fileRuleId

// IngressRoutes is a map which stores database config information
type IngressRoutes map[string]*Route // Key here is resource id --> clusterId--projectId--resourceType--routeId

// Project holds the project level configuration
type Project struct {
	ProjectConfig *ProjectConfig `json:"projectConfig" yaml:"projectConfig" mapstructure:"projectConfig"`

	DatabaseConfigs         DatabaseConfigs         `json:"dbConfigs" yaml:"dbConfigs" mapstructure:"dbConfigs"`
	DatabaseSchemas         DatabaseSchemas         `json:"dbSchemas" yaml:"dbSchemas" mapstructure:"dbSchemas"`
	DatabaseRules           DatabaseRules           `json:"dbRules" yaml:"dbRules" mapstructure:"dbRules"`
	DatabasePreparedQueries DatabasePreparedQueries `json:"dbPreparedQuery" yaml:"dbPreparedQuery" mapstructure:"dbPreparedQuery"`

	EventingConfig   *EventingConfig  `json:"eventingConfig" yaml:"eventingConfig" mapstructure:"eventingConfig"`
	EventingSchemas  EventingSchemas  `json:"eventingSchemas" yaml:"eventingSchemas" mapstructure:"eventingSchemas"`
	EventingRules    EventingRules    `json:"eventingRules" yaml:"eventingRules" mapstructure:"eventingRules"`
	EventingTriggers EventingTriggers `json:"eventingTriggers" yaml:"eventingTriggers" mapstructure:"eventingTriggers"`

	FileStoreConfig *FileStoreConfig `json:"fileStoreConfig" yaml:"fileStoreConfig" mapstructure:"fileStoreConfig"`
	FileStoreRules  FileStoreRules   `json:"fileStoreRules" yaml:"fileStoreRules" mapstructure:"fileStoreRules"`

	Auths Auths `json:"auths" yaml:"auths" mapstructure:"auths"`

	LetsEncrypt *LetsEncrypt `json:"letsencrypt" yaml:"letsencrypt" mapstructure:"letsencrypt"`

	IngressRoutes IngressRoutes       `json:"ingressRoute" yaml:"ingressRoute" mapstructure:"ingressRoute"`
	IngressGlobal *GlobalRoutesConfig `json:"ingressGlobal" yaml:"ingressGlobal" mapstructure:"ingressGlobal"`

	RemoteService Services `json:"remoteServices" yaml:"remoteServices" mapstructure:"remoteServices"`
}

// ProjectConfig stores information of individual project
type ProjectConfig struct {
	ID                 string    `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`
	Name               string    `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name"`
	Secrets            []*Secret `json:"secrets,omitempty" yaml:"secrets,omitempty" mapstructure:"secrets"`
	SecretSource       string    `json:"secretSource,omitempty" yaml:"secretSource,omitempty" mapstructure:"secretSource"`
	IsIntegration      bool      `json:"isIntegration,omitempty" yaml:"isIntegration,omitempty" mapstructure:"isIntegration"`
	AESKey             string    `json:"aesKey,omitempty" yaml:"aesKey,omitempty" mapstructure:"aesKey"`
	DockerRegistry     string    `json:"dockerRegistry,omitempty" yaml:"dockerRegistry,omitempty" mapstructure:"dockerRegistry"`
	ContextTimeGraphQL int       `json:"contextTimeGraphQL,omitempty" yaml:"contextTimeGraphQL,omitempty" mapstructure:"contextTimeGraphQL"` // contextTime sets the timeout of query
}

// DriverConfig stores the parameters for drivers of Databases.
type DriverConfig struct {
	MaxConn        int    `json:"maxConn,omitempty" yaml:"maxConn,omitempty" mapstructure:"maxConn"`                      // for SQL and Mongo
	MaxIdleTimeout int    `json:"maxIdleTimeout,omitempty" yaml:"maxIdleTimeout,omitempty" mapstructure:"maxIdleTimeout"` // for SQL and Mongo
	MinConn        uint64 `json:"minConn,omitempty" yaml:"minConn,omitempty" mapstructure:"minConn"`                      // only for Mongo
	MaxIdleConn    int    `json:"maxIdleConn,omitempty" yaml:"maxIdleConn,omitempty" mapstructure:"maxIdleConn"`          // only for SQL
}

// DatabaseConfig stores information of database config
type DatabaseConfig struct {
	DbAlias      string       `json:"dbAlias,omitempty" yaml:"dbAlias" mapstructure:"dbAlias"`
	Type         string       `json:"type,omitempty" yaml:"type" mapstructure:"type"` // database type
	DBName       string       `json:"name,omitempty" yaml:"name" mapstructure:"name"` // name of the logical database or schema name according to the database type
	Conn         string       `json:"conn,omitempty" yaml:"conn" mapstructure:"conn"`
	IsPrimary    bool         `json:"isPrimary" yaml:"isPrimary" mapstructure:"isPrimary"`
	Enabled      bool         `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	BatchTime    int          `json:"batchTime,omitempty" yaml:"batchTime" mapstructure:"batchTime"`          // time in milli seconds
	BatchRecords int          `json:"batchRecords,omitempty" yaml:"batchRecords" mapstructure:"batchRecords"` // indicates number of records per batch
	Limit        int64        `json:"limit,omitempty" yaml:"limit" mapstructure:"limit"`                      // indicates number of records to send per request
	DriverConf   DriverConfig `json:"driverConf,omitempty" yaml:"driverConf" mapstructure:"driverConf"`
}

// DatabaseSchema stores information of db schemas
type DatabaseSchema struct {
	Table   string `json:"col,omitempty" yaml:"col" mapstructure:"col"`
	DbAlias string `json:"dbAlias,omitempty" yaml:"dbAlias" mapstructure:"dbAlias"`
	Schema  string `json:"schema,omitempty" yaml:"schema" mapstructure:"schema"`
}

// DatabaseRule stores information of db rule
type DatabaseRule struct {
	Table             string           `json:"col,omitempty" yaml:"col" mapstructure:"col"`
	DbAlias           string           `json:"dbAlias,omitempty" yaml:"dbAlias" mapstructure:"dbAlias"`
	IsRealTimeEnabled bool             `json:"isRealtimeEnabled,omitempty" yaml:"isRealtimeEnabled" mapstructure:"isRealtimeEnabled"`
	Rules             map[string]*Rule `json:"rules,omitempty" yaml:"rules" mapstructure:"rules"`
}

// EventingConfig stores information of eventing config
type EventingConfig struct {
	Enabled       bool             `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	DBAlias       string           `json:"dbAlias" yaml:"dbAlias" mapstructure:"dbAlias"`
	InternalRules EventingTriggers `json:"internalRules,omitempty" yaml:"internalRules,omitempty" mapstructure:"internalRules"`
}

// EventingSchema stores information of eventing schema
type EventingSchema struct {
	ID     string `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`
	Schema string `json:"schema" yaml:"schema" mapstructure:"schema"`
}

// FileStoreConfig stores information of file store config
type FileStoreConfig struct {
	Enabled        bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	StoreType      string `json:"storeType" yaml:"storeType" mapstructure:"storeType"`
	Conn           string `json:"conn" yaml:"conn" mapstructure:"conn"`
	Endpoint       string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Bucket         string `json:"bucket" yaml:"bucket" mapstructure:"bucket"`
	Secret         string `json:"secret" yaml:"secret" mapstructure:"secret"`
	DisableSSL     *bool  `json:"disableSSL,omitempty" yaml:"disableSSL,omitempty" mapstructure:"disableSSL"`
	ForcePathStyle *bool  `json:"forcePathStyle,omitempty" yaml:"forcePathStyle,omitempty" mapstructure:"forcePathStyle"`
}

// Secret describes the a secret object
type Secret struct {
	IsPrimary bool   `json:"isPrimary" yaml:"isPrimary" mapstructure:"isPrimary"` // used by the frontend & backend to generate token out of multiple secrets
	Alg       JWTAlg `json:"alg" yaml:"alg" mapstructure:"alg"`                   // RSA256 or HMAC256

	KID string `json:"kid" yaml:"kid" mapstructure:"kid"` // uniquely identifies a secret

	JwkURL string      `json:"jwkUrl" yaml:"jwkUrl" mapstructure:"jwkUrl"`
	JwkKey interface{} `json:"-" yaml:"-"`

	Audience []string `json:"aud" yaml:"aud" mapstructure:"aud"`
	Issuer   []string `json:"iss" yaml:"iss" mapstructure:"iss"`

	// Used for HMAC256 secret
	Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`

	// Use for RSA256
	PublicKey  string `json:"publicKey" yaml:"publicKey" mapstructure:"publicKey"`
	PrivateKey string `json:"privateKey" yaml:"privateKey" mapstructure:"privateKey"`
}

// JWTAlg is type of method used for signing token
type JWTAlg string

const (
	// HS256 is method used for signing token
	HS256 JWTAlg = "HS256"

	// RS256 is method used for signing token
	RS256 JWTAlg = "RS256"

	// JwkURL is the method for identifying a secret that has to be validated against secret kes fetched from url
	JwkURL JWTAlg = "JWK_URL"

	// RS256Public is the method for identifying a secret that has to be validated against with a public key
	RS256Public JWTAlg = "RS256_PUBLIC"
)

// Admin holds the admin config
type Admin struct {
	ClusterConfig *ClusterConfig `json:"clusterConfig" yaml:"clusterConfig" mapstructure:"clusterConfig"`
	LicenseKey    string         `json:"licenseKey" yaml:"licenseKey" mapstructure:"licenseKey"`
	LicenseValue  string         `json:"licenseValue" yaml:"licenseValue" mapstructure:"licenseValue"`
	License       string         `json:"license" yaml:"license" mapstructure:"license"`
	Integrations  Integrations   `json:"integrations" yaml:"integrations" mapstructure:"integrations"`
}

// AdminUser holds the user credentials and scope
type AdminUser struct {
	User   string `json:"user" yaml:"user" mapstructure:"user"`
	Pass   string `json:"pass" yaml:"pass" mapstructure:"pass"`
	Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`
}

// SSL holds the certificate and key file locations
type SSL struct {
	Enabled bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Crt     string `json:"crt" yaml:"crt" mapstructure:"crt"`
	Key     string `json:"key" yaml:"key" mapstructure:"key"`
}

// Deployments store all services information for particular project
type Deployments struct {
	Services interface{} `json:"services" yaml:"services" mapstructure:"services"`
}

// Crud holds the mapping of database level configuration
type Crud map[string]*CrudStub // The key here is the alias for database type

// CrudStub holds the config at the database level
type CrudStub struct {
	Type            string                           `json:"type,omitempty" yaml:"type" mapstructure:"type"` // database type
	DBName          string                           `json:"name,omitempty" yaml:"name" mapstructure:"name"` // name of the logical database or schema name according to the database type
	Conn            string                           `json:"conn,omitempty" yaml:"conn" mapstructure:"conn"`
	Collections     map[string]*TableRule            `json:"collections,omitempty" yaml:"collections" mapstructure:"collections"` // The key here is table name
	PreparedQueries map[string]*DatbasePreparedQuery `json:"preparedQueries,omitempty" yaml:"preparedQueries" mapstructure:"preparedQueries"`
	IsPrimary       bool                             `json:"isPrimary" yaml:"isPrimary" mapstructure:"isPrimary"`
	Enabled         bool                             `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	BatchTime       int                              `json:"batchTime,omitempty" yaml:"batchTime" mapstructure:"batchTime"`          // time in milli seconds
	BatchRecords    int                              `json:"batchRecords,omitempty" yaml:"batchRecords" mapstructure:"batchRecords"` // indicates number of records per batch
	Limit           int64                            `json:"limit,omitempty" yaml:"limit" mapstructure:"limit"`                      // indicates number of records per batch
}

// DatbasePreparedQuery stores information of prepared query
type DatbasePreparedQuery struct {
	ID        string   `json:"id" yaml:"id" mapstructure:"id"`
	SQL       string   `json:"sql" yaml:"sql" mapstructure:"sql"`
	Rule      *Rule    `json:"rule" yaml:"rule" mapstructure:"rule"`
	DbAlias   string   `json:"dbAlias" yaml:"dbAlias" mapstructure:"dbAlias"`
	Arguments []string `json:"args" yaml:"args" mapstructure:"args"`
}

// TableRule contains the config at the collection level
type TableRule struct {
	IsRealTimeEnabled bool             `json:"isRealtimeEnabled,omitempty" yaml:"isRealtimeEnabled" mapstructure:"isRealtimeEnabled"`
	Rules             map[string]*Rule `json:"rules,omitempty" yaml:"rules" mapstructure:"rules"` // The key here is query, insert, update or delete
	Schema            string           `json:"schema,omitempty" yaml:"schema" mapstructure:"schema"`
}

// Rule is the authorisation object at the query level
type Rule struct {
	ID       string                 `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`
	Rule     string                 `json:"rule" yaml:"rule" mapstructure:"rule"`
	Eval     string                 `json:"eval,omitempty" yaml:"eval,omitempty" mapstructure:"eval"`
	Type     string                 `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type"`
	F1       interface{}            `json:"f1,omitempty" yaml:"f1,omitempty" mapstructure:"f1"`
	F2       interface{}            `json:"f2,omitempty" yaml:"f2,omitempty" mapstructure:"f2"`
	Clauses  []*Rule                `json:"clauses,omitempty" yaml:"clauses,omitempty" mapstructure:"clauses"`
	DB       string                 `json:"db,omitempty" yaml:"db,omitempty" mapstructure:"db"`
	Col      string                 `json:"col,omitempty" yaml:"col,omitempty" mapstructure:"col"`
	Find     map[string]interface{} `json:"find,omitempty" yaml:"find,omitempty" mapstructure:"find"`
	URL      string                 `json:"url,omitempty" yaml:"url,omitempty" mapstructure:"url"`
	Fields   interface{}            `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields"`
	Field    string                 `json:"field,omitempty" yaml:"field,omitempty" mapstructure:"field"`
	Value    interface{}            `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value"`
	Clause   *Rule                  `json:"clause,omitempty" yaml:"clause,omitempty" mapstructure:"clause"`
	Name     string                 `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name"`
	Error    string                 `json:"error,omitempty" yaml:"error,omitempty" mapstructure:"error"`
	Store    string                 `json:"store,omitempty" yaml:"store,omitempty" mapstructure:"store"`
	Claims   map[string]interface{} `json:"claims,omitempty" yaml:"claims,omitempty" mapstructure:"claims"`
	Template TemplatingEngine       `json:"template,omitempty" yaml:"template,omitempty" mapstructure:"template"`
	ReqTmpl  string                 `json:"requestTemplate,omitempty" yaml:"requestTemplate,omitempty" mapstructure:"requestTemplate"`
	OpFormat string                 `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" mapstructure:"outputFormat"`
}

// Auths holds the mapping of the sign in method
type Auths map[string]*AuthStub // The key here is the sign in method

// AuthStub holds the config at a single sign in level
type AuthStub struct {
	ID      string `json:"id" yaml:"id" mapstructure:"id"`
	Enabled bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Secret  string `json:"secret" yaml:"secret" mapstructure:"secret"`
}

// ServicesModule holds the config for the service module
type ServicesModule struct {
	Services         Services `json:"externalServices" yaml:"externalServices" mapstructure:"externalServices"`
	InternalServices Services `json:"internalServices" yaml:"internalServices" mapstructure:"internalServices"`
}

// Services holds the config of services
type Services map[string]*Service

// Service holds the config of service
type Service struct {
	ID        string               `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`    // eg. http://localhost:8080
	URL       string               `json:"url,omitempty" yaml:"url,omitempty" mapstructure:"url"` // eg. http://localhost:8080
	Endpoints map[string]*Endpoint `json:"endpoints,omitempty" yaml:"endpoints,omitempty" mapstructure:"endpoints"`
}

// Endpoint holds the config of a endpoint
type Endpoint struct {
	Kind      EndpointKind     `json:"kind" yaml:"kind" mapstructure:"kind"`
	Tmpl      TemplatingEngine `json:"template,omitempty" yaml:"template,omitempty" mapstructure:"template"`
	ReqTmpl   string           `json:"requestTemplate" yaml:"requestTemplate" mapstructure:"requestTemplate"`
	GraphTmpl string           `json:"graphTemplate" yaml:"graphTemplate" mapstructure:"graphTemplate"`
	ResTmpl   string           `json:"responseTemplate" yaml:"responseTemplate" mapstructure:"responseTemplate"`
	OpFormat  string           `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" mapstructure:"outputFormat"`
	Token     string           `json:"token,omitempty" yaml:"token,omitempty" mapstructure:"token"`
	Claims    string           `json:"claims,omitempty" yaml:"claims,omitempty" mapstructure:"claims"`
	Method    string           `json:"method" yaml:"method" mapstructure:"method"`
	Path      string           `json:"path" yaml:"path" mapstructure:"path"`
	Rule      *Rule            `json:"rule,omitempty" yaml:"rule,omitempty" mapstructure:"rule"`
	Headers   Headers          `json:"headers,omitempty" yaml:"headers,omitempty" mapstructure:"headers"`
	Timeout   int              `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout"` // Timeout is in seconds
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

// TemplatingEngine describes the type of endpoint. Default value - go
type TemplatingEngine string

const (
	// TemplatingEngineGo describes the go templating engine
	TemplatingEngineGo TemplatingEngine = "go"
)

// Header describes the operation to be performed on the header
type Header struct {
	Key   string `json:"key" yaml:"key" mapstructure:"key"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
	Op    string `json:"op" yaml:"op" mapstructure:"op"`
}

// Headers describes an array of headers
type Headers []Header

// UpdateHeader updated the header values
func (headers Headers) UpdateHeader(reqHeaders http.Header) {
	for _, h := range headers {
		switch h.Op {
		case "", "set":
			reqHeaders.Set(h.Key, h.Value)
		case "add":
			reqHeaders.Add(h.Key, h.Value)
		case "del":
			reqHeaders.Del(h.Key)
		}
	}
}

// FileStore holds the config for the file store module
type FileStore struct {
	Enabled        bool        `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	StoreType      string      `json:"storeType" yaml:"storeType" mapstructure:"storeType"`
	Conn           string      `json:"conn" yaml:"conn" mapstructure:"conn"`
	Endpoint       string      `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Bucket         string      `json:"bucket" yaml:"bucket" mapstructure:"bucket"`
	Secret         string      `json:"secret" yaml:"secret" mapstructure:"secret"`
	Rules          []*FileRule `json:"rules,omitempty" yaml:"rules" mapstructure:"rules"`
	DisableSSL     *bool       `json:"disableSSL,omitempty" yaml:"disableSSL,omitempty" mapstructure:"disableSSL"`
	ForcePathStyle *bool       `json:"forcePathStyle,omitempty" yaml:"forcePathStyle,omitempty" mapstructure:"forcePathStyle"`
}

// FileRule is the authorization object at the file rule level
type FileRule struct {
	ID     string           `json:"id" yaml:"id" mapstructure:"id"`
	Prefix string           `json:"prefix" yaml:"prefix" mapstructure:"prefix"`
	Rule   map[string]*Rule `json:"rule" yaml:"rule" mapstructure:"rule"` // The key can be create, read, delete
}

// Static holds the config for the static files module
type Static struct {
	Routes         []*StaticRoute `json:"routes" yaml:"routes" mapstructure:"routes"`
	InternalRoutes []*StaticRoute `json:"internalRoutes" yaml:"internalRoutes" mapstructure:"internalRoutes"`
}

// StaticRoute holds the config for each route
type StaticRoute struct {
	ID        string `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`
	Path      string `json:"path" yaml:"path" mapstructure:"path"`
	URLPrefix string `json:"prefix" yaml:"prefix" mapstructure:"prefix"`
	Host      string `json:"host" yaml:"host" mapstructure:"host"`
	Proxy     string `json:"proxy" yaml:"proxy" mapstructure:"proxy"`
	Protocol  string `json:"protocol,omitempty" yaml:"protocol,omitempty" mapstructure:"protocol"`
}

// Eventing holds the config for the eventing module (task queue)
type Eventing struct {
	Enabled       bool                        `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	DBAlias       string                      `json:"dbAlias" yaml:"dbAlias" mapstructure:"dbAlias"`
	Rules         map[string]*EventingTrigger `json:"triggers,omitempty" yaml:"triggers" mapstructure:"triggers"`
	InternalRules map[string]*EventingTrigger `json:"internalTriggers,omitempty" yaml:"internalTriggers,omitempty" mapstructure:"internalTriggers"`
	SecurityRules map[string]*Rule            `json:"securityRules,omitempty" yaml:"securityRules,omitempty" mapstructure:"securityRules"`
	Schemas       map[string]SchemaObject     `json:"schemas,omitempty" yaml:"schemas,omitempty" mapstructure:"schemas"`
}

// EventingTrigger stores information of eventing trigger
type EventingTrigger struct {
	Type            string            `json:"type" yaml:"type" mapstructure:"type"`
	Retries         int               `json:"retries" yaml:"retries" mapstructure:"retries"`
	Timeout         int               `json:"timeout" yaml:"timeout" mapstructure:"timeout"` // Timeout is in milliseconds
	ID              string            `json:"id" yaml:"id" mapstructure:"id"`
	URL             string            `json:"url" yaml:"url" mapstructure:"url"`
	Options         map[string]string `json:"options" yaml:"options" mapstructure:"options"`
	Tmpl            TemplatingEngine  `json:"template,omitempty" yaml:"template,omitempty" mapstructure:"template"`
	RequestTemplate string            `json:"requestTemplate,omitempty" yaml:"requestTemplate,omitempty" mapstructure:"requestTemplate"`
	OpFormat        string            `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" mapstructure:"outputFormat"`
	Claims          string            `json:"claims" yaml:"claims" mapstructure:"claims"`
	Filter          *Rule             `json:"filter" yaml:"filter" mapstructure:"filter"`
}

// SchemaObject is the body of the request for adding schema
type SchemaObject struct {
	ID     string `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`
	Schema string `json:"schema" yaml:"schema" mapstructure:"schema"`
}

// LetsEncrypt describes the configuration for let's encrypt
type LetsEncrypt struct {
	ID                 string   `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`
	WhitelistedDomains []string `json:"domains" yaml:"domains" mapstructure:"domains"`
}
