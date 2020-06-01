package functions

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func goTemplate(tmpl *template.Template, format string, claims, params interface{}) (interface{}, error) {
	// Prepare the object
	object := map[string]interface{}{"body": params, "auth": claims}
	var b strings.Builder
	if err := tmpl.Execute(&b, object); err != nil {
		return nil, utils.LogError("Unable to execute golang template", module, segmentCall, err)
	}

	s := b.String()

	var newParams interface{}
	switch format {
	case "string":
		return s, nil

	case "json":
		if err := json.Unmarshal([]byte(s), &newParams); err != nil {
			return nil, utils.LogError(fmt.Sprintf("Unable to marhsal templated output (%s) to JSON", s), module, segmentCall, err)
		}

	case "yaml", "":
		if err := yaml.Unmarshal([]byte(s), &newParams); err != nil {
			return nil, utils.LogError(fmt.Sprintf("Unable to marhsal templated output (%s) to YAML", s), module, segmentCall, err)
		}

	default:
		return nil, utils.LogError(fmt.Sprintf("Invalid output format (%s) provided", format), module, segmentGoTemplate, nil)
	}

	return newParams, nil
}

func (m *Module) createGoFuncMaps() template.FuncMap {
	return template.FuncMap{
		"hash":       utils.HashString,
		"add":        func(a, b int) int { return a + b },
		"generateId": func() string { return ksuid.New().String() },
		"encrypt":    m.auth.Encrypt,
	}
}
