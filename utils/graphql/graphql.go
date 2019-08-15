package graphql

import (
	"log"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/spaceuptech/space-cloud/model"
)

type Operations map[string]Operation

// Operation is the main graphql module
type Operation struct {
	kind     string // can be read or func request
	readReq  *model.ReadRequest
	funcReq  *model.FunctionsRequest
	children Operations
}

func parseGraphQLQuery(query string) error {
	source := source.NewSource(&source.Source{
		Body: []byte(query),
		Name: "GraphQL request",
	})

	// parse the source
	doc, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		return err
	}

	ops := Operations{}
	mapGraphQLDocument(doc, ops)
	return nil
}

func mapGraphQLDocument(node ast.Node, ops Operation) {
	log.Println(node.GetKind())
	switch node.GetKind() {

	case kinds.Document:
		doc := node.(*ast.Document)
		for _, v := range doc.Definitions {
			mapGraphQLDocument(v, ops)
		}

	case kinds.OperationDefinition:
		op := node.(*ast.OperationDefinition)
		switch op.Operation {
		case "query":
			ops[op.Operation] = Operation{}
			for _, v := range op.SelectionSet.Selections {
				field := v.(*ast.Field)
				mapGraphQLDocument(field, ops)
			}
		}

	case kinds.Field:
		field := node.(*ast.Field)

		// No directive means its a nested field
		if len(field.Directives) == 0 {
			return
		}

		kind := getQueryKind(field.Directives[0])
		op := Operation{kind: kind}

		if kind == "read" {
			op.readReq = getReadRequest(field.Arguments)
		}
		ops = op
	}
}

func getReadRequest(args []*ast.Arguments) *model.ReadRequest {
	find := map[string]interface{}{}
	for _, v := range args {

	}
}

func getQueryKind(directive *ast.Directive) string {
	switch directive.Name.Value {
	case "mongo", "postgres", "mysql":
		return "read"

	default:
		return "func"
	}
}
