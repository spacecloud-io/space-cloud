package graphql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	authHelpers "github.com/spaceuptech/space-cloud/gateway/modules/auth/helpers"
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
		req.MatchWhere = append(req.MatchWhere, returnWhere.Where)
	}

	req.GroupBy, err = extractGroupByClause(ctx, field.Arguments, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	var hasOptions bool
	req.Options, hasOptions, err = generateOptions(ctx, field.Arguments, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req.Options.HasOptions = hasOptions

	req.Aggregate, err = extractAggregate(ctx, field, store, dbType, col)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		req.IsBatch = !(len(req.Aggregate) > 0)
		if req.Options == nil {
			req.Options = &model.ReadOptions{}
		}
		result, metaData, err := graph.crud.Read(ctx, dbAlias, col, req, reqParams)
		if err != nil {
			cb("", "", nil, err)
			return
		}

		if req.Options.Debug && metaData != nil {
			val := store["_query"]
			val.(*utils.Array).Append(structs.Map(metaData))
		}

		// Post process only if joins were not enabled
		if isPostProcessingEnabled(req.PostProcess) && len(req.Options.Join) == 0 {
			_ = authHelpers.PostProcessMethod(ctx, graph.aesKey, req.PostProcess[col], result)
		}

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

	functionMap, selectionSet, err := graph.extractSelectionSet(ctx, field, store, dbAlias, col, &req.Options.Join, req.Options.ReturnType)
	if err != nil {
		cb("", "", nil, err)
		return
	}
	req.Aggregate = functionMap
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
	if len(returnWhere.Where) > 0 {
		req.MatchWhere = append(req.MatchWhere, returnWhere.Where)
	}

	req.PostProcess[col] = actions

	if err := graph.runAuthForJoins(ctx, dbType, dbAlias, token, req, req.Options.Join); err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		//  batch operation cannot be performed with aggregation or joins or when post processing is applied or when cache is not nil
		req.IsBatch = !(len(req.Aggregate) > 0 || len(req.Options.Join) > 0 || req.Cache != nil)
		req.Options.HasOptions = hasOptions
		result, metaData, err := graph.crud.Read(ctx, dbAlias, col, req, reqParams)
		if err != nil {
			cb("", "", nil, err)
			return
		}

		if req.Options.Debug && metaData != nil {
			val := store["_query"]
			val.(*utils.Array).Append(structs.Map(metaData))
		}

		// Post process only if joins were not enabled
		if isPostProcessingEnabled(req.PostProcess) && len(req.Options.Join) == 0 {
			_ = authHelpers.PostProcessMethod(ctx, graph.aesKey, req.PostProcess[col], result)
		}

		cb(dbAlias, col, result, err)
	}()
}

func (graph *Module) runAuthForJoins(ctx context.Context, dbType, dbAlias, token string, req *model.ReadRequest, join []*model.JoinOption) error {
	for _, j := range join {
		returnWhere := model.ReturnWhereStub{Col: j.Table, PrefixColName: len(req.Options.Join) > 0, ReturnWhere: dbType != string(model.Mongo), Where: map[string]interface{}{}}
		actions, _, err := graph.auth.IsReadOpAuthorised(ctx, graph.project, dbAlias, j.Table, token, req, returnWhere)
		if err != nil {
			return err
		}

		if len(returnWhere.Where) > 0 {
			req.MatchWhere = append(req.MatchWhere, returnWhere.Where)
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

	isDebug, err := getDebugParam(ctx, field.Arguments, store)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	req := model.PreparedQueryRequest{Params: params, Debug: isDebug}
	// Check if PreparedQuery op is authorised
	actions, reqParams, err := graph.auth.IsPreparedQueryAuthorised(ctx, graph.project, dbAlias, id, token, &req)
	if err != nil {
		cb("", "", nil, err)
		return
	}

	go func() {
		result, metaData, err := graph.crud.ExecPreparedQuery(ctx, dbAlias, id, &req, reqParams)
		if err != nil {
			cb("", "", nil, err)
			return
		}

		if req.Debug && metaData != nil {
			val := store["_query"]
			val.(*utils.Array).Append(structs.Map(metaData))
		}
		_ = authHelpers.PostProcessMethod(ctx, graph.aesKey, actions, result)
		cb(dbAlias, id, result, err)
	}()
}

func generateReadRequest(ctx context.Context, field *ast.Field, store utils.M) (*model.ReadRequest, bool, error) {
	var err error

	op, err := extractQueryOp(ctx, field.Arguments, store)
	if err != nil {
		return nil, false, err
	}

	// Create a read request object
	readRequest := model.ReadRequest{Operation: op, Options: new(model.ReadOptions), PostProcess: map[string]*model.PostProcess{}}

	readRequest.Find, err = ExtractWhereClause(field.Arguments, store)
	if err != nil {
		return nil, false, err
	}

	readRequest.GroupBy, err = extractGroupByClause(ctx, field.Arguments, store)
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

	readRequest.Cache, err = generateCacheOptions(ctx, field.Directives, store)
	if err != nil {
		return nil, false, err
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
func (graph *Module) checkIfLinkCanBeOptimized(fieldStruct *model.FieldType, dbAlias, col string) (*model.JoinOption, bool) {
	currentTableFieldID := fieldStruct.LinkedTable.From
	referredTableFieldID := fieldStruct.LinkedTable.To
	if fieldStruct.LinkedTable.Field != "" {
		return nil, false
	}
	referredTableName := fieldStruct.LinkedTable.Table
	referredDbAlias := fieldStruct.LinkedTable.DBType
	if dbAlias != referredDbAlias { // join cannot happen over different databases
		return nil, false
	}
	dbType, err := graph.crud.GetDBType(dbAlias)
	if err != nil {
		return nil, false
	}
	linkedOp := utils.All
	if !fieldStruct.IsList {
		linkedOp = utils.One
	}
	if model.DBType(dbType) == model.Mongo {
		return nil, false
	}
	return &model.JoinOption{
		Op:    linkedOp,
		Table: referredTableName,
		As:    fieldStruct.FieldName,
		On: map[string]interface{}{
			fmt.Sprintf("%s.%s", col, currentTableFieldID): fmt.Sprintf("%s.%s", referredTableName, referredTableFieldID),
		},
		Type: "LEFT",
	}, true
}

func (graph *Module) extractSelectionSet(ctx context.Context, field *ast.Field, store utils.M, dbAlias, col string, join *[]*model.JoinOption, returnType string) (map[string][]string, map[string]int32, error) {
	selectMap := map[string]int32{}
	functionMap := make(map[string][]string)
	aggregateFound := new(bool)
	schemaFields, _ := graph.schema.GetSchema(dbAlias, col)
	if field.SelectionSet == nil {
		return nil, nil, nil
	}

	dbType, err := graph.crud.GetDBType(dbAlias)
	if err != nil {
		return nil, nil, err
	}

	for _, selection := range field.SelectionSet.Selections {
		v := selection.(*ast.Field)

		// Skip dbFetchTs fields
		if v.Name.Value == "_dbFetchTs" {
			continue
		}

		// skip aggregate field & fields with directives
		if v.Name.Value == utils.GraphQLAggregate || (len(v.Directives) > 0 && v.Directives[0].Name.Value == utils.GraphQLAggregate) {
			f, err := aggregateSingleField(ctx, v, store, col, model.DBType(dbType), aggregateFound)
			if err != nil {
				return nil, nil, err
			}
			for key, value := range f {
				v, ok := functionMap[key]
				if !ok {
					functionMap[key] = value
				} else {
					functionMap[key] = append(v, value...)
				}
			}
			continue
		}
		if len(v.Directives) > 0 {
			continue
		}

		joinTable, isJointTable := isJointTable(v.Name.Value, *join)
		if schemaFields != nil {
			// skip linked fields but allow joint tables
			fieldStruct, p := schemaFields[v.Name.Value]
			if p && fieldStruct.IsLinked && !isJointTable {
				if v.SelectionSet == nil {
					continue
				}
				// check if the link can be optimised to join
				joinInfo, isOptimized := graph.checkIfLinkCanBeOptimized(fieldStruct, dbAlias, col)
				if !isOptimized {
					continue
				}
				*join = append(*join, joinInfo)
				joinTable = joinInfo
				isJointTable = true
			}
		}

		if isJointTable {
			if v.SelectionSet == nil {
				return nil, nil, errors.New("joint tables cannot have an empty selection set")
			}

			f, m, err := graph.extractSelectionSet(ctx, v, store, dbAlias, joinTable.Table, &joinTable.Join, returnType)
			if err != nil {
				return nil, nil, err
			}
			for key, value := range f {
				v, ok := functionMap[key]
				if !ok {
					functionMap[key] = value
				} else {
					functionMap[key] = append(v, value...)
				}
			}
			for k, v := range m {
				selectMap[k] = v
			}
			continue
		}

		// Need to make thing
		key := v.Name.Value
		if model.DBType(dbType) != model.Mongo {
			if returnType == "table" {
				arr := strings.Split(key, "__")
				if len(arr) != 2 {
					return nil, nil, fmt.Errorf("field name must be of the format `table__column`")
				}
				key = arr[0] + "." + arr[1]
			} else {
				key = col + "." + key
			}
		}
		selectMap[key] = 1
	}
	return functionMap, selectMap, nil
}

func extractAggregate(ctx context.Context, v *ast.Field, store utils.M, dbType, col string) (map[string][]string, error) {
	functionMap := make(map[string][]string)
	aggregateFound := new(bool)
	if v.SelectionSet == nil {
		return nil, nil
	}
	for _, selection := range v.SelectionSet.Selections {
		f, err := aggregateSingleField(ctx, selection.(*ast.Field), store, col, model.DBType(dbType), aggregateFound)
		if err != nil {
			return nil, err
		}
		for key, value := range f {
			v, ok := functionMap[key]
			if !ok {
				functionMap[key] = value
			} else {
				functionMap[key] = append(v, value...)
			}
		}
	}
	return functionMap, nil
}

func aggregateSingleField(ctx context.Context, field *ast.Field, store utils.M, tableName string, dbType model.DBType, aggregateFound *bool) (map[string][]string, error) {
	functionMap := make(map[string][]string)

	// Check if aggregate was found in the directive
	if len(field.Directives) > 0 && field.Directives[0].Name.Value == "aggregate" {
		// Get the required parameters
		returnField := field.Name.Value
		op, fieldName, err := getAggregateArguments(field.Directives[0], store)
		if err != nil {
			return nil, err
		}

		if dbType != model.Mongo {
			if fieldName == "" {
				fieldName = tableName + "." + returnField
			}
			if arr := strings.Split(fieldName, "."); len(arr) <= 1 {
				fieldName = tableName + "." + fieldName
			}
		} else if fieldName == "" {
			fieldName = returnField
		}

		colArray, ok := functionMap[op]
		if !ok {
			colArray = []string{}
		}

		columnName := fmt.Sprintf("%s:%s:table", returnField, strings.Join(strings.Split(fieldName, "__"), "."))
		colArray = append(colArray, columnName)
		functionMap[op] = colArray
		return functionMap, nil
	}

	if field.Name.Value != utils.GraphQLAggregate || field.SelectionSet == nil {
		return nil, nil
	}

	if *aggregateFound {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "GraphQL query cannot have multiple aggregate fields, specify all functions in single aggregate field", nil, nil)
	}
	*aggregateFound = true
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
			returnFieldName := columnField.Name.Value
			columnName := strings.Join(strings.Split(columnField.Name.Value, "__"), ".")
			colArray = append(colArray, fmt.Sprintf("%s:%s", returnFieldName, columnName))
		}
		functionMap[functionField.Name.Value] = colArray
	}
	return functionMap, nil
}
func getAggregateArguments(field *ast.Directive, store utils.M) (string, string, error) {
	var operation, fieldName string
	for _, arg := range field.Arguments {
		switch arg.Name.Value {
		case "op":
			temp, err := utils.ParseGraphqlValue(arg.Value, store)
			if err != nil {
				return "", "", err
			}

			op, ok := temp.(string)
			if !ok {
				return "", "", fmt.Errorf("invalid type provided (%s) for aggregate op", reflect.TypeOf(temp))
			}

			operation = op

		case "field":
			temp, err := utils.ParseGraphqlValue(arg.Value, store)
			if err != nil {
				return "", "", err
			}

			f, ok := temp.(string)
			if !ok {
				return "", "", fmt.Errorf("invalid type provided (%s) for aggregate op", reflect.TypeOf(temp))
			}

			fieldName = f
		}
	}
	if operation == "" {
		return "", "", errors.New("need to provide `op` when using aggregations")
	}

	return operation, fieldName, nil
}

func extractQueryOp(ctx context.Context, args []*ast.Argument, store utils.M) (string, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "op":
			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return "", err
			}
			switch temp.(string) {
			case utils.All, utils.One:
				return temp.(string), nil
			default:
				return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid value provided for field (op)", nil, nil)
			}
		}
	}
	return utils.All, nil
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
			return nil, errors.New("invalid where clause provided")
		}
	}

	return utils.M{}, nil
}

func generateCacheOptions(ctx context.Context, directives []*ast.Directive, store utils.M) (*config.ReadCacheOptions, error) {
	for _, directive := range directives {
		for _, argument := range directive.Arguments {
			switch argument.Name.Value {
			case "cache":
				temp, err := utils.ParseGraphqlValue(argument.Value, store)
				if err != nil {
					return nil, err
				}

				cacheObj, ok := temp.(map[string]interface{})
				if !ok {
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for field cache in arguments expecting (object) got %v", reflect.TypeOf(temp)), err, nil)
				}
				ttlValue, ok := cacheObj["ttl"]
				if !ok {
					ttlValue = 0
				}

				ttl, ok := ttlValue.(int)
				if !ok {
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for field ttl in arguments expecting (integer) got %v", reflect.TypeOf(temp)), err, nil)
				}

				instantInvalidateObj, ok := cacheObj["instantInvalidate"]
				if !ok {
					instantInvalidateObj = false
				}

				instantInvalidate, ok := instantInvalidateObj.(bool)
				if !ok {
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for field instantInvalidate in arguments expecting (bool) got %v", reflect.TypeOf(temp)), err, nil)
				}

				return &config.ReadCacheOptions{
					TTL:               int64(ttl),
					InstantInvalidate: instantInvalidate,
				}, nil
			}
		}
	}
	return nil, nil
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

			join := make([]*model.JoinOption, 0)
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
				return nil, hasOptions, fmt.Errorf("invalid type (%s) for distinct", reflect.TypeOf(temp))
			}

			options.Distinct = &tempString
		case "debug":
			hasOptions = true // Set the flag to true

			isDebug, err := getDebugParam(ctx, []*ast.Argument{v}, store)
			if err != nil {
				return nil, false, err
			}

			options.Debug = isDebug
		}
	}
	return &options, hasOptions, nil
}

func getDebugParam(ctx context.Context, args []*ast.Argument, store utils.M) (bool, error) {
	for _, v := range args {
		if v.Name.Value == "debug" {
			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return false, err
			}

			tempBool, ok := temp.(bool)
			if !ok {
				return false, fmt.Errorf("invalid type (%s) for debug", reflect.TypeOf(temp))
			}
			return tempBool, nil
		}
	}
	return false, nil
}

func isJointTable(table string, join []*model.JoinOption) (*model.JoinOption, bool) {
	for _, j := range join {
		if j.Table == table {
			return j, true
		}
	}

	return nil, false
}
