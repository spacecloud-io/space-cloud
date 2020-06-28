package utils

const (
	// TableEventingLogs is a variable for "event_logs"
	TableEventingLogs string = "event_logs"
	// TableInvocationLogs is a variable for "invocation_logs"
	TableInvocationLogs string = "invocation_logs"
	// SchemaInvocationLogs is a variable for invocaton schema
	SchemaInvocationLogs string = `type invocation_logs {
		_id: ID! @primary
		event_id: ID! @foreign(table: "event_logs", field: "_id")
		invocation_time: DateTime! @createdAt
		request_payload: String
		response_status_code: Integer
		response_body: String
		error_msg: String
		remark: String	
	  }`
	// SchemaEventLogs is a variable for event schema
	SchemaEventLogs string = `type event_logs {
		_id: ID! @primary
		batchid: String
		type: String
		rule_name: String
		token: Integer
		ts: DateTime
		event_ts: DateTime @createdAt
		payload: String
		status: String
		remark: String
		invocations: [invocation_logs]! @link(table: "invocation_logs", from: "_id", to: "event_id")
	  }`
)
