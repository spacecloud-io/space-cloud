package graphql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
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

// GetDBAlias returns the dbAlias of the request
func (graph *Module) GetDBAlias(field *ast.Field) (string, error) {
	if len(field.Directives) == 0 {
		return "", errors.New("database / service directive not provided")
	}
	dbAlias := field.Directives[0].Name.Value

	if _, err := graph.crud.GetDBType(dbAlias); err == nil {
		return dbAlias, nil
	}

	return "", fmt.Errorf("provided db alias (%s) does not exists", dbAlias)
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

func (graph *Module) processLinkedResult(ctx context.Context, field *ast.Field, fieldStruct model.FieldType, token string, req *model.ReadRequest, store utils.M, cb model.GraphQLCallback) {
	graph.execLinkedReadRequest(ctx, field, fieldStruct.LinkedTable.DBType, fieldStruct.LinkedTable.Table, token, req,
		store, createDBCallback(func(dbAlias, col string, result interface{}, err error) {
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
			s, isSchemaPresent := graph.schema.GetSchema(dbAlias, col)

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
						graph.processLinkedResult(ctx, field, *linkedFieldSchema, token, req, store, newCB)
						return
					}
					newCB(obj, nil)
				}(loopIndex, loopValue)
			}

			wgArray.Wait()
			finalArray := newArray.GetAll()
			if !fieldStruct.IsList {
				graph.processQueryResult(ctx, field, token, store, finalArray[0], s, cb)
				return
			}
			graph.processQueryResult(ctx, field, token, store, finalArray, s, cb)
		}))
}

func (graph *Module) processQueryResult(ctx context.Context, field *ast.Field, token string, store utils.M, result interface{}, schema model.Fields, cb model.GraphQLCallback) {
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
						wg.Done()
						continue
					}

					graph.execGraphQLDocument(ctx, f, token, storeNew, schema, createCallback(func(result interface{}, err error) {
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
				wg.Done()
				continue
			}
			graph.execGraphQLDocument(ctx, f, token, storeNew, schema, createCallback(func(result interface{}, err error) {
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
