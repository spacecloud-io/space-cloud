package configman

import (
	"fmt"
	"sync"
)

// The necesary global objects to hold all the controllers

var (
	controllerLock sync.RWMutex

	controllerDefinitions map[string]Types // Key = moduleName
)

// AddControllerApp adds a controller for the specified module
func AddControllerApp(module string, types Types) error {
	controllerLock.Lock()
	defer controllerLock.Unlock()

	// Check if controller is already present
	if _, p := controllerDefinitions[module]; !p {
		return fmt.Errorf("the controller for module '%s' is already present", module)
	}

	controllerDefinitions[module] = types
	return nil
}
