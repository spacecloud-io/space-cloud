package types

// type Meta struct {
// 	Project string `json:"project"`
// 	DbType  string `json:"dbtype"`
// 	Col     string `json:"col"`
// 	Token   string `json:"token"`
// }

type Meta struct {
	DbType, Col string
	Token       string
	Project     string
	Operation   string
}

// CreateRequest is the http body received for a create request
type CreateRequest struct {
	Document  interface{} `json:"doc"`
	Operation string      `json:"op"`
}

// ReadRequest is the http body received for a read request
type ReadRequest struct {
	Find      map[string]interface{} `json:"find"`
	Operation string                 `json:"op"`
	Options   *ReadOptions           `json:"options"`
}

// ReadOptions is the options required for a read request
type ReadOptions struct {
	Select   map[string]int32 `json:"select"`
	Sort     []string         `json:"sort"`
	Skip     *int64           `json:"skip"`
	Limit    *int64           `json:"limit"`
	Distinct *string          `json:"distinct"`
}

// UpdateRequest is the http body received for an update request
type UpdateRequest struct {
	Find      map[string]interface{} `json:"find"`
	Operation string                 `json:"op"`
	Update    map[string]interface{} `json:"update"`
}

// DeleteRequest is the http body received for a delete request
type DeleteRequest struct {
	Find      map[string]interface{} `json:"find"`
	Operation string                 `json:"op"`
}

// AggregateRequest is the http body received for an aggregate request
type AggregateRequest struct {
	Pipeline  interface{} `json:"pipe"`
	Operation string      `json:"op"`
}

// AllRequest is a union of parameters required in the various requests
type AllRequest struct {
	Col       string      `json:"col"`
	Document  interface{} `json:"doc"`
	Operation string      `json:"op"`
	Find      interface{} `json:"find"`
	Update    interface{} `json:"update"`
	Type      string      `json:"type"`
}

// BatchRequest is the http body for a batch request
type BatchRequest struct {
	Requests []AllRequest `json:"reqs"`
}
