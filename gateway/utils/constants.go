package utils

import (
	"golang.org/x/net/context"
)

// BuildVersion is the current version of Space Cloud
const BuildVersion = "0.18.1"

//DLQEventTriggerPrefix used as suffix for DLQ event trigger
const DLQEventTriggerPrefix = "dlq_"

const (
	// One operation returns a single document from the database
	One string = "one"

	// All operation returns multiple documents from the database
	All string = "all"

	// Count operation returns the number of documents which match the condition
	Count string = "count"

	// Distinct operation returns distinct values
	Distinct string = "distinct"

	// Upsert creates a new document if it doesn't exist, else it updates exiting document
	Upsert string = "upsert"
)

// DBType is the type of database used for a particular crud operation
type DBType string

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
)

// Broker is the type of broker used by Space Cloud
type Broker string

const (
	// Nats is the type used for Nats
	Nats Broker = "nats"
)

// FileStoreType is the type of file store used
type FileStoreType string

const (
	// Local is the type used for the local filesystem
	Local FileStoreType = "local"

	// AmazonS3 is the type used for the AmazonS3 storage
	AmazonS3 FileStoreType = "amazon-s3"

	// GCPStorage is the type used for the GCP storage
	GCPStorage FileStoreType = "gcp-storage"
)

// FileOpType is the type of file operation being performed on the file store
type FileOpType string

const (
	// PayloadSize is the size of the payload(in bytes) in file upload and download
	PayloadSize int = 256 * 1024 // 256 kB
)

const (
	// FileRead is the type used for read operations
	FileRead FileOpType = "read"

	// FileCreate is the type used for create operations
	FileCreate FileOpType = "create"

	// FileDelete is the type used for delete operations
	FileDelete FileOpType = "delete"
)

// OperationType is the type of operation being performed on the database
type OperationType string

const (
	// Create is the type used for insert operations
	Create OperationType = "create"

	// Read is the type used for query operation
	Read OperationType = "read"

	// List is the type used for file store list operation
	List OperationType = "list"

	// Update is the type used ofr update operations
	Update OperationType = "update"

	// Delete is the type used for delete operations
	Delete OperationType = "delete"

	// Batch is the type used for batch operations
	Batch OperationType = "batch"

	// Aggregation is the type used for aggregations
	Aggregation OperationType = "aggr"
)

const (

	// RealtimeInsert is for create operations
	RealtimeInsert string = "insert"

	// RealtimeUpdate is for update operations
	RealtimeUpdate string = "update"

	// RealtimeDelete is for delete operations
	RealtimeDelete string = "delete"

	// RealtimeInitial is for the initial data
	RealtimeInitial string = "initial"
)
const (
	// TypeRealtimeSubscribe is the request type for live query subscription
	TypeRealtimeSubscribe string = "realtime-subscribe"

	// TypeRealtimeUnsubscribe is the request type for live query subscription
	TypeRealtimeUnsubscribe string = "realtime-unsubscribe"

	// TypeRealtimeFeed is the response type for realtime feed
	TypeRealtimeFeed string = "realtime-feed"
)

// DefaultConfigFilePath is the default path to load / store the config file
const DefaultConfigFilePath string = "config.yaml"

// OrchestratorType is the type of the orchestrator
type OrchestratorType string

const (
	// Kubernetes is the type used for a kubernetes deployement
	Kubernetes OrchestratorType = "kubernetes"
)

const (
	// PortHTTP is the port used for the http server
	PortHTTP string = "4122"

	// PortHTTPConnect is the port used for the http server with consul connect
	PortHTTPConnect string = "4124"

	// PortHTTPSecure is the port used for the http server with tls
	PortHTTPSecure string = "4126"
)

// InternalUserID is the auth.id used for internal requests
const InternalUserID string = "internal-sc-user"

const (
	// GqlConnectionKeepAlive send every 20 second to client over websocket
	GqlConnectionKeepAlive string = "ka" // Server -> Client
	// GqlConnectionInit is used by graphql over websocket protocol
	GqlConnectionInit string = "connection_init" // Client -> Server
	// GqlConnectionAck is used by graphql over websocket protocol
	GqlConnectionAck string = "connection_ack" // Server -> Client
	// GqlConnectionError is used by graphql over websocket protocol
	GqlConnectionError string = "connection_error" // Server -> Client
	// GqlConnectionTerminate is used by graphql over websocket protocol
	GqlConnectionTerminate string = "connection_terminate" // Client -> Server
	// GqlStart is used by graphql over websocket protocol
	GqlStart string = "start" // Client -> Server
	// GqlData is used by graphql over websocket protocol
	GqlData string = "data" // Server -> Client
	// GqlError is used by graphql over websocket protocol
	GqlError string = "error" // Server -> Client
	// GqlComplete is used by graphql over websocket protocol
	GqlComplete string = "complete" // Server -> Client
	// GqlStop is used by graphql over websocket protocol
	GqlStop string = "stop" // Client -> Server
)

// GraphQLGroupByArgument is used by graphql group clause
const GraphQLGroupByArgument = "group"

// GraphQLAggregate is used by graphql aggregate clause
const GraphQLAggregate = "aggregate"

// FieldType is the type for storing sql inspection information
type FieldType struct {
	FieldName    string `db:"Field"`
	FieldType    string `db:"Type"`
	FieldNull    string `db:"Null"`
	FieldKey     string `db:"Key"`
	FieldDefault string `db:"Default"`
}

// ForeignKeysType is the type for storing  foreignkeys information of sql inspection
type ForeignKeysType struct {
	TableName      string `db:"TABLE_NAME"`
	ColumnName     string `db:"COLUMN_NAME"`
	ConstraintName string `db:"CONSTRAINT_NAME"`
	DeleteRule     string `db:"DELETE_RULE"`
	RefTableName   string `db:"REFERENCED_TABLE_NAME"`
	RefColumnName  string `db:"REFERENCED_COLUMN_NAME"`
}

// IndexType is the type use to indexkey information of sql inspection
type IndexType struct {
	TableName  string `db:"TABLE_NAME"`
	ColumnName string `db:"COLUMN_NAME"`
	IndexName  string `db:"INDEX_NAME"`
	Order      int    `db:"SEQ_IN_INDEX"`
	Sort       string `db:"SORT"`
	IsUnique   string `db:"IS_UNIQUE"`
}

// DatabaseCollections stores all callections of sql or postgres or mongo
type DatabaseCollections struct {
	TableName string `db:"table_name" json:"tableName"`
}

// MaxEventTokens describes the maximum number of event workers
const MaxEventTokens int = 100

const (
	// EventDBCreate is fired for create request
	EventDBCreate string = "DB_INSERT"

	// EventDBUpdate is fired for update request
	EventDBUpdate string = "DB_UPDATE"

	// EventDBDelete is fired for delete request
	EventDBDelete string = "DB_DELETE"

	// EventFileCreate is fired for create request
	EventFileCreate string = "FILE_CREATE"

	// EventFileDelete is fired for delete request
	EventFileDelete string = "FILE_DELETE"
)

const (
	// EventStatusIntent signifies that the event hasn't been staged yet
	EventStatusIntent string = "intent"

	// EventStatusStaged signifies that the event can be processed
	EventStatusStaged string = "staged"

	// EventStatusProcessed signifies that the event has been successfully been processed and can be deleted
	EventStatusProcessed string = "processed"

	// EventStatusFailed signifies that the event has failed and should not be processed
	EventStatusFailed string = "failed"

	// EventStatusCancelled signifies that the event has been cancelled and should not be processed
	EventStatusCancelled string = "cancel"
)

// RequestKind specifies the kind of the request
type RequestKind string

const (
	// RequestKindDirect is used when an http request is to be made directly
	RequestKindDirect RequestKind = "direct"

	// RequestKindConsulConnect is used when an http request is to be made via consul connect
	RequestKindConsulConnect RequestKind = "consul-connect"
)

// SpaceCloudServiceName is the service name space cloud will register itself with in service discovery mechanisms
const SpaceCloudServiceName string = "space-cloud"

// TypeMakeHTTPRequest makes a http request
type TypeMakeHTTPRequest func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error

// GetSecrets gets fileStore and database secrets from runner
type GetSecrets func(project, secretName, key string) (string, error)
