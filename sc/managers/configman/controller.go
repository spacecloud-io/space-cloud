package configman

import (
	"fmt"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/model"
)

// The necesary global objects to hold all the controllers
var (
	controllerLock        sync.RWMutex
	controllerDefinitions = map[string]model.Types{} // Key = moduleName

	openapiDoc = openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "SpaceCloud config and operation APIs",
			Description: "Specification of all the config and operation APIs exposed by SpaceCloud",
			Version:     "v0.22.0",
		},
		Components: openapi3.NewComponents(),
		Paths:      make(openapi3.Paths),
	}
)

// RegisterConfigController adds a controller for the specified module
func RegisterConfigController(module string, types model.Types) error {
	controllerLock.Lock()
	defer controllerLock.Unlock()

	// Check if controller is already present
	if _, p := controllerDefinitions[module]; p {
		return fmt.Errorf("the controller for module '%s' is already present", module)
	}

	controllerDefinitions[module] = types

	// Add the routes to the openapi docs
	addOpenAPIPath(module, types)
	return nil
}
