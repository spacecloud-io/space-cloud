package graphql

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) execLinkedReadRequest(ctx context.Context, field *ast.Field, dbType, col, token string, req *model.ReadRequest, store utils.M, loader *loaderMap, cb dbCallback) {
	// Check if read op is authorised
	if _, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbType, col, token, req); err != nil {
		cb("", "", nil, err)
		return
	}

	dataLoader := loader.get(getFieldName(field)+"."+store["path"].(string)+".linked."+col, graph)

	go func() {
		// Create dataloader key
		key := model.ReadRequestKey{DBType: dbType, Col: col, HasOptions: false, Req: *req}
		result, err := dataLoader.Load(ctx, key)()
		cb(dbType, col, result, err)
	}()
}

func (graph *Module) execReadRequest(ctx context.Context, field *ast.Field, token string, store utils.M, loader *loaderMap, cb dbCallback) {
	dbType, err := GetDBType(field)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	col, err := getCollection(field)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req, hasOptions, err := generateReadRequest(field, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	// Check if read op is authorised
	if _, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbType, col, token, req); err != nil {
		cb("", "", nil, err)
		return
	}

	dataLoader := loader.get(getFieldName(field)+"."+store["path"].(string), graph)

	go func() {
		// Create dataloader key
		key := model.ReadRequestKey{DBType: dbType, Col: col, HasOptions: hasOptions, Req: *req}
		result, err := dataLoader.Load(ctx, key)()
		cb(dbType, col, result, err)
	}()
}

func generateReadRequest(field *ast.Field, store utils.M) (*model.ReadRequest, bool, error) {
	var err error

	// Create a read request object
	readRequest := model.ReadRequest{Operation: utils.All, Options: new(model.ReadOptions)}

	readRequest.Find, err = ExtractWhereClause(field.Arguments, store)
	if err != nil {
		return nil, false, err
	}

	var hasOptions bool
	readRequest.Options, hasOptions, err = generateOptions(field.Arguments, store)
	if err != nil {
		return nil, false, err
	}

	return &readRequest, hasOptions, nil
}

// ExtractWhereClause return the where arg of graphql schema
func ExtractWhereClause(args []*ast.Argument, store utils.M) (map[string]interface{}, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "where":
			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			if obj, ok := temp.(utils.M); ok {
				return obj, nil
			}
			if obj, ok := temp.(map[string]interface{}); ok {
				return obj, nil
			}
			return nil, errors.New("Invalid where clause provided")
		}
	}

	return utils.M{}, nil
}

func generateOptions(args []*ast.Argument, store utils.M) (*model.ReadOptions, bool, error) {
	hasOptions := false // Flag to see if options exist
	options := model.ReadOptions{}
	for _, v := range args {
		switch v.Name.Value {
		case "skip":
			hasOptions = true // Set the flag to true

			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			tempInt, ok := temp.(int)
			if !ok {
				return nil, hasOptions, errors.New("Invalid type for skip")
			}

			tempInt64 := int64(tempInt)
			options.Skip = &tempInt64

		case "limit":
			hasOptions = true // Set the flag to true

			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			tempInt, ok := temp.(int)
			if !ok {
				return nil, hasOptions, errors.New("Invalid type for skip")
			}

			tempInt64 := int64(tempInt)
			options.Limit = &tempInt64

			// TODO: implement sort, distinct, etc.
		}
	}
	return &options, hasOptions, nil
}
