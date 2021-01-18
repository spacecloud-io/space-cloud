package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func checkType(ctx context.Context, dbAlias, dbType, col string, value interface{}, fieldValue *model.FieldType) (interface{}, error) {
	switch v := value.(type) {
	case int:
		// TODO: int64
		switch fieldValue.Kind {
		case model.TypeDateTime:
			return time.Unix(int64(v)/1000, 0), nil
		case model.TypeInteger, model.TypeBigInteger, model.TypeSmallInteger:
			return value, nil
		case model.TypeFloat:
			return float64(v), nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got Integer", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}

	case string:
		switch fieldValue.Kind {
		case model.TypeDateTime:
			unitTimeInRFC3339Nano, err := time.Parse(time.RFC3339Nano, v)
			if err != nil {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid datetime format recieved for field %s in collection %s - use RFC3339 fromat", fieldValue.FieldName, col), nil, nil)
			}
			return unitTimeInRFC3339Nano, nil
		case model.TypeID, model.TypeString, model.TypeTime, model.TypeDate, model.TypeUUID:
			return value, nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got String", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}

	case float32, float64:
		switch fieldValue.Kind {
		case model.TypeDateTime:
			return time.Unix(int64(v.(float64))/1000, 0), nil
		case model.TypeFloat:
			return value, nil
		case model.TypeInteger, model.TypeSmallInteger, model.TypeBigInteger:
			return int64(value.(float64)), nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got Float", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}
	case bool:
		switch fieldValue.Kind {
		case model.TypeBoolean:
			return value, nil
		default:
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s - wanted %s got Bool", fieldValue.FieldName, col, fieldValue.Kind), nil, nil)
		}

	case time.Time, *time.Time:
		return v, nil

	case map[string]interface{}:
		if fieldValue.Kind == model.TypeJSON {
			if model.DBType(dbType) == model.Mongo {
				return value, nil
			}
			data, err := json.Marshal(value)
			if err != nil {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error checking type in schema module unable to marshal data for field having type json", err, nil)
			}
			return string(data), nil
		}
		if fieldValue.Kind != model.TypeObject {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type received for field %s in collection %s", fieldValue.FieldName, col), nil, nil)
		}

		return SchemaValidator(ctx, dbAlias, dbType, col, fieldValue.NestedObject, v)

	case []interface{}:
		if !fieldValue.IsList {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid type (array) received for field %s in collection %s", fieldValue.FieldName, col), nil, nil)
		}

		arr := make([]interface{}, len(v))
		for index, value := range v {
			val, err := checkType(ctx, dbAlias, dbType, col, value, fieldValue)
			if err != nil {
				return nil, err
			}
			arr[index] = val
		}
		return arr, nil
	default:
		if !fieldValue.IsFieldTypeRequired {
			return nil, nil
		}

		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("no matching type found for field %s in collection %s", fieldValue.FieldName, col), nil, nil)
	}
}

func validateArrayOperations(ctx context.Context, dbAlias, dbType, col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection %s", col), nil, nil)
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field %s from collection %s is not defined in the schema", fieldKey, col), nil, nil)
		}

		switch t := fieldValue.(type) {
		case []interface{}:
			if schemaDocValue.IsForeign && !schemaDocValue.IsList {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for field %s in collection %s", fieldKey, col), nil, nil)
			}
			for _, value := range t {
				if _, err := checkType(ctx, dbAlias, dbType, col, value, schemaDocValue); err != nil {
					return err
				}
			}
			return nil
		case interface{}:
			if _, err := checkType(ctx, dbAlias, dbType, col, t, schemaDocValue); err != nil {
				return err
			}
		default:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for field %s in collection %s", fieldKey, col), nil, nil)
		}
	}

	return nil
}

func validateMathOperations(ctx context.Context, col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	for fieldKey, fieldValue := range v {
		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field %s from collection %s is not defined in the schema", fieldKey, col), nil, nil)
		}
		if schemaDocValue.Kind == model.TypeInteger && reflect.TypeOf(fieldValue).Kind() == reflect.Float64 {
			fieldValue = int(fieldValue.(float64))
		}
		switch fieldValue.(type) {
		case int:
			if schemaDocValue.Kind != model.TypeInteger && schemaDocValue.Kind != model.TypeFloat {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s got Integer", fieldKey, col, schemaDocValue.Kind), nil, nil)
			}
			return nil
		case float32, float64:
			if schemaDocValue.Kind != model.TypeFloat {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s got Float", fieldKey, col, schemaDocValue.Kind), nil, nil)
			}
			return nil
		default:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind), nil, nil)
		}
	}

	return nil
}

func validateDateOperations(ctx context.Context, col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	for fieldKey := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field %s from collection %s is not defined in the schema", fieldKey, col), nil, nil)
		}

		if schemaDocValue.Kind != model.TypeDateTime {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind), nil, nil)
		}
	}

	return nil
}

func validateUnsetOperation(ctx context.Context, dbType, col string, doc interface{}, schemaDoc model.Fields) error {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	// For mongo we need to check if the field to be removed is required
	if dbType == string(model.Mongo) {
		for fieldName := range v {
			columnInfo, ok := schemaDoc[strings.Split(fieldName, ".")[0]]
			if ok {
				if columnInfo.IsFieldTypeRequired {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot use $unset on field which is required/mandatory", nil, nil)
				}
			}
		}
		return nil
	}

	if dbType == string(model.Postgres) || dbType == string(model.MySQL) || dbType == string(model.SQLServer) {
		for fieldName := range v {
			columnInfo, ok := schemaDoc[strings.Split(fieldName, ".")[0]]
			if ok {
				if columnInfo.Kind == model.TypeJSON {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Cannot use $unset on field which has type (%s)", model.TypeJSON), nil, nil)
				}
			} else {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field (%s) doesn't exists in schema of (%s)", fieldName, col), nil, nil)
			}
		}
	}
	return nil
}

func validateSetOperation(ctx context.Context, dbAlias, dbType, col string, doc interface{}, SchemaDoc model.Fields) (interface{}, error) {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	newMap := map[string]interface{}{}
	for key, value := range v {
		// We could get a a key with value like `a.b`, where the user intends to set the field `b` inside object `a`. This holds true for working with json
		// types in postgres. However, no such key would be present in the schema. Hence take the top level key to validate the schema
		SchemaDocValue, ok := SchemaDoc[strings.Split(key, ".")[0]]
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field (%s) from collection (%s) is not defined in the schema", key, col), nil, nil)
		}
		// check type
		newDoc, err := checkType(ctx, dbAlias, dbType, col, value, SchemaDocValue)
		if err != nil {
			return nil, err
		}
		newMap[key] = newDoc
	}

	for fieldKey, fieldValue := range SchemaDoc {
		if fieldValue.IsUpdatedAt {
			newMap[fieldKey] = time.Now().UTC()
		}
	}

	return newMap, nil
}

func isFieldPresentInUpdate(field string, updateDoc map[string]interface{}) bool {
	for _, operatorTemp := range updateDoc {
		operator := operatorTemp.(map[string]interface{})
		if _, p := operator[field]; p {
			return true
		}
	}

	return false
}

func getCollectionSchema(doc *ast.Document, dbName, collectionName string) (model.Fields, error) {
	var isCollectionFound bool

	fieldMap := model.Fields{}
	for _, v := range doc.Definitions {
		colName := v.(*ast.ObjectDefinition).Name.Value
		if colName != collectionName {
			continue
		}

		// Mark the collection as found
		isCollectionFound = true

		for _, field := range v.(*ast.ObjectDefinition).Fields {

			if field.Type == nil {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Type not provided in graphql SDL for collection/table schema (%s) with field (%s)", collectionName, field.Name.Value), nil, nil)
			}

			fieldTypeStuct := model.FieldType{
				FieldName:  field.Name.Value,
				TypeIDSize: 50,
			}
			if len(field.Directives) > 0 {
				// Loop over all directives

				for _, directive := range field.Directives {
					switch directive.Name.Value {
					case model.DirectivePrimary:
						fieldTypeStuct.IsPrimary = true
						fieldTypeStuct.PrimaryKeyInfo = &model.TableProperties{}
						for _, argument := range directive.Arguments {
							if argument.Name.Value == "autoIncrement" {
								val, _ := utils.ParseGraphqlValue(argument.Value, nil)
								switch t := val.(type) {
								case bool:
									fieldTypeStuct.PrimaryKeyInfo.IsAutoIncrement = t
								default:
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, argument.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": argument.Name.Value})
								}

							}
							if argument.Name.Value == "order" {
								val, _ := utils.ParseGraphqlValue(argument.Value, nil)
								t, ok := val.(int)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, argument.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": argument.Name.Value})
								}
								if t == 0 {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Order cannot be zero for field (%s) directinve @(%s)", fieldTypeStuct.FieldName, directive.Name.Value), nil, map[string]interface{}{"arg": argument.Name.Value})
								}
								fieldTypeStuct.PrimaryKeyInfo.Order = t
							}
						}
					case model.DirectiveArgs:
						fieldTypeStuct.Args = new(model.FieldArgs)
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "precision":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								size, ok := val.(int)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
								fieldTypeStuct.Args.Precision = size
							case "scale":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								size, ok := val.(int)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
								fieldTypeStuct.Args.Scale = size
							}
						}
					case model.DirectiveCreatedAt:
						fieldTypeStuct.IsCreatedAt = true
					case model.DirectiveUpdatedAt:
						fieldTypeStuct.IsUpdatedAt = true
					case model.DirectiveVarcharSize:
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "value":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								size, ok := val.(int)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
								fieldTypeStuct.TypeIDSize = size
							}
						}
					case model.DirectiveDefault:
						fieldTypeStuct.IsDefault = true

						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "value":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.Default = val
							}
						}
						if fieldTypeStuct.Default == nil {
							return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Default directive must be accompanied with value field", nil, nil)
						}
					case model.DirectiveIndex, model.DirectiveUnique:
						if fieldTypeStuct.IndexInfo == nil {
							fieldTypeStuct.IndexInfo = make([]*model.TableProperties, 0)
						}
						indexInfo := &model.TableProperties{Order: model.DefaultIndexOrder, Sort: model.DefaultIndexSort, IsIndex: directive.Name.Value == model.DirectiveIndex, IsUnique: directive.Name.Value == model.DirectiveUnique}
						var ok bool
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "name", "group":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								indexInfo.Group, ok = val.(string)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
							case "order":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								indexInfo.Order, ok = val.(int)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected int", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
							case "sort":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								sort, ok := val.(string)
								if !ok || (sort != "asc" && sort != "desc") {
									if !ok {
										return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
									}
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unknow value provided for field (%s) directinve @(%s) argument (%s) got (%v) expected either (asc) or (desc)", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
								indexInfo.Sort = sort
							}
						}
						if ok {
							indexInfo.Field = field.Name.Value
							fieldTypeStuct.IndexInfo = append(fieldTypeStuct.IndexInfo, indexInfo)
						}

					case model.DirectiveLink:
						fieldTypeStuct.IsLinked = true
						fieldTypeStuct.LinkedTable = &model.TableProperties{DBType: dbName}
						kind, err := getFieldType(dbName, field.Type, &fieldTypeStuct, doc)
						if err != nil {
							return nil, err
						}
						fieldTypeStuct.LinkedTable.Table = kind

						// Load the from and to fields. If either is not available, we will throw an error.
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "table":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.Table = val.(string)
							case "from":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.From = val.(string)
							case "to":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.To = val.(string)
							case "field":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.Field = val.(string)
							case "db":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.DBType = val.(string)
							}
						}

						// Throw an error if from and to are unavailable
						if fieldTypeStuct.LinkedTable.From == "" || fieldTypeStuct.LinkedTable.To == "" {
							return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Link directive must be accompanied with (to) and (from) arguments", nil, nil)
						}

					case model.DirectiveForeign:
						fieldTypeStuct.IsForeign = true
						fieldTypeStuct.JointTable = &model.TableProperties{}
						fieldTypeStuct.JointTable.Table = strings.Split(field.Name.Value, "_")[0]
						fieldTypeStuct.JointTable.To = "id"
						fieldTypeStuct.JointTable.OnDelete = "NO ACTION"

						// Load the joint table name and field
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "table":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.Table = val.(string)

							case "field", "to":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.To = val.(string)

							case "onDelete":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.OnDelete = val.(string)
								if fieldTypeStuct.JointTable.OnDelete == "cascade" {
									fieldTypeStuct.JointTable.OnDelete = "CASCADE"
								} else {
									fieldTypeStuct.JointTable.OnDelete = "NO ACTION"
								}
							}
						}
						fieldTypeStuct.JointTable.ConstraintName = GetConstraintName(collectionName, fieldTypeStuct.FieldName)
					default:
						return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unknown directive (%s) provided for field (%s)", directive.Name.Value, fieldTypeStuct.FieldName), nil, nil)
					}
				}
			}

			kind, err := getFieldType(dbName, field.Type, &fieldTypeStuct, doc)
			if err != nil {
				return nil, err
			}
			fieldTypeStuct.Kind = kind
			switch kind {
			case model.TypeTime, model.TypeDateTime:
				if fieldTypeStuct.Args == nil {
					fieldTypeStuct.Args = new(model.FieldArgs)
				}
				if fieldTypeStuct.Args.Precision == 0 {
					fieldTypeStuct.Args.Precision = model.DefaultDateTimePrecision
				}
			case model.TypeFloat:
				if fieldTypeStuct.Args == nil {
					fieldTypeStuct.Args = new(model.FieldArgs)
				}
				if fieldTypeStuct.Args.Scale == 0 {
					fieldTypeStuct.Args.Scale = model.DefaultScale
				}
				if fieldTypeStuct.Args.Precision == 0 {
					fieldTypeStuct.Args.Precision = model.DefaultPrecision
				}
			}

			fieldMap[field.Name.Value] = &fieldTypeStuct
		}
	}

	// Throw an error if the collection wasn't found
	if !isCollectionFound {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Collection/Table (%s) not found in schema or you have provided an unknown data type (%s)", collectionName, collectionName), nil, nil)
	}
	return fieldMap, nil
}

func getFieldType(dbName string, fieldType ast.Type, fieldTypeStuct *model.FieldType, doc *ast.Document) (string, error) {
	switch fieldType.GetKind() {
	case kinds.NonNull:
		fieldTypeStuct.IsFieldTypeRequired = true
		return getFieldType(dbName, fieldType.(*ast.NonNull).Type, fieldTypeStuct, doc)
	case kinds.List:
		// Lists are not allowed for primary and foreign keys
		if fieldTypeStuct.IsPrimary || fieldTypeStuct.IsForeign {
			return "", helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Primary and foreign keys directives cannot be added on field (%s) with type lists", fieldTypeStuct.FieldName), nil, nil)
		}

		fieldTypeStuct.IsList = true
		return getFieldType(dbName, fieldType.(*ast.List).Type, fieldTypeStuct, doc)

	case kinds.Named:
		myType := fieldType.(*ast.Named).Name.Value
		switch myType {
		case model.TypeString, model.TypeEnum:
			return model.TypeString, nil
		case model.TypeID:
			return model.TypeID, nil
		case model.TypeDateTime:
			return model.TypeDateTime, nil
		case model.TypeFloat:
			return model.TypeFloat, nil
		case model.TypeInteger:
			return model.TypeInteger, nil
		case model.TypeSmallInteger:
			return model.TypeSmallInteger, nil
		case model.TypeBigInteger:
			return model.TypeBigInteger, nil
		case model.TypeBoolean:
			return model.TypeBoolean, nil
		case model.TypeJSON:
			return model.TypeJSON, nil
		case model.TypeTime:
			return model.TypeTime, nil
		case model.TypeDate:
			return model.TypeDate, nil
		case model.TypeUUID:
			return model.TypeUUID, nil

		default:
			if fieldTypeStuct.IsLinked {
				// Since the field is actually a link. We'll store the type as is. This type must correspond to a table or a primitive type
				// or else the link won't work. It's upto the user to make sure of that.
				return myType, nil
			}

			// The field is a nested type. Update the nestedObject field and return typeObject. This is a side effect.
			nestedschemaField, err := getCollectionSchema(doc, dbName, myType)
			if err != nil {
				return "", err
			}
			fieldTypeStuct.NestedObject = nestedschemaField

			return model.TypeObject, nil
		}
	default:
		return "", helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid field kind `%s` provided for field `%s`", fieldType.GetKind(), fieldTypeStuct.FieldName), nil, nil)
	}
}
