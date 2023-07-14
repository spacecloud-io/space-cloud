package types

import "github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"

type (
	// OPAOutput describes the desired structure of opa policy output
	OPAOutput struct {
		Reason string `mapstructure:"reason"`
		Allow  bool   `mapstructure:"allow"`
		Deny   bool   `mapstructure:"deny"`
	}

	// PluginOPAParams describes the params of the `opa` plugin
	PluginOPAParams struct {
		Rego      string                `json:"rego"`
		PolicyRef *v1alpha1.ResourceRef `json:"ref"`
	}
)
