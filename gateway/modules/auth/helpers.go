package auth

import (
	"context"
	"reflect"

	"github.com/spaceuptech/helpers"

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
