package driver

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

type Driver interface {
	// GenerateTypes takes openapi spec and generates types in the respective languages
	GenerateTypes(*openapi3.T) (string, string, error)

	// GenerateAPI takes openapi spec and generates apis in the respective languages
	GenerateAPIs(*openapi3.T) (string, string, error)

	// GeneratePlugins takes list of plugins provided by SpaceCloud and generates plugins
	// in the respective languages
	GeneratePlugins([]v1alpha1.HTTPPlugin) (string, string, error)
}
