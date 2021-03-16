package types

const (
	// All is used when all the records needs to be worked on
	All string = "all"

	// One is used when oly a single record needs to be worked on
	One string = "one"

	// Read is used to Read documents
	Read string = "read"

	Aggregate string = "aggr"

	// Upsert is used to upsert documents
	Batch string = "batch"

	// Count is used to count the number of documents returned
	Count string = "count"

	// Distinct is used to get the distinct values
	Distinct string = "distinct"

	// Upsert is used to upsert documents
	Upsert string = "upsert"

	// Delete is used to delete documents
	Delete string = "delete"

	// Update is used to update documents
	Update string = "update"

	// Create is used to create documents
	Create string = "create"
)

const (
	// TypeRealtimeSubscribe is the request type for live query subscription
	TypeRealtimeSubscribe string = "realtime-subscribe"

	// TypeRealtimeUnsubscribe is the request type for live query subscription
	TypeRealtimeUnsubscribe string = "realtime-unsubscribe"

	// TypeRealtimeFeed is the response type for realtime feed
	TypeRealtimeFeed string = "realtime-feed"
)

const (
	// RealtimeInsert is for create operations
	RealtimeInsert string = "insert"

	// RealtimeUpdate is for update operations
	RealtimeUpdate string = "update"

	// RealtimeDelete is for delete operations
	RealtimeDelete string = "delete"

	// RealtimeInitial is for initial operations
	RealtimeInitial string = "initial"
)
