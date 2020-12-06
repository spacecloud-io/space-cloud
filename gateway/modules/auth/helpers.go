package auth

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	tmpl2 "github.com/spaceuptech/space-cloud/gateway/utils/tmpl"
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

func (m *Module) executeTemplate(ctx context.Context, rule *config.Rule, templateString string, newArgs map[string]interface{}) (interface{}, error) {
	var obj interface{}
	switch rule.Template {
	// If nothing provided default templating engine is go
	case config.TemplatingEngineGo, "":
		// Create a new template object
		t := template.New(rule.Name)
		t = t.Funcs(tmpl2.CreateGoFuncMaps(m))
		t, err := t.Parse(templateString)
		if err != nil {
			return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to parse provided template in security rule (Webhook)", err, nil))
		}
		if rule.OpFormat == "" {
			rule.OpFormat = "json"
		}
		var tempArgs interface{} = newArgs
		obj, err = tmpl2.GoTemplate(ctx, t, rule.OpFormat, newArgs["token"].(string), newArgs["auth"], tempArgs)
		if err != nil {
			return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to execute provided template in security rule (Webhook)", err, nil))
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step for security rule (Webhook) & using the default body.", rule.Template), nil)
		obj = newArgs
	}
	return obj, nil
}
