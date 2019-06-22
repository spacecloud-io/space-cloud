package model

// RegisterRequest is the struct which carries the space cloud register payload
type RegisterRequest struct {
	ID      string `json:"id"` // This is the space cloud id
	Secret  string `json:"secret"`
	Account string `json:"account"`
}

// RegisterResponse is the response to the register request
type RegisterResponse struct {
	Ack   bool   `json:"ack"`
	Error string `json:"error"`
}

// ProjectFeed is the body sent to push a project config
type ProjectFeed struct {
	Config  interface{} `json:"config"`
	Project string      `json:"project"`
}
