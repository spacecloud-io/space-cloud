package utils

import (
	"golang.org/x/net/context"
)

// BuildVersion is the current version of Space Cloud
const BuildVersion = "0.14.0"

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

	// MySQL is the type used for MySQL
	MySQL DBType = "sql-mysql"

	// Postgres is the type used for PostgresQL
	Postgres DBType = "sql-postgres"

	// SqlServer is the type used for MsSQL
	SqlServer DBType = "sql-sqlserver"
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
	// GQL_CONNECTION_KEEP_ALIVE send every 20 second to client over websocket
	GQL_CONNECTION_KEEP_ALIVE string = "ka" // Server -> Client
	// GQL_CONNECTION_INIT is used by graphql over websocket protocol
	GQL_CONNECTION_INIT string = "connection_init" // Client -> Server
	// GQL_CONNECTION_ACK is used by graphql over websocket protocol
	GQL_CONNECTION_ACK string = "connection_ack" // Server -> Client
	// GQL_CONNECTION_ERROR is used by graphql over websocket protocol
	GQL_CONNECTION_ERROR string = "connection_error" // Server -> Client
	// GQL_CONNECTION_TERMINATE is used by graphql over websocket protocol
	GQL_CONNECTION_TERMINATE string = "connection_terminate" // Client -> Server
	// GQL_START is used by graphql over websocket protocol
	GQL_START string = "start" // Client -> Server
	// GQL_DATA is used by graphql over websocket protocol
	GQL_DATA string = "data" // Server -> Client
	// GQL_ERROR is used by graphql over websocket protocol
	GQL_ERROR string = "error" // Server -> Client
	// GQL_COMPLETE is used by graphql over websocket protocol
	GQL_COMPLETE string = "complete" // Server -> Client
	// GQL_STOP is used by graphql over websocket protocol
	GQL_STOP string = "stop" // Client -> Server
)

// FieldType is the type for storing sql inspection information
type FieldType struct {
	FieldName    string  `db:"Field"`
	FieldType    string  `db:"Type"`
	FieldNull    string  `db:"Null"`
	FieldKey     string  `db:"Key"`
	FieldDefault *string `db:"Default"`
	FieldExtra   string  `db:"Extra"`
}

// ForeignKeysType is the type for storing  foreignkeys information of sql inspection
type ForeignKeysType struct {
	TableName      string `db:"TABLE_NAME"`
	ColumnName     string `db:"COLUMN_NAME"`
	ConstraintName string `db:"CONSTRAINT_NAME"`
	RefTableName   string `db:"REFERENCED_TABLE_NAME"`
	RefColumnName  string `db:"REFERENCED_COLUMN_NAME"`
}

// DatabaseCollections stores all callections of sql or postgres or mongo
type DatabaseCollections struct {
	TableName string `db:"table_name" json:"tableName"`
}

// MaxEventTokens describes the maximum number of event workers
const MaxEventTokens int = 100

const (
	// EventCreate is fired for create request
	EventCreate string = "DB_INSERT"

	// EventUpdate is fired for update request
	EventUpdate string = "DB_UPDATE"

	// EventDelete is fired for delete request
	EventDelete string = "DB_DELETE"
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

type RequestKind string

const (
	// RequestKindDirect is used when an http request is to be made directly
	RequestKindDirect RequestKind = "direct"

	// RequestKindConsulConnect is used when an http request is to be made via consul connect
	RequestKindConsulConnect RequestKind = "consul-connect"
)

// SpaceCloudServiceName is the service name space cloud will register itself with in service discovery mechanisms
const SpaceCloudServiceName string = "space-cloud"

type MakeHttpRequest func(ctx context.Context, method, url, token string, params, vPtr interface{}) error
