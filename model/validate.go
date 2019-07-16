package model

// RegisterRequest is the struct which carries the space cloud register payload
type RegisterRequest struct {
	ID     string `json:"id"` // This is the space cloud id
	Key    string `json:"key"`
	UserID string `json:"userId"`
	Mode   int    `json:"mode"`
}

// RegisterResponse is the response to the register request
type RegisterResponse struct {
	Ack   bool   `json:"ack"`
	Error string `json:"error"`
}
