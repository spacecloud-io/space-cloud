package configman

import (
	"sync"
)

// The necesary global objects to hold all the controllers
var (
	controllerLock sync.RWMutex

	registeredOperationControllers = []string{}
	registeredConfigControllers    = []string{}
)

// RegisterOperationController adds an operation controller for the specified module
func RegisterOperationController(module string) {
	controllerLock.Lock()
	defer controllerLock.Unlock()

	registeredOperationControllers = append(registeredOperationControllers, module)
}

// RegisterConfigController adds an operation controller for the specified module
func RegisterConfigController(module string) {
	controllerLock.Lock()
	defer controllerLock.Unlock()

	registeredConfigControllers = append(registeredConfigControllers, module)
}
