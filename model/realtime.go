package model

// FeedData is the format to send realtime data
type FeedData struct {
	ID        string
	Type      string
	Payload   map[string]interface{}
	TimeStamp int64
	Group     string
}

// RealtimeRequest is the object sent for realtime requests
type RealtimeRequest struct {
	Group string                 `json:"group"` // Group is the collection name
	Type  string                 `json:"type"`  // Can either be subscribe or unsubscribe
	ID    string                 `json:"id"`    // id is the query id
	Where map[string]interface{} `json:"where"`
}

// RealtimeResponse is the object sent for realtime requests
type RealtimeResponse struct {
	Group string `json:"group"` // Group is the collection name
	ID    string `json:"id"`    // id is the query id
}

// Message is the request body of the message
type Message struct {
	Type string
	Data interface{}
	ID   string `json:"id"` // the request id
}
