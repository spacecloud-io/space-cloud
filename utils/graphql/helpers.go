package graphql

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
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
		o := utils.M{}

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
		enumValue := value.(*ast.EnumValue)
		if strings.Contains(enumValue.Value, "__") {
			temp := strings.ReplaceAll(enumValue.Value, "__", ".")
			return utils.LoadValue(temp, store)
		}
		return enumValue.Value, nil

	case kinds.StringValue:
		stringValue := value.(*ast.StringValue)
		if strings.Contains(stringValue.Value, "__") {
			temp := strings.ReplaceAll(stringValue.Value, "__", ".")
			return utils.LoadValue(temp, store)
		}
		return stringValue.Value, nil

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

	case kinds.Variable:
		t := value.(*ast.Variable)
		return utils.LoadValue("vars."+t.Name.Value, store)

	default:
		return nil, errors.New("Invalid data type `" + value.GetKind() + "` for value " + string(value.GetLoc().Source.Body)[value.GetLoc().Start:value.GetLoc().End])
	}
}

func (graph *Module) processQueryResult(field *ast.Field, token string, store utils.M, result interface{}) (interface{}, error) {
	switch val := result.(type) {
	case []interface{}:
		array := make([]interface{}, len(val))
		for i, v := range val {
			obj := map[string]interface{}{}

			for _, sel := range field.SelectionSet.Selections {
				storeNew := shallowClone(store)
				storeNew[getFieldName(field)] = v
				storeNew["coreParentKey"] = getFieldName(field)

				f := sel.(*ast.Field)

				// if f.Name.Value == "__typename" {
				// 	continue
				// }

				output, err := graph.execGraphQLDocument(f, token, storeNew)
				if err != nil {
					return nil, err
				}

				obj[getFieldName(f)] = output
			}

			array[i] = obj
		}

		return array, nil

	case map[string]interface{}, utils.M:
		obj := map[string]interface{}{}

		for _, sel := range field.SelectionSet.Selections {
			storeNew := shallowClone(store)
			storeNew[getFieldName(field)] = val
			storeNew["coreParentKey"] = getFieldName(field)

			f := sel.(*ast.Field)
			// if f.Name.Value == "__typename" {
			// 	continue
			// }
			output, err := graph.execGraphQLDocument(f, token, storeNew)
			if err != nil {
				return nil, err
			}

			obj[getFieldName(f)] = output
		}
		return obj, nil

	default:
		log.Println("Type of val in helpers", reflect.TypeOf(val))
		return nil, errors.New("Incorrect result type")
	}
}

func adjustObjectKey(key string) string {
	if strings.HasPrefix(key, "_") {
		key = "$" + key[1:]
	}

	key = strings.ReplaceAll(key, "__", ".")

	return key
}
