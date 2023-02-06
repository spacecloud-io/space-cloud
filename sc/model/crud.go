package model

import "github.com/spacecloud-io/space-cloud/config"

// CreateRequest is the http body received for a create request
type CreateRequest struct {
	Document  interface{} `json:"doc"`
	Operation string      `json:"op"`
	IsBatch   bool        `json:"isBatch"`
}

// ReadRequest is the http body received for a read request
type ReadRequest struct {
	GroupBy     []interface{}            `json:"group"`
	Aggregate   map[string][]string      `json:"aggregate"`
	Find        map[string]interface{}   `json:"find"`
	Operation   string                   `json:"op"`
	Options     *ReadOptions             `json:"options"`
	IsBatch     bool                     `json:"isBatch"`
	Extras      map[string]interface{}   `json:"extras"`
	PostProcess map[string]*PostProcess  `json:"postProcess"`
	MatchWhere  []map[string]interface{} `json:"matchWhere"`
	Cache       *config.ReadCacheOptions `json:"cache"`
}

// ReadOptions is the options required for a read request
type ReadOptions struct {
	// Debug field is used internally to show
	// _query meta data in the graphql
	Debug      bool             `json:"debug"`
	Select     map[string]int32 `json:"select"`
	Sort       []string         `json:"sort"`
	Skip       *int64           `json:"skip"`
	Limit      *int64           `json:"limit"`
	Distinct   *string          `json:"distinct"`
	Join       []*JoinOption    `json:"join"`
	ReturnType string           `json:"returnType"`
	HasOptions bool             `json:"hasOptions"` // used internally
}

// JoinOption describes the way a join needs to be performed
type JoinOption struct {
	// Op can be either All or One
	// This field decides the way the result of join is returned
	// If op is all, the result is returned as an array
	// If op is one, the result is returned as an object
	Op    string                 `json:"Op" mapstructure:"Op"`
	Type  string                 `json:"type" mapstructure:"type"`
	Table string                 `json:"table" mapstructure:"table"`
	As    string                 `json:"as" mapstructure:"as"`
	On    map[string]interface{} `json:"on" mapstructure:"on"`
	Join  []*JoinOption          `json:"join" mapstructure:"join"`
}

// UpdateRequest is the http body received for an update request
type UpdateRequest struct {
	Find      map[string]interface{} `json:"find"`
	Operation string                 `json:"op"`
	Update    map[string]interface{} `json:"update"`
}

// DeleteRequest is the http body received for a delete request
type DeleteRequest struct {
	Find      map[string]interface{} `json:"find"`
	Operation string                 `json:"op"`
}

// PreparedQueryRequest is the http body received for a PreparedQuery request
type PreparedQueryRequest struct {
	Params map[string]interface{} `json:"params"`
	// This field is used internally to show
	// _query meta data in the graphql
	Debug bool
}

// AggregateRequest is the http body received for an aggregate request
type AggregateRequest struct {
	Pipeline  interface{} `json:"pipe"`
	Operation string      `json:"op"`
}

// AllRequest is a union of parameters required in the various requests
type AllRequest struct {
	Col       string                 `json:"col"`
	Document  interface{}            `json:"doc"`
	Operation string                 `json:"op"`
	Find      map[string]interface{} `json:"find"`
	Update    map[string]interface{} `json:"update"`
	Type      string                 `json:"type"`
	DBAlias   string                 `json:"dBAlias"`
	Extras    map[string]interface{} `json:"extras"`
}

// SQLMetaData stores sql query information
type SQLMetaData struct {
	Col       string        `json:"col" structs:"col"`
	SQL       string        `json:"sql" structs:"sql"`
	DbAlias   string        `json:"db" structs:"db"`
	Args      []interface{} `json:"args" structs:"args"`
	QueryTime string        `json:"queryTime" structs:"queryTime"`
}

// BatchRequest is the http body for a batch request
type BatchRequest struct {
	Requests []*AllRequest `json:"reqs"`
}

// DBType is the type of database used for a particular crud operation
type DBType string

// DatabaseCollections stores all callections of sql or postgres or mongo
type DatabaseCollections struct {
	TableName string `db:"table_name" json:"tableName"`
}

const (
	// Mongo is the type used for MongoDB
	Mongo DBType = "mongo"

	// EmbeddedDB is the type used for EmbeddedDB
	EmbeddedDB DBType = "embedded"

	// MySQL is the type used for MySQL
	MySQL DBType = "mysql"

	// Postgres is the type used for PostgresQL
	Postgres DBType = "postgres"

	// SQLServer is the type used for MsSQL
	SQLServer DBType = "sqlserver"

	// DefaultValidate is used for default validation operation
	DefaultValidate = "default"

	// DefaultFetchLimit is the default value to be used as a limit to fetch rows/collection in each read query
	DefaultFetchLimit = 1000
)
