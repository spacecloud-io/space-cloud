package utils

// InternalUserID is the auth.id used for internal requests
const InternalUserID string = "internal-sc-user"

// GraphQLGroupByArgument is used by graphql group clause
const GraphQLGroupByArgument = "group"

// GraphQLAggregate is used by graphql aggregate clause
const GraphQLAggregate = "aggregate"

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

const (
	// TableEventingLogs is a variable for "event_logs"
	TableEventingLogs string = "event_logs"

	// TableInvocationLogs is a variable for "invocation_logs"
	TableInvocationLogs string = "invocation_logs"
)
