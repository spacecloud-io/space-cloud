package utils

// BuildVersion is the current version of Space Cloud
const BuildVersion = "0.11.0"

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

	// Aggregation is the type used for aggregations
	Aggregation OperationType = "aggr"
)

const (
	// RealtimeWorkerCount are the number of goroutines to process realtime data
	RealtimeWorkerCount int = 10

	// FunctionsWorkerCount are the number of goroutines to process functions data
	FunctionsWorkerCount int = 10

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

	// TypeServiceRegister is the request type for service registration
	TypeServiceRegister string = "service-register"

	// TypeServiceUnregister is the request type for service removal
	TypeServiceUnregister string = "service-unregister"

	// TypeServiceRequest is type triggering a service's function
	TypeServiceRequest string = "service-request"
)

// RealTimeProtocol is the type of protocol requested for realtime.
type RealTimeProtocol string

const (
	// Websocket for realtime implementation.
	Websocket RealTimeProtocol = "Websocket"

	// GRPC for realtime implementation.
	GRPC RealTimeProtocol = "GRPC"

	// GRPCService for Service implementation.
	GRPCService RealTimeProtocol = "GRPC-Service"
)

// DefaultConfigFilePath is the default path to load / store the config file
const DefaultConfigFilePath string = "config.yaml"

// RaftCommandType is the type received in the raft commands
type RaftCommandType string

const (
	// RaftCommandSet is used to set a projects config
	RaftCommandSet RaftCommandType = "set"

	// RaftCommandSetDeploy is used to set the deploy config
	RaftCommandSetDeploy RaftCommandType = "set-deploy"

	// RaftCommandSetOperation is used to set the deploy config
	RaftCommandSetOperation RaftCommandType = "set-operation"

	// RaftCommandSetStatic is used to set the deploy config
	RaftCommandSetStatic RaftCommandType = "set-static"

	// RaftCommandAddInternalRouteOperation is used to add internal routes
	RaftCommandAddInternalRouteOperation RaftCommandType = "add-internal-route"

	// RaftCommandDelete is used to delete a projects config
	RaftCommandDelete RaftCommandType = "delete"
)

// RaftSnapshotDirectory is where the snapshot of the log is stored
const RaftSnapshotDirectory string = "raft-store"

// OrchestratorType is the type of the orchestrator
type OrchestratorType string

const (
	// Kubernetes is the type used for a kubernetes deployement
	Kubernetes OrchestratorType = "kubernetes"
)

const (
	// ScopeDeploy is te scope used for the deploy module
	ScopeDeploy string = "deploy"
)

// TypeRegisterRequest is the space cloud register request
const TypeRegisterRequest string = "register"

const (
	// PortHTTP is the port used for the http server
	PortHTTP string = "4122"

	// PortGRPC is the port used for the grpc server
	PortGRPC string = "4124"

	// PortHTTPSecure is the port used for the http server with tls
	PortHTTPSecure string = "4126"

	// PortGRPCSecure is the port used for the grpc server with tls
	PortGRPCSecure string = "4128"

	// PortNatsServer is the port used for nats
	PortNatsServer int = 4222

	// PortNatsCluster is the port used by nats for clustering
	PortNatsCluster int = 4248

	// PortGossip is used for the membership protocol
	PortGossip string = "4232"

	// PortRaft is used internally by raft
	PortRaft string = "4234"
)

// InternalUserID is the auth.id used for internal requests
const InternalUserID string = "internal-sc-user"
