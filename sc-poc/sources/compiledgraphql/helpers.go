package compiledgraphql

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/spacecloud-io/space-cloud/modules/graphql"
)

func (s *Source) compile(fn graphql.CompilerFn) error {
	// Get the compiled graphql query
	var defaultVariables map[string]interface{}
	if len(s.Spec.Graphql.DefaultVariables.Raw) > 0 {
		if err := json.Unmarshal(s.Spec.Graphql.DefaultVariables.Raw, &defaultVariables); err != nil {
			return err
		}
	}
	compiledQuery, err := fn(s.Spec.Graphql.Query, s.Spec.Graphql.OperationName, defaultVariables, true)
	if err != nil {
		return err
	}

	s.compiledQuery = compiledQuery
	return nil
}

func (s *Source) getSchemas() (requestSchema, responseSchema *openapi3.SchemaRef) {
	requestSchema = new(openapi3.SchemaRef)
	data, _ := s.compiledQuery.VariableSchema.MarshalJSON()
	_ = requestSchema.UnmarshalJSON(data)

	responseSchema = new(openapi3.SchemaRef)
	data, _ = s.compiledQuery.ResponseSchema.MarshalJSON()
	_ = responseSchema.UnmarshalJSON(data)
	return
}

func (s *Source) call(ctx context.Context, vars map[string]any) (any, error) {
	result := s.compiledQuery.Execute(ctx, vars)
	if result.HasErrors() {
		return nil, errors.New("unable to execute graphql request")
	}

	return result.Data, nil
}
