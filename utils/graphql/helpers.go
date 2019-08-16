package graphql

import (
	"errors"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/spaceuptech/space-cloud/utils"
)

func shallowClone(obj m) m {
	temp := m{}
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
	col := field.Name.Value
	if len(field.Directives[0].Arguments) > 0 {
		for _, v := range field.Directives[0].Arguments {
			if v.Name.Value == "col" {
				c, ok := v.Value.GetValue().(string)
				if !ok {
					return "", errors.New("Invalid value for collection: " + string(v.Value.GetLoc().Source.Body)[v.Value.GetLoc().Start:v.Value.GetLoc().End])
				}
				col = c
			}
		}
	}

	return col, nil
}

func parseValue(value ast.Value, store m) (interface{}, error) {
	switch value.GetKind() {
	case kinds.ObjectValue:
		o := m{}

		obj := value.(*ast.ObjectValue)

		for _, v := range obj.Fields {
			temp, err := parseValue(v.Value, store)
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
			val, err := parseValue(v, store)
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

	default:
		return nil, errors.New("Invalid data type `" + value.GetKind() + "` for value " + string(value.GetLoc().Source.Body)[value.GetLoc().Start:value.GetLoc().End])
	}
}

func adjustObjectKey(key string) string {
	if strings.HasPrefix(key, "_") {
		key = "$" + key[1:]
	}

	key = strings.ReplaceAll(key, "__", ".")

	return key
}
