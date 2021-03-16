package types

// ServiceRequest is the api call request
type ServiceRequest struct {
	Params  interface{} `json:"params"`
	Timeout int         `json:"timeout"`
}

// WebsocketMessage is the body for a websocket request
type WebsocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	ID   string      `json:"id"` // the request id
}

//FunctionsPayload is the struct transmitted via the broker
type FunctionsPayload struct {
	ID      string                 `json:"id"`
	Auth    map[string]interface{} `json:"auth"`
	Params  interface{}            `json:"params"`
	Service string                 `json:"service"`
	Func    string                 `json:"func"`
	Error   string                 `json:"error"`
}
