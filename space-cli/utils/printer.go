package utils

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
)

// CreateSpecObject returns the string equivalent of the git op object
func CreateSpecObject(api, objType string, meta map[string]string, spec interface{}) (*model.SpecObject, error) {
	v := model.SpecObject{
		API:  api,
		Type: objType,
		Meta: meta,
		Spec: spec,
	}

	return &v, nil
}

func PrintYaml(objs []*model.SpecObject) error {
	for _, val := range objs {
		b, err := yaml.Marshal(val)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		fmt.Println("---")
	}
	return nil
}
