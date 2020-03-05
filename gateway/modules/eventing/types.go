package eventing

const (
	eventingLogs     string = "event_logs"
	invocationLogs   string = "invocation_logs"
	invocationSchema string = `type invocation_logs {
		_id: ID! @primary
		event_id: ID! @foreign(table: "event_logs", field: "_id")
		invocation_time: DateTime!
		request_payload: String
		response_status_code: Integer
		response_body: String
		remark: String	
	  }`
	eventSchema string = `type event_logs {
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
