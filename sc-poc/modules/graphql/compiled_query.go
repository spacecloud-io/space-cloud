package graphql

import (
	"context"
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/invopop/jsonschema"

	"github.com/spacecloud-io/space-cloud/modules/graphql/rootvalue"
	"github.com/spacecloud-io/space-cloud/utils"
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

	// Internal types
	schema graphql.Schema
}

// GetCompiledQuery returns a compiled query
func (a *App) GetCompiledQuery(name string) (q *CompiledQuery, p bool) {
	q, p = a.compiledQueries[name]
	return
}

// Compile compiles the graphql request for processing
func (a *App) Compile(query, operationName string, defaultValues map[string]any, enableExtraction bool) (*CompiledQuery, error) {
	source := source.NewSource(&source.Source{
		Body: []byte(query),
		Name: "GraphQL request",
	})
	graphqlDoc, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		return nil, err
	}

	operationAST := graphqlDoc.Definitions[0].(*ast.OperationDefinition)

	// Set the default values to an empty map if it isn't provided
	if defaultValues == nil {
		defaultValues = make(map[string]interface{})
	}

	isAuthRequired, injectedClaims, exportedVars := preprocessForAuth(graphqlDoc)

	var variableSchema, responseSchema *jsonschema.Schema
	var extensions map[string]interface{}
	if enableExtraction {

		// Prepare map of fields to be ignored
		ignoredFields := make(map[string]struct{}, len(injectedClaims)+len(exportedVars))
		for _, variable := range injectedClaims {
			ignoredFields[variable] = struct{}{}
		}
		for variable, v := range exportedVars {
			ignoredFields[variable] = v
		}

		variableSchema, err = a.convertVariablesToJSONSchema(operationAST, ignoredFields)
		if err != nil {
			return nil, err
		}

		responseSchema, extensions, err = a.convertGraphqlOutputToJSONSchema(operationAST)
		if err != nil {
			return nil, err
		}
	}

	return &CompiledQuery{
		IsAuthRequired: isAuthRequired,
		InjectedClaims: injectedClaims,

		VariableSchema: variableSchema,
		ResponseSchema: responseSchema,
		Extensions:     extensions,

		Doc:           graphqlDoc,
		Query:         query,
		OperationName: operationName,
		OperationType: operationAST.Operation,
		DefaultValues: defaultValues,

		schema: a.schema,
	}, nil
}

func (compiledQuery *CompiledQuery) AuthenticateRequest(ctx context.Context, vars map[string]interface{}) error {
	if compiledQuery.IsAuthRequired {
		authResult, p := utils.GetAuthenticationResult(ctx)
		if !p || !authResult.IsAuthenticated {
			// Send an error if request is unauthenticated
			return errors.New("unable to authenticate request")
		}

		// Inject the claims in the variables
		for claim, variable := range compiledQuery.InjectedClaims {
			v, p := authResult.Claims[claim]
			if !p {
				// We need to throw an error if requested claim is not present in token
				return fmt.Errorf("token does not contain required claim '%s'", claim)
			}

			_ = utils.StoreValue(variable, v, vars)
		}
	}
	return nil
}

// Execute executes the compiled graphql query
func (compiledQuery *CompiledQuery) Execute(ctx context.Context, vars map[string]interface{}) *graphql.Result {

	newVars := utils.MergeMaps(compiledQuery.DefaultValues, vars)

	// Execute the query
	rootValue := rootvalue.New(compiledQuery.Doc)
	result := graphql.Execute(graphql.ExecuteParams{
		Context:       ctx,
		Schema:        compiledQuery.schema,
		AST:           compiledQuery.Doc,
		OperationName: compiledQuery.OperationName,
		Args:          newVars,
		Root:          rootValue,
	})

	if rootValue.HasErrors() {
		result.Errors = append(result.Errors, rootValue.GetFormatedErrors()...)
	}

	return result
}
