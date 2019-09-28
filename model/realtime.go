package model

// FeedData is the format to send realtime data
type FeedData struct {
	QueryID   string      `json:"id" structs:"id"`
	DocID     string      `json:"docId" structs:"docId"`
	Type      string      `json:"type" structs:"type"`
	Payload   interface{} `json:"payload" structs:"payload"`
	TimeStamp int64       `json:"time" structs:"time"`
	Group     string      `json:"group" structs:"group"`
	DBType    string      `json:"dbType" structs:"dbType"`
	TypeName  string      `json:"__typename,omitempty" structs:"__typename,omitempty"`
}

// RealtimeRequest is the object sent for realtime requests
type RealtimeRequest struct {
	Token   string                 `json:"token"`
	DBType  string                 `json:"dbType"`
	Project string                 `json:"project"`
	Group   string                 `json:"group"` // Group is the collection name
	Type    string                 `json:"type"`  // Can either be subscribe or unsubscribe
	ID      string                 `json:"id"`    // id is the query id
	Where   map[string]interface{} `json:"where"`
	Options LiveQueryOptions       `json:"options"`
}

// RealtimeResponse is the object sent for realtime requests
type RealtimeResponse struct {
	Group string      `json:"group"` // Group is the collection name
	ID    string      `json:"id"`    // id is the query id
	Ack   bool        `json:"ack"`
	Error string      `json:"error"`
	Docs  []*FeedData `json:"docs"`
}

// Message is the request body of the message
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	ID   string      `json:"id"` // the request id
}

// LiveQueryOptions is to set the options for realtime requests
type LiveQueryOptions struct {
	SkipInitial bool `json:"skipInitial"`
}
