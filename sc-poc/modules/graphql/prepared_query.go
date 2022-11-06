package graphql

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/invopop/jsonschema"

	"github.com/spacecloud-io/space-cloud/modules/auth"
	"github.com/spacecloud-io/space-cloud/modules/graphql/rootvalue"
	"github.com/spacecloud-io/space-cloud/utils"
)

// Compile compiles the graphql request for processing
func (a *App) Compile(query, operationName string, defaultValues map[string]interface{}, enableExtraction bool) (*CompiledQuery, error) {
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

	var variableSchema, responseSchema *jsonschema.Schema
	var extensions map[string]interface{}
	if enableExtraction {
		// TODO: This should not contain exported or injected fields
		variableSchema, err = a.convertVariablesToJSONSchema(operationAST)
		if err != nil {
			return nil, err
		}

		responseSchema, extensions, err = a.convertGraphqlOutputToJSONSchema(operationAST)
		if err != nil {
			return nil, err
		}
	}

	isAuthRequired, injectedClaims := preprocessForAuth(graphqlDoc)
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
	}, nil
}

func AuthenticateRequest(r *http.Request, compiledQuery *CompiledQuery, vars map[string]interface{}) error {
	if compiledQuery.IsAuthRequired {
		authResult, p := auth.GetAuthenticationResult(r)
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
func (a *App) Execute(ctx context.Context, compiledQuery *CompiledQuery, vars map[string]interface{}) *graphql.Result {

	newVars := utils.MergeMaps(compiledQuery.DefaultValues, vars)

	// Execute the query
	rootValue := rootvalue.New(compiledQuery.Doc)
	result := graphql.Execute(graphql.ExecuteParams{
		Context:       ctx,
		Schema:        a.schema,
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
