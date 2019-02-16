package config

// Config holds the entire configuration
type Config struct {
	Projects map[string]*Project `json:"projects"` // The key here is the project id
}

// Project holds the project level configuration
type Project struct {
	ID      string `json:"id"`
	Secret  string `json:"secret"`
	Modules *Modules `json:"modules"`
}

// Modules holds the config of all the modules of that environment
type Modules struct {
	Crud      Crud       `json:"crud"`
	Auth      Auth       `json:"auth"`
	FaaS      *FaaS      `json:"faas"`
	Realtime  *Realtime  `json:"realtime"`
	FileStore *FileStore `json:"fileStore"`
}

// Crud holds the mapping of database level configuration
type Crud map[string]*CrudStub // The key here is the database type

// CrudStub holds the config at the database level
type CrudStub struct {
	Connection  string                `json:"conn"`
	Collections map[string]*TableRule `json:"collections"` // The key here is table name
	IsPrimary   bool                  `json:"isPrimary"`
}

// TableRule containes the config at the collection level
type TableRule struct {
	IsRealTimeEnabled bool             `json:"isRealtimeEnabled"`
	Rules             map[string]*Rule `json:"rules"` // The key here is query, insert, update or delete
}

// Rule is the authorisation object at the query level
type Rule struct {
	Rule      string                 `json:"rule"`
	Eval      string                 `json:"eval"`
	FieldType string                 `json:"type"`
	F1        interface{}            `json:"f1"`
	F2        interface{}            `json:"f2"`
	Clauses   []*Rule                `json:"clauses"`
	DbType    string                 `json:"db"`
	Col       string                 `json:"col"`
	Find      map[string]interface{} `json:"find"`
}

// Auth holds the mapping of the sign in method
type Auth map[string]*AuthStub // The key here is the sign in method

// AuthStub holds the config at a single sign in level
type AuthStub struct {
	Enabled bool   `json:"enabled"`
	ID      string `json:"id"`
	Secret  string `json:"secret"`
}

// FaaS holds the config for the FaaS module
type FaaS struct {
	Enabled bool   `json:"enabled"`
	Nats    string `json:"nats"`
}

// Realtime holds the config for the realtime module
type Realtime struct {
	Enabled bool   `json:"enabled"`
	Kafka   string `json:"kafka"`
}

// FileStore holds the config for the file store module
type FileStore struct {
	Enabled    bool                 `json:"enabled"`
	StoreType  string               `json:"storeType"`
	Connection string               `json:"conn"`
	Rules      map[string]*FileRule `json:"rules"`
}

// FileRule is the authorization object at the file rule level
type FileRule struct {
	Prefix string           `json:"prefix"`
	Rule   map[string]*Rule `json:"rule"` // The key can be create, read, delete
}
