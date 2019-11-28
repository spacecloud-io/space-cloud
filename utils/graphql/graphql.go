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

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/modules/schema"
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

// GetProjectID sets the project configuration
func (graph *Module) GetProjectID() string {
	return graph.project
}

// ExecGraphQLQuery executes the provided graphql query
func (graph *Module) ExecGraphQLQuery(ctx context.Context, req *model.GraphQLRequest, token string, cb callback) {

	source := source.NewSource(&source.Source{
		Body: []byte(req.Query),
		Name: req.OperationName,
	})
	// parse the source
	doc, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		cb(nil, err)
		return
	}

	graph.execGraphQLDocument(ctx, doc, token, utils.M{"vars": req.Variables, "path": ""}, newLoaderMap(), nil, createCallback(cb))
}

type callback func(op interface{}, err error)
type dbCallback func(dbType, col string, op interface{}, err error)

func createCallback(cb callback) callback {
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

	return func(dbType, col string, result interface{}, err error) {
		lock.Lock()
		defer lock.Unlock()

		// Check if callback has already been invoked once
		if isCalled {
			return
		}

		// Set the flag to prevent duplicate invocation
		isCalled = true
		cb(dbType, col, result, err)
	}
}

func (graph *Module) execGraphQLDocument(ctx context.Context, node ast.Node, token string, store utils.M, loader *loaderMap, schema schema.SchemaFields, cb callback) {
	switch node.GetKind() {

	case kinds.Document:
		doc := node.(*ast.Document)
		if len(doc.Definitions) > 0 {
			graph.execGraphQLDocument(ctx, doc.Definitions[0], token, store, loader, nil, createCallback(cb))
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

				graph.execGraphQLDocument(ctx, field, token, store, loader, nil, createCallback(func(result interface{}, err error) {
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
		if len(field.Directives) > 0 {
			kind := getQueryKind(field.Directives[0])
			if kind == "read" {
				graph.execReadRequest(ctx, field, token, store, loader, createDBCallback(func(dbType, col string, result interface{}, err error) {
					if err != nil {
						cb(nil, err)
						return
					}

					// Load the schema
					s, _ := graph.auth.Schema.GetSchema(dbType, col)

					graph.processQueryResult(ctx, field, token, store, result, loader, s, cb)
				}))
				return
			}

			if kind == "func" {
				graph.execFuncCall(ctx, token, field, store, createCallback(func(result interface{}, err error) {
					if err != nil {
						cb(nil, err)
						return
					}

					graph.processQueryResult(ctx, field, token, store, result, loader, nil, cb)
				}))
				return
			}

			cb(nil, errors.New("incorrect query type"))
			return
		}

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
				req := &model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{linkedInfo.To: val}}
				graph.processLinkedResult(ctx, field, fieldStruct, token, req, store, loader, cb)
				return
			}
		}

		currentValue, err := utils.LoadValue(fmt.Sprintf("%s.%s", store["coreParentKey"], field.Name.Value), store)
		if err != nil {
			cb(nil, nil)
			return
		}
		if field.SelectionSet == nil {
			cb(currentValue, nil)
			return
		}

		obj := utils.NewObject()

		// Create a wait group
		var wg sync.WaitGroup
		wg.Add(len(field.SelectionSet.Selections))

		for _, sel := range field.SelectionSet.Selections {
			storeNew := shallowClone(store)
			storeNew[getFieldName(field)] = currentValue
			storeNew["coreParentKey"] = getFieldName(field)

			f := sel.(*ast.Field)

			graph.execGraphQLDocument(ctx, f, token, storeNew, loader, schema, createCallback(func(object interface{}, err error) {
				defer wg.Done()

				if err != nil {
					cb(nil, err)
					return
				}

				obj.Set(getFieldName(f), object)
			}))
		}

		// Wait then return the result
		wg.Wait()
		cb(obj.GetAll(), nil)
		return

	default:
		cb(nil, errors.New("Invalid node type "+node.GetKind()+": "+string(node.GetLoc().Source.Body)[node.GetLoc().Start:node.GetLoc().End]))
		return
	}
}

func getQueryKind(directive *ast.Directive) string {
	switch directive.Name.Value {

	case "postgres", "mysql", "mongo", "sqlserver":
		return "read"

	default:
		return "func"
	}
}
