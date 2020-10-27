package auth

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) getFields(ctx context.Context, fields interface{}, args map[string]interface{}) ([]interface{}, error) {
	switch v := fields.(type) {
	case []interface{}:
		return v, nil
	case string:
		value, err := utils.LoadValue(v, args)
		if err != nil {
			return nil, err
		}
		arr, ok := value.([]interface{})
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid value provided for field (Fields) in security rules, value obtained from args object should be an array of values", nil, map[string]interface{}{"argsObjectType": reflect.TypeOf(value)})
		}
		return arr, nil
	default:
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid value provided for field (Fields) in security rules, it can be either args object or array of string", nil, nil)
	}
}

func getMatchFields(f1, f2 interface{}) (interface{}, interface{}) {
	if f1String, ok := f1.(string); ok {
		if strings.HasPrefix(f1String, "args.find") {
			return f1, f2
		}
	}
	if f2String, ok := f2.(string); ok {
		if strings.HasPrefix(f2String, "args.find") {
			return f2, f1
		}
	}
	return f1, f2
}

func getRuleFieldForReturnWhere(field interface{}, args map[string]interface{}, stub model.ReturnWhereStub, isField1 bool) interface{} {
	fString, ok := field.(string)
	if ok {
		// Check if its in the find clause
		if strings.HasPrefix(fString, "args.find") {
			// Remove the prefix first
			fString = strings.TrimPrefix(fString, "args.find.")

			// Add the table name if its required
			if stub.PrefixColName {
				fString = stub.Col + "." + fString
			}

			return fString
		}

		// Check if we can load it from a variable
		val, err := utils.LoadValue(fString, args)
		if err == nil {
			field = val
		}
	}
	if isField1 {
		return fmt.Sprintf("'%v'", field)
	}
	return field
}
