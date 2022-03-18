package configman

import (
	"fmt"
	"sync"
)

// The necesary global objects to hold all the controllers

var (
	controllerLock sync.RWMutex

	// Controller apps stores the controllers baked into the binary
	controllerApps map[string]string // moduleName -> appName

	controllerDefinitions map[string]Types // Key = moduleName
)

// AddControllerApp adds a controller for the specified module
func AddControllerApp(module, appName string, types Types) error {
	controllerLock.Lock()
	defer controllerLock.Unlock()

	// Check if controller is already present
	if _, p := controllerApps[module]; !p {
		return fmt.Errorf("the controller for module '%s' is already present", module)
	}

	controllerApps[module] = appName
	controllerDefinitions[module] = types
	return nil
}
