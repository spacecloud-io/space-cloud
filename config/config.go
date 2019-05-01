package config

import (
	"strings"
	"encoding/json"
)

// Config holds the entire configuration
type Config struct {
	Projects map[string]*Project `json:"projects" yaml:"projects,omitempty"` // The key here is the project id
}

// Project holds the project level configuration
type Project struct {
	ID      string   `json:"id" yaml:"id"`
	Secret  string   `json:"secret" yaml:"secret"`
	Modules *Modules `json:"modules" yaml:"modules"`
	SSL     *SSL     `json:"ssl" yaml:"ssl"`
}

// SSL holds the certificate and key file locations
type SSL struct {
	Crt string `json:"crt" yaml:"crt"`
	Key string `json:"key" yaml:"key"`
}

func (proj *Project) GoString () string {
	output := strings.Builder{}
	enc := json.NewEncoder(&output)
	enc.SetIndent("", "  ")
	enc.Encode(proj)
	return output.String()
}//-- end func Project.String

// Modules holds the config of all the modules of that environment
type Modules struct {
	Crud      Crud       `json:"crud,omitempty" yaml:"crud,omitempty"`
	Auth      Auth       `json:"auth,omitempty" yaml:"auth,omitempty"`
	FaaS      *FaaS      `json:"faas,omitempty" yaml:"faas,omitempty"`
	Realtime  *Realtime  `json:"realtime,omitempty" yaml:"realtime,omitempty"`
	FileStore *FileStore `json:"fileStore,omitempty" yaml:"fileStore,omitempty"`
}

// Crud holds the mapping of database level configuration
type Crud map[string]*CrudStub // The key here is the database type

// CrudStub holds the config at the database level
type CrudStub struct {
	IsPrimary   bool `json:"isPrimary" yaml:"isPrimary"`
	Conn        string `json:"conn,omitempty" yaml:"conn,omitempty"`
	// see conn.go for ConnConfig
	Connection	*ConnConfig	`json:"connection,omitempty" yaml:"connection,omitempty"`
	Collections map[string]*TableRule `json:"collections,omitempty" yaml:"collections,omitempty"`
	// The key here is table name
}

// TableRule containes the config at the collection level
type TableRule struct {
	IsRealTimeEnabled bool `json:"isRealtimeEnabled" yaml:"isRealtimeEnabled"`
	Rules map[string]*Rule `json:"rules,omitempty" yaml:"rules,omitempty"`
	// The key here is query, insert, update or delete
}

// Rule is the authorisation object at the query level
type Rule struct {
	Rule    string `json:"rule" yaml:"rule"`
	Eval    string `json:"eval,omitempty" yaml:"eval,omitempty"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	F1      interface{} `json:"f1,omitempty" yaml:"f1,omitempty"`
	F2      interface{}  `json:"f2,omitempty" yaml:"f2,omitempty"`
	Clauses []*Rule `json:"clauses,omitempty" yaml:"clauses,omitempty"`
	DB      string `json:"db,omitempty" yaml:"db,omitempty"`
	Col     string `json:"col,omitempty" yaml:"col,omitempty"`
	Find    map[string]interface{} `json:"find,omitempty" yaml:"find,omitempty"`
}

// Auth holds the mapping of the sign in method
type Auth map[string]*AuthStub // The key here is the sign in method

// AuthStub holds the config at a single sign in level
type AuthStub struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	ID      string `json:"id" yaml:"id"`
	Secret  string `json:"secret,omitempty" yaml:"secret,omitempty"`
}

// FaaS holds the config for the FaaS module
type FaaS struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Nats    string `json:"nats" yaml:"nats"`
}

// Realtime holds the config for the realtime module
type Realtime struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Broker  string `json:"broker" yaml:"broker"`
	Conn    string `json:"conn,omitempty" yaml:"conn,omitempty"`
}

// FileStore holds the config for the file store module
type FileStore struct {
	Enabled   bool `json:"enabled" yaml:"enabled"`
	StoreType string `json:"storeType" yaml:"storeType"`
	Conn      string `json:"conn" yaml:"conn"`
	Rules     map[string]*FileRule `json:"rules,omitempty" yaml:"rules,omitempty"`
}

// FileRule is the authorization object at the file rule level
type FileRule struct {
	Prefix string `json:"prefix" yaml:"prefix"`
	Rule   map[string]*Rule `json:"rule,omitempty" yaml:"rule,omitempty"`
	// The key can be create, read, delete
}

