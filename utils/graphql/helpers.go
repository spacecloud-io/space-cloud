package graphql

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/schema"
	"github.com/spaceuptech/space-cloud/utils"
)

func shallowClone(obj utils.M) utils.M {
	temp := utils.M{}
	for k, v := range obj {
		temp[k] = v
	}

	return temp
}

func getFieldName(field *ast.Field) string {
	if field.Alias != nil {
		return field.Alias.Value
	}

	return field.Name.Value
}

// GetDBType returns the dbType of the request
func GetDBType(field *ast.Field) (string, error) {
	if len(field.Directives) == 0 {
		return "", errors.New("database / service directive not provided")
	}
	dbType := field.Directives[0].Name.Value
	switch dbType {
	case "postgres", "mysql", "sqlserver":
		return "sql-" + dbType, nil

	default:
		return dbType, nil
	}
}

func getCollection(field *ast.Field) (string, error) {
	if len(field.Directives[0].Arguments) > 0 {
		for _, v := range field.Directives[0].Arguments {
			if v.Name.Value == "col" {
				col, ok := v.Value.GetValue().(string)
				if !ok {
					return "", errors.New("Invalid value for collection: " + string(v.Value.GetLoc().Source.Body)[v.Value.GetLoc().Start:v.Value.GetLoc().End])
				}
				return col, nil
			}
		}
	}
	return field.Name.Value, nil
}

// ParseValue returns an interface that can be casted to string
func ParseValue(value ast.Value, store utils.M) (interface{}, error) {
	switch value.GetKind() {
	case kinds.ObjectValue:
		o := map[string]interface{}{}

		obj := value.(*ast.ObjectValue)

		for _, v := range obj.Fields {
			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return nil, err
			}

			o[adjustObjectKey(v.Name.Value)] = temp
		}

		return o, nil

	case kinds.ListValue:
		listValue := value.(*ast.ListValue)

		array := make([]interface{}, len(listValue.Values))
		for i, v := range listValue.Values {
			val, err := ParseValue(v, store)
			if err != nil {
				return nil, err
			}

			array[i] = val
		}
		return array, nil

	case kinds.EnumValue:
		v := value.(*ast.EnumValue).Value
		if strings.Contains(v, "__") {
			v = strings.ReplaceAll(v, "__", ".")
		}
		val, err := utils.LoadValue(v, store)
		if err == nil {
			return val, nil
		}

		return v, nil

	case kinds.StringValue:
		v := value.(*ast.StringValue).Value
		if strings.Contains(v, "__") {
			v = strings.ReplaceAll(v, "__", ".")
		}
		val, err := utils.LoadValue(v, store)
		if err == nil {
			return val, nil
		}

		return v, nil

	case kinds.IntValue:
		intValue := value.(*ast.IntValue)

		// Convert string to int
		val, err := strconv.Atoi(intValue.Value)
		if err != nil {
			return nil, err
		}

		return val, nil

	case kinds.FloatValue:
		floatValue := value.(*ast.FloatValue)

		// Convert string to int
		val, err := strconv.ParseFloat(floatValue.Value, 64)
		if err != nil {
			return nil, err
		}

		return val, nil

	case kinds.BooleanValue:
		boolValue := value.(*ast.BooleanValue)
		return boolValue.Value, nil

	case kinds.Variable:
		t := value.(*ast.Variable)
		return utils.LoadValue("vars."+t.Name.Value, store)

	default:
		return nil, errors.New("Invalid data type `" + value.GetKind() + "` for value " + string(value.GetLoc().Source.Body)[value.GetLoc().Start:value.GetLoc().End])
	}
}

func (graph *Module) processLinkedResult(ctx context.Context, field *ast.Field, fieldStruct *schema.SchemaFieldType, token string, req *model.ReadRequest, store utils.M, loader *loaderMap, cb callback) {
	graph.execLinkedReadRequest(ctx, field, fieldStruct.LinkedTable.DBType, fieldStruct.LinkedTable.Table, token, req,
		store, loader, createDBCallback(func(dbType, col string, result interface{}, err error) {
			if err != nil {
				cb(nil, err)
				return
			}

			array := result.([]interface{})

			if len(array) == 0 {
				cb(nil, nil)
				return
			}

			// Check the linked table has a schema
			s, isSchemaPresent := graph.auth.Schema.GetSchema(dbType, col)

			length := len(array)
			if !fieldStruct.IsList {
				length = 1
			}

			// Create a wait group
			var wgArray sync.WaitGroup
			wgArray.Add(length)

			newArray := utils.NewArray(length)
			for loopIndex := 0; loopIndex < length; loopIndex++ {
				loopValue := array[loopIndex]

				go func(i int, v interface{}) {

					newCB := createCallback(func(result interface{}, err error) {
						defer wgArray.Done()

						if err != nil {
							cb(nil, err)
							return
						}

						newArray.Set(i, result)
					})

					obj := v.(map[string]interface{})
					if fieldStruct.LinkedTable.Field != "" {

						if !isSchemaPresent {
							// Simply return the field in the  document received
							value, p := obj[fieldStruct.LinkedTable.Field]
							if !p {
								newCB(nil, nil)
								return
							}

							newCB(value, nil)
							return
						}

						// Check if the linked field itself is a link
						linkedFieldSchema, p := s[fieldStruct.LinkedTable.Field]
						if !p || !linkedFieldSchema.IsLinked {
							// Simply return the field in the  document received
							value, p := obj[fieldStruct.LinkedTable.Field]
							if !p {
								newCB(nil, nil)
								return
							}

							// Process the value
							newCB(value, nil)
							return
						}

						// The field itself is linked. Need to query that from the database now
						linkedInfo := linkedFieldSchema.LinkedTable
						findVar, err := utils.LoadValue("args."+linkedInfo.From, map[string]interface{}{"args": obj})
						if err != nil {
							newCB(nil, nil)
							return
						}
						req := &model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{linkedInfo.To: findVar}}
						graph.processLinkedResult(ctx, field, linkedFieldSchema, token, req, store, loader, newCB)
						return
					}
					newCB(obj, nil)
					return
				}(loopIndex, loopValue)
			}

			wgArray.Wait()
			finalArray := newArray.GetAll()
			if !fieldStruct.IsList {
				graph.processQueryResult(ctx, field, token, store, finalArray[0], loader, s, cb)
				return
			}
			graph.processQueryResult(ctx, field, token, store, finalArray, loader, s, cb)
		}))
}

func (graph *Module) processQueryResult(ctx context.Context, field *ast.Field, token string, store utils.M, result interface{}, loader *loaderMap, schema schema.SchemaFields, cb callback) {
	addFieldPath(store, getFieldName(field))

	switch val := result.(type) {
	case []interface{}:
		array := utils.NewArray(len(val))

		// Create a wait group
		var wgArray sync.WaitGroup
		wgArray.Add(len(val))

		for loopIndex, loopValue := range val {
			go func(i int, v interface{}) {
				defer wgArray.Done()

				if field.SelectionSet == nil {
					array.Set(i, v)
					return
				}

				obj := utils.NewObject()

				// Create a wait group
				var wg sync.WaitGroup
				wg.Add(len(field.SelectionSet.Selections))

				for _, sel := range field.SelectionSet.Selections {
					storeNew := shallowClone(store)
					storeNew[getFieldName(field)] = v
					storeNew["coreParentKey"] = getFieldName(field)

					f := sel.(*ast.Field)

					if f.Name.Value == "__typename" {
						obj.Set(f.Name.Value, strings.Title(field.Name.Value))
						continue
					}

					graph.execGraphQLDocument(ctx, f, token, storeNew, loader, schema, createCallback(func(result interface{}, err error) {
						defer wg.Done()

						if err != nil {
							cb(nil, err)
							return
						}

						obj.Set(getFieldName(f), result)
					}))
				}

				wg.Wait()
				array.Set(i, obj.GetAll())
			}(loopIndex, loopValue)
		}

		wgArray.Wait()
		cb(array.GetAll(), nil)
		return

	case map[string]interface{}, utils.M:
		obj := utils.NewObject()

		if field.SelectionSet == nil {
			cb(val, nil)
			return
		}

		// Create a wait group
		var wg sync.WaitGroup
		wg.Add(len(field.SelectionSet.Selections))

		for _, sel := range field.SelectionSet.Selections {
			storeNew := shallowClone(store)
			storeNew[getFieldName(field)] = val
			storeNew["coreParentKey"] = getFieldName(field)

			f := sel.(*ast.Field)
			if f.Name.Value == "__typename" {
				obj.Set(f.Name.Value, strings.Title(field.Name.Value))
				continue
			}
			graph.execGraphQLDocument(ctx, f, token, storeNew, loader, schema, createCallback(func(result interface{}, err error) {
				defer wg.Done()

				if err != nil {
					cb(nil, err)
					return
				}

				obj.Set(getFieldName(f), result)
			}))
		}
		wg.Wait()
		cb(obj.GetAll(), nil)
		return

	default:
		cb(result, nil)
		return
	}
}

func addFieldPath(store utils.M, field string) {
	if _, p := store["path"]; !p {
		store["path"] = ""
	}

	store["path"] = store["path"].(string) + "." + field
}

func adjustObjectKey(key string) string {
	if strings.HasPrefix(key, "_") && key != "_id" {
		key = "$" + key[1:]
	}

	key = strings.ReplaceAll(key, "__", ".")

	return key
}
