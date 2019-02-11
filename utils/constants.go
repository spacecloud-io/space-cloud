package utils

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

// FileStoreType is the type of file store used
type FileStoreType string

const (
	// Local is the type used for the local filesystem
	Local FileStoreType = "local"
)

// OperationType is the type of operation being performed on the database
type OperationType string

const (
	// Create is the type used for insert operations
	Create OperationType = "insert"

	// Read is the type used for query operation
	Read OperationType = "query"

	// Update is the type used ofr update operations
	Update OperationType = "update"

	// Delete is the type used for delete operations
	Delete OperationType = "delete"

	// Aggregation is the type used for aggregations
	Aggregation OperationType = "aggr"
)
