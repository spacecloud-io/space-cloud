package tmpl

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/ghodss/yaml"
	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GoTemplate executes a go template
func GoTemplate(module, segment string, tmpl *template.Template, format, token string, claims, params interface{}) (interface{}, error) {
	// Prepare the object
	object := map[string]interface{}{"args": params, "auth": claims, "token": token}
	var b strings.Builder
	if err := tmpl.Execute(&b, object); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), "Unable to execute golang template", err, nil)
	}

	s := b.String()

	var newParams interface{}
	switch format {
	case "string":
		return s, nil

	case "json":
		if err := json.Unmarshal([]byte(s), &newParams); err != nil {
			return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Unable to marhsal templated output (%s) to JSON", s), err, nil)
		}

	case "yaml", "":
		if err := yaml.Unmarshal([]byte(s), &newParams); err != nil {
			return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Unable to marhsal templated output (%s) to YAML", s), err, nil)
		}

	default:
		return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Invalid output format (%s) provided", format), nil, nil)
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
	if auth != nil {
		m["encrypt"] = auth.Encrypt
	}

	return m
}
