package model

import "encoding/json"

// GraphQLRequest is the payload received in a graphql request
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// ReadRequestKey is the key type for the dataloader
type ReadRequestKey struct {
	DBType     string
	Col        string
	HasOptions bool
	Req        ReadRequest
}

// String returns a guaranteed unique string that can be used to identify an object
func (key ReadRequestKey) String() string {
	data, _ := json.Marshal(key)
	return string(data)
}

// Raw returns the raw, underlaying value of the key
func (key ReadRequestKey) Raw() interface{} {
	return key
}

// GraphQLCallback is used as a callback for graphql requests
type GraphQLCallback func(op interface{}, err error)
