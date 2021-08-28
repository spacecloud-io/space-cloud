package utils

import (
	"golang.org/x/net/context"
)

// BuildVersion is the current version of Space Cloud
const BuildVersion = "0.21.5"

// DLQEventTriggerPrefix used as suffix for DLQ event trigger
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

// DefaultContextTime used for creating default context time for endpoints
const DefaultContextTime = 100

// DefaultCacheTTLTimeout used for setting default ttl expiry if not specified
// the value is in seconds, current value corresponds to 5 minutes
const DefaultCacheTTLTimeout = 60 * 5

// AdminSecretKID describes the kid to be used for admin secrets
const AdminSecretKID = "sc-admin-kid"
