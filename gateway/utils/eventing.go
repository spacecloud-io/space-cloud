package utils

const (
	TableEventingLogs    string = "event_logs"
	TableInvocationLogs  string = "invocation_logs"
	SchemaInvocationLogs string = `type invocation_logs {
		_id: ID! @primary
		event_id: ID! @foreign(table: "event_logs", field: "_id")
		invocation_time: DateTime!
		request_payload: String
		response_status_code: Integer
		response_body: String
		error_msg: String
		remark: String	
	  }`
	SchemaEventLogs string = `type event_logs {
		_id: ID! @primary
		batchid: String
		type: String
		rule_name: String
		token: Integer
		timestamp: Integer
		event_timestamp: Integer
		payload: String
		status: String
		retries: Integer
		url: String
		remark: String
		invocations: [invocation_logs]! @link(table: "invocation_logs", from: "_id", to: "event_id")
	  }`
)
