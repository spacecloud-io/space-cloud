package configman

import (
	"fmt"
)

func loadTypeDefinition(group, module, typeName string) (*TypeDefinition, error) {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	defs, p := controllerDefinitions[moduleKey(group, module)]
	if !p {
		return nil, fmt.Errorf("provided module '%s' does not exist", module)
	}

	typeDef, p := defs[typeName]
	if !p {
		return nil, fmt.Errorf("type '%s' does not exist in module '%s'", typeName, module)
	}

	return typeDef, nil
}

func unsyncLoadController(group, module string, appLoader loadApp) (interface{}, error) {
	// First check if a internal controller exists for the module
	key := moduleKey(group, module)
	appName, p := controllerApps[key]
	if p {
		// Try loading the app
		app, err := appLoader(appName)
		if err != nil {
			return nil, err
		}

		return app, nil
	}

	return nil, fmt.Errorf("no controller exists for provided module '%s'", key)
}

func moduleKey(group, module string) string {
	return fmt.Sprintf("%s/%s", group, module)
}

func loadHook(group, module string, typeDef *TypeDefinition, phase HookPhase, loadApp loadApp) (HookImpl, error) {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	// Check if hooks are defined for that phase
	if typeDef.Hooks == nil {
		return nil, nil
	}
	if _, p := typeDef.Hooks[phase]; !p {
		return nil, nil
	}

	ctrl, err := unsyncLoadController(group, module, loadApp)
	if err != nil {
		return nil, err
	}

	hookImpl, ok := ctrl.(HookImpl)
	if !ok {
		return nil, fmt.Errorf("controller '%s' doesn't implement hook functionality", moduleKey(group, module))
	}

	return hookImpl, nil
}
