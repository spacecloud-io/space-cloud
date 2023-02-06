package file

import "github.com/spacecloud-io/space-cloud/model"

type config struct {
	Config map[string]configModule `json:"config" yaml:"config" mapstructure:"config"`
}

type configModule map[string][]*model.ResourceObject
