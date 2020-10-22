package graphql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (graph *Module) execLinkedReadRequest(ctx context.Context, field *ast.Field, dbAlias, col, token string, req *model.ReadRequest, store utils.M, cb dbCallback) {
	dbType, _ := graph.crud.GetDBType(dbAlias)
	// Check if read op is authorised
	returnWhere := model.ReturnWhereStub{Col: col, PrefixColName: false, ReturnWhere: dbType != string(model.Mongo), Where: map[string]interface{}{}}
	actions, reqParams, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbAlias, col, token, req, returnWhere)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req.PostProcess[col] = actions

	if len(returnWhere.Where) > 0 {
		for k, v := range returnWhere.Where {
			req.Find[k] = v
		}
	}

	req.GroupBy, err = extractGroupByClause(ctx, field.Arguments, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	req.Aggregate, err = extractAggregate(ctx, field, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {

		req.IsBatch = !(len(req.Aggregate) > 0)
		if req.Options == nil {
			req.Options = &model.ReadOptions{}
		}
		req.Options.HasOptions = false
		result, err := graph.crud.Read(ctx, dbAlias, col, req, reqParams)

		cb(dbAlias, col, result, err)
	}()
}

func (graph *Module) execReadRequest(ctx context.Context, field *ast.Field, token string, store utils.M, cb dbCallback) {
	dbAlias, err := graph.GetDBAlias(ctx, field, token, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	col, err := getCollection(field)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req, hasOptions, err := generateReadRequest(ctx, field, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	selectionSet, err := graph.extractSelectionSet(field, dbAlias, col, req.Options.Join, len(req.Options.Join) > 0, req.Options.ReturnType)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	if len(selectionSet) > 0 {
		req.Options.Select = selectionSet
	}

	// Check if read op is authorised
	dbType, _ := graph.crud.GetDBType(dbAlias)

	returnWhere := model.ReturnWhereStub{Col: col, PrefixColName: len(req.Options.Join) > 0, ReturnWhere: dbType != string(model.Mongo), Where: map[string]interface{}{}}
	actions, reqParams, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbAlias, col, token, req, returnWhere)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	for k, v := range returnWhere.Where {
		req.Find[k] = v
	}
	req.PostProcess[col] = actions

	if err := graph.runAuthForJoins(ctx, dbType, dbAlias, token, req, req.Options.Join); err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		//  batch operation cannot be performed with aggregation or joins or when post processing is applied
		req.IsBatch = !(len(req.Aggregate) > 0 || len(req.Options.Join) > 0 || isPostProcessingEnabled(req.PostProcess))
		req.Options.HasOptions = hasOptions
		result, err := graph.crud.Read(ctx, dbAlias, col, req, reqParams)
		cb(dbAlias, col, result, err)
	}()
}

func (graph *Module) runAuthForJoins(ctx context.Context, dbType, dbAlias, token string, req *model.ReadRequest, join []model.JoinOption) error {
	for _, j := range join {
		returnWhere := model.ReturnWhereStub{Col: j.Table, PrefixColName: len(req.Options.Join) > 0, ReturnWhere: dbType != string(model.Mongo), Where: map[string]interface{}{}}
		actions, _, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbAlias, j.Table, token, req, returnWhere)
		if err != nil {
			return err
		}
		for k, v := range returnWhere.Where {
			req.Find[k] = v
		}
		req.PostProcess[j.Table] = actions

		if err := graph.runAuthForJoins(ctx, dbType, dbAlias, token, req, j.Join); err != nil {
			return err
		}
	}

	return nil
}

func (graph *Module) execPreparedQueryRequest(ctx context.Context, field *ast.Field, token string, store utils.M, cb dbCallback) {
	dbAlias, err := graph.GetDBAlias(ctx, field, token, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	id := field.Name.Value

	params, err := getFuncParams(ctx, field, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req := model.PreparedQueryRequest{Params: params}
	// Check if PreparedQuery op is authorised
	actions, reqParams, err := graph.auth.IsPreparedQueryAuthorised(ctx, graph.project, dbAlias, id, token, &req)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		result, err := graph.crud.ExecPreparedQuery(ctx, dbAlias, id, &req, reqParams)
		_ = graph.auth.PostProcessMethod(ctx, actions, result)
		cb(dbAlias, id, result, err)
	}()
}

func generateReadRequest(ctx context.Context, field *ast.Field, store utils.M) (*model.ReadRequest, bool, error) {
	var err error

	// Create a read request object
	readRequest := model.ReadRequest{Operation: utils.All, Options: new(model.ReadOptions), PostProcess: map[string]*model.PostProcess{}}

	readRequest.Find, err = ExtractWhereClause(ctx, field.Arguments, store)
	if err != nil {
		return nil, false, err
	}

	readRequest.GroupBy, err = extractGroupByClause(ctx, field.Arguments, store)
	if err != nil {
		return nil, false, err
	}

	readRequest.Aggregate, err = extractAggregate(ctx, field, store)
	if err != nil {
		return nil, false, err
	}

	var hasOptions bool
	readRequest.Options, hasOptions, err = generateOptions(ctx, field.Arguments, store)
	if err != nil {
		return nil, false, err
	}
	// if distinct option has been set then set operation to distinct from all
	if hasOptions && readRequest.Options.Distinct != nil {
		readRequest.Operation = utils.Distinct
	}

	// Get extra arguments
	readRequest.Extras = generateArguments(ctx, field, store)

	return &readRequest, hasOptions, nil
}

func generateArguments(ctx context.Context, field *ast.Field, store utils.M) map[string]interface{} {
	obj := map[string]interface{}{}
	for _, arg := range field.Arguments {
		switch arg.Name.Value {
		case "where", "group", "skip", "limit", "sort", "distinct": // read & delete
			continue
		case "op", "set", "inc", "mul", "max", "min", "currentTimestamp", "currentDate", "push", "rename", "unset": // update
			continue
		case "docs": // create
			continue
		default:
			val, err := utils.ParseGraphqlValue(arg.Value, store)
			if err != nil {
				helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to extract argument from graphql query (%s)", arg.Name.Value), map[string]interface{}{"argValue": arg.Value.GetValue()})
				continue
			}

			obj[arg.Name.Value] = val
		}
	}
	return obj
}

func (graph *Module) extractSelectionSet(field *ast.Field, dbAlias, col string, join []model.JoinOption, isJoin bool, returnType string) (map[string]int32, error) {
	selectMap := map[string]int32{}
	schemaFields, _ := graph.schema.GetSchema(dbAlias, col)
	if field.SelectionSet == nil {
		return nil, nil
	}
	for _, selection := range field.SelectionSet.Selections {
		v := selection.(*ast.Field)
		// skip aggregate field & fields with directives
		if v.Name.Value == utils.GraphQLAggregate || len(v.Directives) > 0 {
			continue
		}

		joinTable, isJointTable := isJointTable(v.Name.Value, join)

		if schemaFields != nil {
			// skip linked fields but allow joint tables
			fieldStruct, p := schemaFields[v.Name.Value]
			if p && fieldStruct.IsLinked && !isJointTable {
				continue
			}
		}

		if isJointTable {
			if v.SelectionSet == nil {
				return nil, errors.New("joint tables cannot have an empty selection set")
			}

			m, err := graph.extractSelectionSet(v, dbAlias, joinTable.Table, joinTable.Join, isJoin, returnType)
			if err != nil {
				return nil, err
			}

			for k, v := range m {
				selectMap[k] = v
			}
			continue
		}

		// Need to make thing
		key := v.Name.Value
		if isJoin {
			if returnType == "table" {
				arr := strings.Split(key, "__")
				if len(arr) != 2 {
					return nil, fmt.Errorf("field name must be of the format `table__column`")
				}
				key = arr[0] + "." + arr[1]
			} else {
				key = col + "." + key
			}
		}
		selectMap[key] = 1
	}
	return selectMap, nil
}

func extractAggregate(ctx context.Context, v *ast.Field, store utils.M) (map[string][]string, error) {
	functionMap := make(map[string][]string)
	aggregateFound := false
	if v.SelectionSet == nil {
		return nil, nil
	}
	for _, selection := range v.SelectionSet.Selections {
		field := selection.(*ast.Field)

		// Check if aggregate was found in the directive
		if len(field.Directives) > 0 && field.Directives[0].Name.Value == "aggregate" {
			// Get the required parameters
			op, err := getAggregateArguments(field.Directives[0], store)
			if err != nil {
				return nil, err
			}

			colArray, ok := functionMap[op]
			if !ok {
				colArray = []string{}
			}
			colArray = append(colArray, strings.Join(strings.Split(field.Name.Value, "__"), ".")+":table")
			functionMap[op] = colArray
			continue
		}

		if field.Name.Value != utils.GraphQLAggregate || field.SelectionSet == nil {
			continue
		}

		if aggregateFound {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "GraphQL query cannot have multiple aggregate fields, specify all functions in single aggregate field", nil, nil)
		}
		aggregateFound = true
		// get function name
		for _, selection := range field.SelectionSet.Selections {
			functionField := selection.(*ast.Field)
			_, ok := functionMap[functionField.Name.Value]
			if ok {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Cannot repeat the same function (%s) twice. Specify all columns within single function field", functionField.Name.Value), nil, nil)
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
				colArray = append(colArray, strings.Join(strings.Split(columnField.Name.Value, "__"), "."))
			}
			functionMap[functionField.Name.Value] = colArray
		}
	}
	return functionMap, nil
}

func getAggregateArguments(field *ast.Directive, store utils.M) (string, error) {
	for _, arg := range field.Arguments {
		switch arg.Name.Value {
		case "op":
			temp, err := utils.ParseGraphqlValue(arg.Value, store)
			if err != nil {
				return "", err
			}

			op, ok := temp.(string)
			if !ok {
				return "", fmt.Errorf("invalid type provided (%s) for aggregate op", reflect.TypeOf(temp))
			}

			return op, nil
		}
	}

	return "", errors.New("need to provide `op` when using aggregations")
}

func extractGroupByClause(ctx context.Context, args []*ast.Argument, store utils.M) ([]interface{}, error) {
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
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("GraphQL (%s) argument is of type %v, but it should be of type array ([])", utils.GraphQLGroupByArgument, reflect.TypeOf(temp)), nil, nil)
		}
	}

	return make([]interface{}, 0), nil
}

// ExtractWhereClause return the where arg of graphql schema
func ExtractWhereClause(ctx context.Context, args []*ast.Argument, store utils.M) (map[string]interface{}, error) {
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

func generateOptions(ctx context.Context, args []*ast.Argument, store utils.M) (*model.ReadOptions, bool, error) {
	hasOptions := false // Flag to see if options exist
	options := model.ReadOptions{}
	for _, v := range args {
		switch v.Name.Value {
		case "join":
			hasOptions = true

			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			join := make([]model.JoinOption, 0)
			if err := mapstructure.Decode(temp, &join); err != nil {
				return nil, hasOptions, err
			}
			options.Join = join
		case "returnType":
			// We won't set hasOptions to true for this one
			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			retType, ok := temp.(string)
			if !ok {
				return nil, hasOptions, fmt.Errorf("invalid type provided for returnType; expecting string got (%s)", reflect.TypeOf(temp))
			}
			options.ReturnType = retType
		case "skip":
			hasOptions = true // Set the flag to true

			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			switch t := temp.(type) {
			case float64:
				// This condition occurs if we provide value of limit operator from graphql variables
				tempInt64 := int64(t)
				options.Skip = &tempInt64
			case int:
				tempInt64 := int64(t)
				options.Skip = &tempInt64
			default:
				return nil, hasOptions, fmt.Errorf("invalid type provided for skip expecting integer got (%s)", reflect.TypeOf(temp))
			}

		case "limit":
			hasOptions = true // Set the flag to true

			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			switch t := temp.(type) {
			case float64:
				// This condition occurs if we provide value of limit operator from graphql variables
				tempInt64 := int64(t)
				options.Limit = &tempInt64
			case int:
				tempInt64 := int64(t)
				options.Limit = &tempInt64
			default:
				return nil, hasOptions, fmt.Errorf("invalid type provided for limit expecting integer got (%s)", reflect.TypeOf(temp))
			}

		case "sort":
			hasOptions = true // Set the flag to true

			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, hasOptions, err
			}

			tempInt, ok := temp.([]interface{})
			if !ok {
				return nil, hasOptions, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for (sort) expecting array got (%v)", reflect.TypeOf(temp)), nil, nil)
			}

			sortArray := make([]string, len(tempInt))
			for i, value := range tempInt {
				valueString, ok := value.(string)
				if !ok {
					return nil, hasOptions, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for (sort) each value in array should be string got (%v)", reflect.TypeOf(value)), nil, nil)
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

func isJointTable(table string, join []model.JoinOption) (model.JoinOption, bool) {
	for _, j := range join {
		if j.Table == table {
			return j, true
		}
	}

	return model.JoinOption{}, false
}
