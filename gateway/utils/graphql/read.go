package graphql

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (graph *Module) execLinkedReadRequest(ctx context.Context, field *ast.Field, dbAlias, col, token string, req *model.ReadRequest, store utils.M, cb dbCallback) {
	// Check if read op is authorised
	actions, _, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbAlias, col, token, req)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		req.IsBatch = true
		if req.Options == nil {
			req.Options = &model.ReadOptions{}
		}
		req.Options.HasOptions = false
		result, err := graph.crud.Read(ctx, dbAlias, graph.project, col, req)
		_ = graph.auth.PostProcessMethod(actions, result)

		cb(dbAlias, col, result, err)
	}()
}

func (graph *Module) execReadRequest(ctx context.Context, field *ast.Field, token string, store utils.M, cb dbCallback) {
	dbAlias, err := graph.GetDBAlias(field)
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
	actions, _, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbAlias, col, token, req)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		req.IsBatch = true
		req.Options.HasOptions = hasOptions
		result, err := graph.crud.Read(ctx, dbAlias, graph.project, col, req)
		_ = graph.auth.PostProcessMethod(actions, result)
		cb(dbAlias, col, result, err)
	}()
}

func (graph *Module) execPreparedQueryRequest(ctx context.Context, field *ast.Field, token string, store utils.M, cb dbCallback) {
	dbAlias, err := graph.GetDBAlias(field)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	id, err := getCollection(field)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	param, err := getFuncParams(field, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req := model.PreparedQueryRequest{Params: param}
	// Check if PreparedQuery op is authorised
	actions, _, err := graph.auth.IsPreparedQueryAuthorised(ctx, graph.project, dbAlias, id, token, &req)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		result, err := graph.crud.ExecPreparedQuery(ctx, graph.project, dbAlias, id, &req)
		_ = graph.auth.PostProcessMethod(actions, result)
		cb(dbAlias, id, result, err)
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
	// if distinct option has been set then set operation to distinct from all
	if hasOptions && readRequest.Options.Distinct != nil {
		readRequest.Operation = utils.Distinct
	}

	return &readRequest, hasOptions, nil
}

// ExtractWhereClause return the where arg of graphql schema
func ExtractWhereClause(args []*ast.Argument, store utils.M) (map[string]interface{}, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "where":
			temp, err := utils.ParseGraphqlValue(v.Value, store)
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

			temp, err := utils.ParseGraphqlValue(v.Value, store)
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

			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			tempInt, ok := temp.(int)
			if !ok {
				return nil, hasOptions, errors.New("Invalid type for skip")
			}

			tempInt64 := int64(tempInt)
			options.Limit = &tempInt64

		case "sort":
			hasOptions = true // Set the flag to true

			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			tempInt, ok := temp.([]interface{})
			if !ok {
				return nil, hasOptions, fmt.Errorf("Invalid type (%s) for sort", reflect.TypeOf(temp))
			}

			sortArray := make([]string, len(tempInt))
			for i, value := range tempInt {
				valueString, ok := value.(string)
				if !ok {
					return nil, hasOptions, fmt.Errorf("Invalid type (%s) for sort", reflect.TypeOf(value))
				}
				sortArray[i] = valueString
			}

			options.Sort = sortArray

		case "distinct":
			hasOptions = true // Set the flag to true

			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			tempString, ok := temp.(string)
			if !ok {
				return nil, hasOptions, errors.New("Invalid type for distinct")
			}

			options.Distinct = &tempString
		}
	}
	return &options, hasOptions, nil
}
