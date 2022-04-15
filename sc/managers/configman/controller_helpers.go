package configman

import (
	"fmt"

	"github.com/spacecloud-io/space-cloud/model"
)

func loadConfigTypeDefinition(configControllerDefinitions map[string]model.ConfigTypes, module, typeName string) (*model.ConfigTypeDefinition, error) {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	defs, p := configControllerDefinitions[module]
	if !p {
		return nil, fmt.Errorf("provided module '%s' does not exist", module)
	}

	typeDef, p := defs[typeName]
	if !p {
		return nil, fmt.Errorf("type '%s' does not exist in module '%s'", typeName, module)
	}

	return typeDef, nil
}

func loadOperationTypeDefinition(operationControllerDefinitions map[string]model.OperationTypes, module, typeName string) (*model.OperationTypeDefinition, error) {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	defs, p := operationControllerDefinitions[module]
	if !p {
		return nil, fmt.Errorf("provided module '%s' does not exist", module)
	}

	typeDef, p := defs[typeName]
	if !p {
		return nil, fmt.Errorf("type '%s' does not exist in module '%s'", typeName, module)
	}

	return typeDef, nil
}
