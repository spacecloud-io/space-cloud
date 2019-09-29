package config

import "github.com/spaceuptech/space-cloud/utils"

// Config holds the entire configuration
type Config struct {
	Projects []*Project `json:"projects" yaml:"projects"` // The key here is the project id
	SSL      *SSL       `json:"ssl" yaml:"ssl"`
	Admin    *Admin     `json:"admin" yaml:"admin"`
	Deploy   Deploy     `json:"deploy" yaml:"deploy"`
	Static   *Static    `json:"gateway" yaml:"gateway"`
	Cluster  string     `json:"cluster" yaml:"cluster"`
	NodeID   string     `json:"nodeId" yaml:"nodeId"`
}

// Deploy holds the deployment environment config
type Deploy struct {
	Orchestrator utils.OrchestratorType `json:"orchestrator,omitempty" yaml:"orchestrator,omitempty"`
	Registry     Registry               `json:"registry,omitempty" yaml:"registry,omitempty"`
	Namespace    string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Enabled      bool                   `json:"enabled" yaml:"enabled"`
}

// Registry holds the details of the registry
type Registry struct {
	URL   string  `json:"url" yaml:"url"`
	ID    string  `json:"id" yaml:"id"`
	Key   string  `json:"key" yaml:"key"`
	Token *string `json:"token,omitempty" yaml:"token,omitempty"`
}

// Project holds the project level configuration
type Project struct {
	Secret  string   `json:"secret" yaml:"secret"`
	ID      string   `json:"id" yaml:"id"`
	Name    string   `json:"name" yaml:"name"`
	Modules *Modules `json:"modules" yaml:"modules"`
}

// Admin stores the admin credentials
type Admin struct {
	Secret    string          `json:"secret" yaml:"secret"`
	Operation OperationConfig `json:"operatiop"`
	Users     []AdminUser     `json:"users" yaml:"users"`
}

// OperationConfig holds the operation mode config
type OperationConfig struct {
	Mode   int    `json:"mode" yaml:"mode"`
	UserID string `json:"userId" yaml:"userId"`
	Key    string `json:"key" yaml:"key"`
}

// AdminUser holds the user credentials and scope
type AdminUser struct {
	User   string       `json:"user" yaml:"user"`
	Pass   string       `json:"pass" yaml:"pass"`
	Scopes ProjectScope `json:"scopes" yaml:"scopes"`
}

// ProjectScope contains the project level scope
type ProjectScope map[string][]string // (project name -> []scopes)

// SSL holds the certificate and key file locations
type SSL struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Crt     string `json:"crt" yaml:"crt"`
	Key     string `json:"key" yaml:"key"`
}

// Modules holds the config of all the modules of that environment
type Modules struct {
	Crud      Crud       `json:"crud" yaml:"crud"`
	Auth      Auth       `json:"auth" yaml:"auth"`
	Functions *Functions `json:"functions" yaml:"functions"`
	FileStore *FileStore `json:"fileStore" yaml:"fileStore"`
	Pubsub    *Pubsub    `json:"pubsub" yaml:"pubsub"`
	Eventing  Eventing   `json:"eventing,omitempty" yaml:"eventing,omitempty"`
}

// Crud holds the mapping of database level configuration
type Crud map[string]*CrudStub // The key here is the database type

// CrudStub holds the config at the database level
type CrudStub struct {
	Conn        string                `json:"conn" yaml:"conn"`
	Collections map[string]*TableRule `json:"collections" yaml:"collections"` // The key here is table name
	IsPrimary   bool                  `json:"isPrimary" yaml:"isPrimary"`
	Enabled     bool                  `json:"enabled" yaml:"enabled"`
}

// TableRule contains the config at the collection level
type TableRule struct {
	IsRealTimeEnabled bool             `json:"isRealtimeEnabled" yaml:"isRealtimeEnabled"`
	Rules             map[string]*Rule `json:"rules" yaml:"rules"` // The key here is query, insert, update or delete
	Schema            string           `json:"schema" yaml:"schema"`
}

// Rule is the authorisation object at the query level
type Rule struct {
	Rule    string                 `json:"rule" yaml:"rule"`
	Eval    string                 `json:"eval,omitempty" yaml:"eval,omitempty"`
	Type    string                 `json:"type,omitempty" yaml:"type,omitempty"`
	F1      interface{}            `json:"f1,omitempty" yaml:"f1,omitempty"`
	F2      interface{}            `json:"f2,omitempty" yaml:"f2,omitempty"`
	Clauses []*Rule                `json:"clauses,omitempty" yaml:"clauses,omitempty"`
	DB      string                 `json:"db,omitempty" yaml:"db,omitempty"`
	Col     string                 `json:"col,omitempty" yaml:"col,omitempty"`
	Find    map[string]interface{} `json:"find,omitempty" yaml:"find,omitempty"`
	Service string                 `json:"service,omitempty" yaml:"service,omitempty"`
	Func    string                 `json:"func,omitempty" yaml:"func,omitempty"`
}

// Auth holds the mapping of the sign in method
type Auth map[string]*AuthStub // The key here is the sign in method

// AuthStub holds the config at a single sign in level
type AuthStub struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	ID      string `json:"id" yaml:"id"`
	Secret  string `json:"secret" yaml:"secret"`
}

// Functions holds the config for the functions module
type Functions struct {
	Services Services `json:"services" yaml:"services"`
}

// Services holds the config of services
type Services map[string]*Service

// Service holds the config of service
type Service struct {
	Functions map[string]Function `json:"functions" yaml:"functions"`
}

// Function holds the config of a function
type Function struct {
	Rule *Rule `json:"rule" yaml:"rule"`
}

// FileStore holds the config for the file store module
type FileStore struct {
	Enabled   bool        `json:"enabled" yaml:"enabled"`
	StoreType string      `json:"storeType" yaml:"storeType"`
	Conn      string      `json:"conn" yaml:"conn"`
	Endpoint  string      `json:"endpoint" yaml:"endpoint"`
	Bucket    string      `json:"bucket" yaml:"bucket"`
	Rules     []*FileRule `json:"rules" yaml:"rules"`
}

// FileRule is the authorization object at the file rule level
type FileRule struct {
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

// Pubsub holds the config for the realtime module
type Pubsub struct {
	Enabled bool          `json:"enabled" yaml:"enabled"`
	Broker  utils.Broker  `json:"broker" yaml:"broker"`
	Conn    string        `json:"conn" yaml:"conn"`
	Rules   []*PubsubRule `json:"rules" yaml:"rules"`
}

// PubsubRule is the authorization object at the pubsub rule level
type PubsubRule struct {
	Subject string           `json:"subject" yaml:"subject"` // The channel subject
	Rule    map[string]*Rule `json:"rule" yaml:"rule"`       // The key can be publish or subscribe
}

// Eventing holds the config for the eventing module (task queue)
type Eventing struct {
	Enabled       bool                    `json:"enabled" yaml:"enabled"`
	DBType        string                  `json:"dbType" yaml:"dbType"`
	Col           string                  `json:"col" yaml:"col"`
	Rules         map[string]EventingRule `json:"rules" yaml:"rules"`
	InternalRules map[string]EventingRule `json:"internalRules,omitempty" yaml:"internalRules,omitempty"`
}

// EventingRule holds an eventing rule
type EventingRule struct {
	Type     string            `json:"type" yaml:"type"`
	Retries  int               `json:"retries" yaml:"retries"`
	Service  string            `json:"service" yaml:"service"`
	Function string            `json:"func" yaml:"func"`
	Options  map[string]string `json:"options" yaml:"options"`
}
