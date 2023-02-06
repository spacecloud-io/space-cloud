package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/spacecloud-io/space-cloud/utils"
)

type (
	m map[string]interface{}
	t []interface{}
)

var anyType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Any",
	Description: "Type any for unknown types",
	Serialize: func(value interface{}) interface{} {
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		return value
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		val, _ := utils.ParseGraphqlValue(valueAST, map[string]interface{}{})
		return val
	},
})
