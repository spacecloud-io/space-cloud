package graphql

import "github.com/graphql-go/graphql"

type (
	// Source describes the implementation of source from the graphql module
	Source interface {
		GetGraphQLTypes() *Types
	}

	// Types describes all the types required by graphql provider
	Types struct {
		QueryTypes    graphql.Fields
		MutationTypes graphql.Fields
		AllTypes      map[string]graphql.Type
	}

	// Compiler describes the impmentation of a source which requires access to the query compiler
	Compiler interface {
		GraphqlCompiler(fn CompilerFn) error
		GetCompiledQuery() *CompiledQuery
	}

	// CompiledQueryReceiver describes the implementation of source which requires compiledQueries
	CompiledQueryReceiver interface {
		SetCompiledQueries(map[string]*CompiledQuery)
	}

	CompilerFn func(query, operationName string, defaultValues map[string]any, enableExtraction bool) (*CompiledQuery, error)
)

type (
	m map[string]interface{}
	t []interface{}
)

type request struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}
