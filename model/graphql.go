package model

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"op"`
	Variables     map[string]interface{} `json:"variables"`
}
