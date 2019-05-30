package model

// FunctionsRequest is the api call request
type FunctionsRequest struct {
	Params  interface{} `json:"params"`
	Timeout int         `json:"timeout"`
}

// FunctionsPayload is the struct transmitted via the broker
type FunctionsPayload struct {
	ID      string                 `json:"id"`
	Auth    map[string]interface{} `json:"auth"`
	Params  interface{}            `json:"params"`
	Service string                 `json:"service"`
	Func    string                 `json:"func"`
	Error   string                 `json:"error"`
}

// ServiceRegisterRequest is the register service request
type ServiceRegisterRequest struct {
	Service string `json:"service"`
}
