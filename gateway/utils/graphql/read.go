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
		result, err := graph.crud.Read(ctx, dbAlias, col, req)
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

	selectionSet := graph.extractSelectionSet(field, dbAlias, col)
	if len(selectionSet) > 0 {
		req.Options.Select = selectionSet
	}

	// Check if read op is authorised
	actions, _, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbAlias, col, token, req)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		//  batch operation cannot be performed with aggregation
		req.IsBatch = !(len(req.Aggregate) > 0)
		req.Options.HasOptions = hasOptions
		result, err := graph.crud.Read(ctx, dbAlias, col, req)
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

	id := field.Name.Value

	params, err := getFuncParams(field, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req := model.PreparedQueryRequest{Params: params}
	// Check if PreparedQuery op is authorised
	actions, _, err := graph.auth.IsPreparedQueryAuthorised(ctx, graph.project, dbAlias, id, token, &req)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		result, err := graph.crud.ExecPreparedQuery(ctx, dbAlias, id, &req)
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

	readRequest.GroupBy, err = extractGroupByClause(field.Arguments, store)
	if err != nil {
		return nil, false, err
	}

	readRequest.Aggregate, err = extractAggregate(field)
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

func (graph *Module) extractSelectionSet(field *ast.Field, dbAlias, col string) map[string]int32 {
	selectMap := map[string]int32{}
	schemaFields, _ := graph.schema.GetSchema(dbAlias, col)
	if field.SelectionSet == nil {
		return nil
	}
	for _, selection := range field.SelectionSet.Selections {
		v := selection.(*ast.Field)
		// skip aggregate field & fields with directives
		if v.Name.Value == utils.GraphQLAggregate || len(v.Directives) > 0 {
			continue
		}
		if schemaFields != nil {
			// skip linked fields
			fieldStruct, p := schemaFields[v.Name.Value]
			if p && fieldStruct.IsLinked {
				continue
			}
		}
		selectMap[v.Name.Value] = 1
	}
	return selectMap
}

func extractAggregate(v *ast.Field) (map[string][]string, error) {
	functionMap := make(map[string][]string)
	aggregateFound := false
	if v.SelectionSet == nil {
		return nil, nil
	}
	for _, selection := range v.SelectionSet.Selections {
		field := selection.(*ast.Field)
		if field.Name.Value != utils.GraphQLAggregate || field.SelectionSet == nil {
			continue
		}
		if aggregateFound {
			return nil, utils.LogError("GraphQL query cannot have multiple aggregate fields, specify all functions in single aggregate field", "graphql", "extractAggregate", nil)
		}
		aggregateFound = true
		// get function name
		for _, selection := range field.SelectionSet.Selections {
			functionField := selection.(*ast.Field)
			_, ok := functionMap[functionField.Name.Value]
			if ok {
				return nil, utils.LogError(fmt.Sprintf("Cannot repeat the same function (%s) twice. Specify all columns within single function field", functionField.Name.Value), "graphql", "extractAggregate", nil)
			}

			if functionField.Name.Value == "count" && functionField.SelectionSet == nil {
				functionMap[functionField.Name.Value] = []string{""}
				continue
			}

			if functionField.SelectionSet == nil {
				return nil, nil
			}
			colArray := make([]string, 0)
			// get column name
			for _, selection := range functionField.SelectionSet.Selections {
				columnField := selection.(*ast.Field)
				colArray = append(colArray, columnField.Name.Value)
			}
			functionMap[functionField.Name.Value] = colArray
		}
	}
	return functionMap, nil
}

func extractGroupByClause(args []*ast.Argument, store utils.M) ([]interface{}, error) {
	for _, v := range args {
		switch v.Name.Value {
		case utils.GraphQLGroupByArgument:
			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			if obj, ok := temp.([]interface{}); ok {
				return obj, nil
			}
			return nil, utils.LogError(fmt.Sprintf("GraphQL (%s) argument is of type %v, but it should be of type array ([])", utils.GraphQLGroupByArgument, reflect.TypeOf(temp)), "graphql", "extractGroupByClause", nil)
		}
	}

	return make([]interface{}, 0), nil
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
