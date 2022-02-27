package model

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
