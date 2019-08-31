package graphql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is the object for the GraphQL module
type Module struct {
	project   string
	auth      *auth.Module
	crud      *crud.Module
	functions *functions.Module
}

// New creates a new GraphQL module
func New(a *auth.Module, c *crud.Module, f *functions.Module) *Module {
	return &Module{auth: a, crud: c, functions: f}
}

// SetConfig sets the project configuration
func (graph *Module) SetConfig(project string) {
	graph.project = project
}

// SetConfig sets the project configuration
func (graph *Module) GetProjectID() string {
	return graph.project
}

// ExecGraphQLQuery executes the provided graphql query
func (graph *Module) ExecGraphQLQuery(req *model.GraphQLRequest, token string) (interface{}, error) {

	source := source.NewSource(&source.Source{
		Body: []byte(req.Query),
		Name: req.OperationName,
	})
	// parse the source
	doc, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		return nil, err
	}

	return graph.execGraphQLDocument(doc, token, utils.M{"vars": req.Variables})
}

func (graph *Module) execGraphQLDocument(node ast.Node, token string, store utils.M) (interface{}, error) {
	switch node.GetKind() {

	case kinds.Document:
		doc := node.(*ast.Document)
		for _, v := range doc.Definitions {
			return graph.execGraphQLDocument(v, token, store)
		}
		return nil, errors.New("No definitions provided")

	case kinds.OperationDefinition:
		op := node.(*ast.OperationDefinition)
		switch op.Operation {
		case "query":
			obj := map[string]interface{}{}
			for _, v := range op.SelectionSet.Selections {

				field := v.(*ast.Field)

				result, err := graph.execGraphQLDocument(field, token, store)
				if err != nil {
					return nil, err
				}

				obj[getFieldName(field)] = result
			}

			return obj, nil
		case "mutation":

			return graph.handleMutation(node, token, store)

		default:
			return nil, errors.New("Invalid operation: " + op.Operation)
		}

	case kinds.Field:
		field := node.(*ast.Field)

		// No directive means its a nested field

		if len(field.Directives) > 0 {
			kind := getQueryKind(field.Directives[0])
			if kind == "read" {
				result, err := graph.execReadRequest(field, token, store)
				if err != nil {
					return nil, err
				}

				return graph.processQueryResult(field, token, store, result)
			}

			if kind == "func" {
				result, err := graph.execFuncCall(field, store)
				if err != nil {
					return nil, err
				}

				return graph.processQueryResult(field, token, store, result)
			}

			return nil, errors.New("Incorrect query type")
		}

		currentValue, err := utils.LoadValue(fmt.Sprintf("%s.%s", store["coreParentKey"], field.Name.Value), store)
		if err != nil {
			return nil, err
		}
		if field.SelectionSet == nil {
			return currentValue, nil
		}

		obj := utils.M{}
		for _, sel := range field.SelectionSet.Selections {
			storeNew := shallowClone(store)
			storeNew[getFieldName(field)] = currentValue
			storeNew["coreParentKey"] = getFieldName(field)

			f := sel.(*ast.Field)

			output, err := graph.execGraphQLDocument(f, token, storeNew)
			if err != nil {
				return nil, err
			}

			obj[getFieldName(f)] = output
		}

		return obj, nil

	default:
		return nil, errors.New("Invalid node type " + node.GetKind() + ": " + string(node.GetLoc().Source.Body)[node.GetLoc().Start:node.GetLoc().End])
	}
}

func getQueryKind(directive *ast.Directive) string {
	switch directive.Name.Value {

	case "postgres", "mysql", "mongo":
		return "read"

	default:
		return "func"
	}
}
