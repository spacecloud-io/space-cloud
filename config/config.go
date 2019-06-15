package config

import "github.com/spaceuptech/space-cloud/utils"

// Config holds the entire configuration
type Config struct {
	Projects map[string]*Project `json:"projects" yaml:"projects"` // The key here is the project id
}

// Project holds the project level configuration
type Project struct {
	ID      string   `json:"id" yaml:"id"`
	Secret  string   `json:"secret" yaml:"secret"`
	Modules *Modules `json:"modules" yaml:"modules"`
	SSL     *SSL     `json:"ssl" yaml:"ssl"`
	Admin   *Admin   `json:"admin" yaml:"admin"`
}

// Admin stores the admin credentials
type Admin struct {
	User string `json:"user" yaml:"user"`
	Pass string `json:"pass" yaml:"pass"`
	Role string `json:"role" yaml:"role"`
}

// SSL holds the certificate and key file locations
type SSL struct {
	Crt string `json:"crt" yaml:"crt"`
	Key string `json:"key" yaml:"key"`
}

// Modules holds the config of all the modules of that environment
type Modules struct {
	Crud      Crud       `json:"crud" yaml:"crud"`
	Auth      Auth       `json:"auth" yaml:"auth"`
	Functions *Functions `json:"functions" yaml:"functions"`
	Realtime  *Realtime  `json:"realtime" yaml:"realtime"`
	FileStore *FileStore `json:"fileStore" yaml:"fileStore"`
	Static    *Static    `json:"static" yaml:"static"`
}

// Crud holds the mapping of database level configuration
type Crud map[string]*CrudStub // The key here is the database type

// CrudStub holds the config at the database level
type CrudStub struct {
	Conn        string                `json:"conn" yaml:"conn"`
	Collections map[string]*TableRule `json:"collections" yaml:"collections"` // The key here is table name
	IsPrimary   bool                  `json:"isPrimary" yaml:"isPrimary"`
}

// TableRule contains the config at the collection level
type TableRule struct {
	IsRealTimeEnabled bool             `json:"isRealtimeEnabled" yaml:"isRealtimeEnabled"`
	Rules             map[string]*Rule `json:"rules" yaml:"rules"` // The key here is query, insert, update or delete
}

// Rule is the authorisation object at the query level
type Rule struct {
	Rule    string                 `json:"rule" yaml:"rule"`
	Eval    string                 `json:"eval" yaml:"eval"`
	Type    string                 `json:"type" yaml:"type"`
	F1      interface{}            `json:"f1" yaml:"f1"`
	F2      interface{}            `json:"f2" yaml:"f2"`
	Clauses []*Rule                `json:"clauses" yaml:"clauses"`
	DB      string                 `json:"db" yaml:"db"`
	Col     string                 `json:"col" yaml:"col"`
	Find    map[string]interface{} `json:"find" yaml:"find"`
	Service string                 `json:"service" yaml:"service"`
	Func    string                 `json:"func" yaml:"func"`
}

// Auth holds the mapping of the sign in method
type Auth map[string]*AuthStub // The key here is the sign in method

// AuthStub holds the config at a single sign in level
type AuthStub struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	ID      string `json:"id" yaml:"id"`
	Secret  string `json:"secret" yaml:"secret"`
}

// Functions holds the config for the Functions module
type Functions struct {
	Enabled bool         `json:"enabled" yaml:"enabled"`
	Broker  utils.Broker `json:"broker" yaml:"broker"`
	Conn    string       `json:"conn" yaml:"conn"`
	Rules   FuncRules    `json:"rules" yaml:"rules"`
}

// FuncRules is the rules for the functions module
type FuncRules map[string]map[string]*Rule // service -> function -> rule

// Realtime holds the config for the realtime module
type Realtime struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Broker  string `json:"broker" yaml:"broker"`
	Conn    string `json:"Conn" yaml:"Conn"`
}

// FileStore holds the config for the file store module
type FileStore struct {
	Enabled   bool                 `json:"enabled" yaml:"enabled"`
	StoreType string               `json:"storeType" yaml:"storeType"`
	Conn      string               `json:"conn" yaml:"conn"`
	Rules     map[string]*FileRule `json:"rules" yaml:"rules"`
}

// FileRule is the authorization object at the file rule level
type FileRule struct {
	Prefix string           `json:"prefix" yaml:"prefix"`
	Rule   map[string]*Rule `json:"rule" yaml:"rule"` // The key can be create, read, delete
}

// Static holds the config for the static files module
type Static struct {
	Enabled bool           `json:"enabled" yaml:"enabled"`
	Gzip    bool           `json:"gzip" yaml:"gzip"`
	Routes  []*StaticRoute `json:"routes" yaml:"routes"`
}

// StaticRoute holds the config for each route
type StaticRoute struct {
	Path      string `json:"path" yaml:"path"`
	URLPrefix string `json:"prefix" yaml:"prefix"`
	Host      string `json:"host" yaml:"host"`
	Proxy     string `json:"proxy" yaml:"proxy"`
}
