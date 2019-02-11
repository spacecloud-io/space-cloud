package model

// FaaSRequest is the api call request
type FaaSRequest struct {
	Params  interface{} `json:"params"`
	Timeout int         `json:"timeout"`
}
