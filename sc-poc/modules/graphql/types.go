package graphql

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/invopop/jsonschema"
)

type (
	m map[string]interface{}
	t []interface{}
)

// CompiledQuery stores the result of a compiled graphql query
type CompiledQuery struct {
	// Fields related to authentication
	IsAuthRequired bool
	InjectedClaims map[string]string

	// Variable & response type defs
	VariableSchema *jsonschema.Schema
	ResponseSchema *jsonschema.Schema
	Extensions     map[string]interface{}

	// Graphql ast
	Query         string
	OperationName string
	OperationType string
	DefaultValues map[string]interface{}
	Doc           *ast.Document
}

type request struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type (
	introspectionResponse struct {
		Data introspectionResponseData `json:"data"`
	}

	introspectionResponseData struct {
		Schema introspectionResponseSchema `json:"__schema"`
	}

	introspectionResponseSchema struct {
		QueryType        *introspectionResponseTypeRef     `json:"queryType"`
		MutationType     *introspectionResponseTypeRef     `json:"mutationType"`
		SubscriptionType *introspectionResponseTypeRef     `json:"subscriptionType"`
		Types            []*introspectionResponseType      `json:"types"`
		Directives       []*introspectionResponseDirective `json:"directives"`
	}

	introspectionResponseType struct {
		Kind          string                             `json:"kind"`
		Name          string                             `json:"name"`
		Description   string                             `json:"description"`
		Fields        []*introspectionResponseField      `json:"fields"`
		InputFields   []*introspectionResponseInputValue `json:"inputFields"`
		Interfaces    []*introspectionResponseTypeRef    `json:"interfaces"`
		EnumValues    []*introspectionResponseEnumValues `json:"enumValues"`
		PossibleTypes []*introspectionResponseTypeRef    `json:"possibleTypes"`
	}

	introspectionResponseField struct {
		Name              string                             `json:"name"`
		Description       string                             `json:"description"`
		Args              []*introspectionResponseInputValue `json:"args"`
		TypeRef           *introspectionResponseTypeRef      `json:"type"`
		IsDeprecated      bool                               `json:"isDeprecated"`
		DeprecationReason string                             `json:"deprecationReason"`
	}

	introspectionResponseInputValue struct {
		Name         string                        `json:"name"`
		Description  string                        `json:"description"`
		TypeRef      *introspectionResponseTypeRef `json:"type"`
		DefaultValue interface{}                   `json:"defaultValue"`
	}

	introspectionResponseTypeRef struct {
		Kind   string                        `json:"kind"`
		Name   string                        `json:"name"`
		OfType *introspectionResponseTypeRef `json:"ofType"`
	}

	introspectionResponseEnumValues struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		IsDeprecated      bool   `json:"isDeprecated"`
		DeprecationReason string `json:"deprecationReason"`
	}

	introspectionResponseDirective struct {
		Name        string                             `json:"name"`
		Description string                             `json:"description"`
		Locations   []string                           `json:"locations"`
		Args        []*introspectionResponseInputValue `json:"args"`
	}
)

func (t *introspectionResponseType) GetName() string {
	return t.Name
}

func (t *introspectionResponseType) GetKind() string {
	return t.Kind
}

func (t *introspectionResponseTypeRef) GetName() string {
	return t.Name
}

func (t *introspectionResponseTypeRef) GetKind() string {
	return t.Kind
}
