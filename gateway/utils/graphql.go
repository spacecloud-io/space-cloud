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
		o := map[string]interface{}{}

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
		v := value.(*ast.EnumValue).Value
		if strings.Contains(v, "__") {
			v = strings.ReplaceAll(v, "__", ".")
		}
		val, err := LoadValue(v, store)
		if err == nil {
			return val, nil
		}

		return v, nil

	case kinds.StringValue:
		v := value.(*ast.StringValue).Value
		if strings.Contains(v, "__") {
			v = strings.ReplaceAll(v, "__", ".")
		}
		val, err := LoadValue(v, store)
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
		return LoadValue("vars."+t.Name.Value, store)

	default:
		return nil, errors.New("Invalid data type `" + value.GetKind() + "` for value " + string(value.GetLoc().Source.Body)[value.GetLoc().Start:value.GetLoc().End])
	}
}

func adjustObjectKey(key string) string {
	if strings.HasPrefix(key, "_") && key != "_id" {
		key = "$" + key[1:]
	}

	key = strings.ReplaceAll(key, "__", ".")

	return key
}
