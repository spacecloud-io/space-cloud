package tmpl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/getlantern/deepcopy"
	"github.com/ghodss/yaml"
	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// ExecTemplate executes a template and returns the string result
func ExecTemplate(ctx context.Context, tmpl *template.Template, object map[string]interface{}) (string, error) {
	var b strings.Builder
	if err := tmpl.Execute(&b, object); err != nil {
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to execute golang template", err, nil)
	}

	return b.String(), nil
}

// GoTemplate executes a go template
func GoTemplate(ctx context.Context, tmpl *template.Template, format, token string, claims, params interface{}) (interface{}, error) {
	// Prepare the object
	object := map[string]interface{}{"args": params, "auth": claims, "token": token}
	s, err := ExecTemplate(ctx, tmpl, object)
	if err != nil {
		return nil, err
	}

	var newParams interface{}
	switch format {
	case "string":
		return s, nil

	case "json":
		if err := json.Unmarshal([]byte(s), &newParams); err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to marhsal templated output (%s) to JSON", s), err, nil)
		}

	case "yaml", "":
		if err := yaml.Unmarshal([]byte(s), &newParams); err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to marhsal templated output (%s) to YAML", s), err, nil)
		}

	default:
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid output format (%s) provided", format), nil, nil)
	}

	return newParams, nil
}

type authModule interface {
	Encrypt(value string) (string, error)
}

// CreateGoFuncMaps creates the helper functions that can be used in go templates
func CreateGoFuncMaps(auth authModule) template.FuncMap {
	m := sprig.TxtFuncMap()
	m["hash"] = utils.HashString
	m["generateId"] = func() string { return ksuid.New().String() }
	m["marshalJSON"] = func(a interface{}) (string, error) {
		data, err := json.Marshal(a)
		return string(data), err
	}
	m["copy"] = func(a interface{}) (interface{}, error) {
		var b interface{}
		err := deepcopy.Copy(&b, a)
		return b, err
	}
	if auth != nil {
		m["encrypt"] = auth.Encrypt
	}

	return m
}
