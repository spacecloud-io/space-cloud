package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

// M is a type for map
type M map[string]interface{}

// ParseGraphqlValue returns an interface that can be casted to string
func ParseGraphqlValue(value ast.Value, store M) (interface{}, error) {
	switch value.GetKind() {
	case kinds.ObjectValue:
		o := M{}

		obj := value.(*ast.ObjectValue)

		for _, v := range obj.Fields {
			temp, err := ParseGraphqlValue(v.Value, store)
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
			val, err := ParseGraphqlValue(v, store)
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
			return LoadValue(temp, store)
		}
		return enumValue.Value, nil

	case kinds.StringValue:
		stringValue := value.(*ast.StringValue)
		if strings.Contains(stringValue.Value, "__") {
			temp := strings.ReplaceAll(stringValue.Value, "__", ".")
			return LoadValue(temp, store)
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
