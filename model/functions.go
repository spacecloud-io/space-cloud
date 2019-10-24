package model

// FunctionsRequest is the api call request
type FunctionsRequest struct {
	Params  interface{} `json:"params"`
	Timeout int         `json:"timeout"`
}
