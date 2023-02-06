package graphql

import "github.com/spacecloud-io/space-cloud/modules/graphql/types"

// GraphqlLoaderKey describes a key used by the graphql dataloader
type GraphqlLoaderKey struct {
	FieldName    string
	Query        string
	AllowedVars  map[string]struct{}
	ExportedVars map[string]*types.StoreValue
}

// String returns the string representation of the key
func (k *GraphqlLoaderKey) String() string {
	return k.FieldName
}

// Raw returns the raw query value
func (k *GraphqlLoaderKey) Raw() interface{} {
	return k.Query
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
