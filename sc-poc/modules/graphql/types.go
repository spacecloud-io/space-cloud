package graphql

import "github.com/graphql-go/graphql"

type (
	// Source describes the implementation of source from the graphql module
	Source interface {
		GetTypes() (queryTypes, mutationTypes graphql.Fields)
		GetAllTypes() map[string]graphql.Type
	}

	// Compiler describes the impmentation of a source which requires access to the query compiler
	Compiler interface {
		GraphqlCompiler(fn CompilerFn) error
		GetCompiledQuery() *CompiledQuery
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
