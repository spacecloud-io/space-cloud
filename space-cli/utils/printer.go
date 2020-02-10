package utils

import (
	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
)

// GetYamlObject returns the string equivalent of the git op object
func GetYamlObject(api, objType string, meta map[string]string, spec interface{}) (string, error) {
	v := model.GitOp{
		Api:  api,
		Type: objType,
		Meta: meta,
		Spec: spec,
	}

	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
