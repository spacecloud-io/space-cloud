package driver

import "github.com/getkin/kin-openapi/openapi3"

type Driver interface {
	// GenerateTypes takes openapi spec and generates types in the respective languages
	GenerateTypes(*openapi3.T) (string, string, error)

	// GenerateAPI takes openapi spec and generates apis in the respective languages
	GenerateAPIs(*openapi3.T) (string, string, error)
}
