package graphql

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Module is the object for the GraphQL module
type Module struct {
	project   string
	auth      AuthInterface
	crud      CrudInterface
	functions FunctionInterface
	schema    SchemaInterface
}

// New creates a new GraphQL module
func New(a AuthInterface, c CrudInterface, f FunctionInterface, s SchemaInterface) *Module {
	return &Module{auth: a, crud: c, functions: f, schema: s}
}

// SetConfig sets the project configuration
func (graph *Module) SetConfig(project string) {
	graph.project = project
}

// GetProjectID sets the project configuration
func (graph *Module) GetProjectID() string {
	return graph.project
}

// ExecGraphQLQuery executes the provided graphql query
func (graph *Module) ExecGraphQLQuery(ctx context.Context, req *model.GraphQLRequest, token string, cb model.GraphQLCallback) {

	s := source.NewSource(&source.Source{
		Body: []byte(req.Query),
		Name: req.OperationName,
	})
	// parse the source
	doc, err := parser.Parse(parser.ParseParams{Source: s})
	if err != nil {
		cb(nil, err)
		return
	}

	graph.execGraphQLDocument(ctx, doc, token, utils.M{"vars": req.Variables, "path": "", "directive": ""}, nil, createCallback(cb))
}

type dbCallback func(dbAlias, col string, op interface{}, err error)

func createCallback(cb model.GraphQLCallback) model.GraphQLCallback {
	var lock sync.Mutex
	var isCalled bool

	return func(result interface{}, err error) {
		lock.Lock()
		defer lock.Unlock()

		// Check if callback has already been invoked once
		if isCalled {
			return
		}

		// Set the flag to prevent duplicate invocation
		isCalled = true
		cb(result, err)
	}
}

func createDBCallback(cb dbCallback) dbCallback {
	var lock sync.Mutex
	var isCalled bool

	return func(dbAlias, col string, result interface{}, err error) {
		lock.Lock()
		defer lock.Unlock()

		// Check if callback has already been invoked once
		if isCalled {
			return
		}

		// Set the flag to prevent duplicate invocation
		isCalled = true
		cb(dbAlias, col, result, err)
	}
}

func (graph *Module) execGraphQLDocument(ctx context.Context, node ast.Node, token string, store utils.M, schema model.Fields, cb model.GraphQLCallback) {
	switch node.GetKind() {

	case kinds.Document:
		doc := node.(*ast.Document)
		if len(doc.Definitions) > 0 {
			graph.execGraphQLDocument(ctx, doc.Definitions[0], token, store, nil, createCallback(cb))
			return
		}
		cb(nil, errors.New("No definitions provided"))
		return

	case kinds.OperationDefinition:
		op := node.(*ast.OperationDefinition)
		switch op.Operation {
		case "query":
			obj := utils.NewObject()

			// Create a wait group
			var wg sync.WaitGroup
			wg.Add(len(op.SelectionSet.Selections))

			for _, v := range op.SelectionSet.Selections {

				field := v.(*ast.Field)

				graph.execGraphQLDocument(ctx, field, token, store, nil, createCallback(func(result interface{}, err error) {
					defer wg.Done()
					if err != nil {
						cb(nil, err)
						return
					}

					// Set the result in the field
					obj.Set(getFieldName(field), result)
				}))
			}

			// Wait then return the result
			wg.Wait()
			cb(obj.GetAll(), nil)
			return
		case "mutation":
			graph.handleMutation(ctx, node, token, store, cb)
			return

		default:
			cb(nil, errors.New("Invalid operation: "+op.Operation))
			return
		}

	case kinds.Field:
		field := node.(*ast.Field)

		// No directive means its a nested field
		if len(field.Directives) > 0 && field.Directives[0].Name.Value != "aggregate" {
			directive, err := graph.getDirectiveName(ctx, field.Directives[0], token, store)
			if err != nil {
				cb(nil, err)
				return
			}

			kind := graph.getQueryKind(directive, field.Name.Value)
			if kind == "read" {
				graph.execReadRequest(ctx, field, token, store, createDBCallback(func(dbAlias, col string, result interface{}, err error) {
					if err != nil {
						cb(nil, err)
						return
					}

					// Load the schema
					s, _ := graph.schema.GetSchema(dbAlias, col)

					graph.processQueryResult(ctx, field, token, store, result, s, cb)
				}))
				return
			}

			if kind == "prepared-queries" {
				graph.execPreparedQueryRequest(ctx, field, token, store, createDBCallback(func(dbAlias, col string, result interface{}, err error) {
					if err != nil {
						cb(nil, err)
						return
					}

					graph.processQueryResult(ctx, field, token, store, result, nil, cb)
				}))
				return
			}

			if kind == "func" {
				graph.execFuncCall(ctx, token, field, store, createCallback(func(result interface{}, err error) {
					if err != nil {
						cb(nil, err)
						return
					}

					graph.processQueryResult(ctx, field, token, store, result, nil, cb)
				}))
				return
			}

			cb(nil, errors.New("incorrect query type"))
			return
		}

		currentValue, err := utils.LoadValue(fmt.Sprintf("%s.%s", store["coreParentKey"], field.Name.Value), store)
		if err != nil {
			if schema != nil {
				fieldStruct, p := schema[field.Name.Value]
				if p && fieldStruct.IsLinked {
					linkedInfo := fieldStruct.LinkedTable
					loadKey := fmt.Sprintf("%s.%s", store["coreParentKey"], linkedInfo.From)
					val, err := utils.LoadValue(loadKey, store)
					if err != nil {
						cb(nil, nil)
						return
					}
					req := &model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{linkedInfo.To: val}, PostProcess: map[string]*model.PostProcess{}}
					graph.processLinkedResult(ctx, field, *fieldStruct, token, req, store, cb)
					return
				}
			}

			// if the field isn't found in the store means that field did not exist in the result. so return nil as error
			cb(nil, nil)
			return
		}
		if field.SelectionSet == nil {
			cb(currentValue, nil)
			return
		}

		graph.processQueryResult(ctx, field, token, store, currentValue, schema, cb)
		return

	default:
		cb(nil, errors.New("Invalid node type "+node.GetKind()+": "+string(node.GetLoc().Source.Body)[node.GetLoc().Start:node.GetLoc().End]))
		return
	}
}

func (graph *Module) getQueryKind(directive, fieldName string) string {
	_, err := graph.crud.GetDBType(directive)
	if err == nil {
		if graph.crud.IsPreparedQueryPresent(directive, fieldName) {
			return "prepared-queries"
		}
		return "read"
	}
	return "func"
}
