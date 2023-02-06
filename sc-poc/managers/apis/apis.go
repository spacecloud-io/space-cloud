package apis

import (
	"net/http"
	"sync"
)

// The necesary global objects to hold all registered apps
var (
	appsLock sync.RWMutex

	registeredApps apps
)

// RegisterApp marks the app which have routers
func RegisterApp(name string, priority int) {
	appsLock.Lock()
	defer appsLock.Unlock()

	registeredApps = append(registeredApps, app{name, priority})
	registeredApps.sort()
}

// Vars returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	if rv := r.Context().Value(pathParamsKey); rv != nil {
		return rv.(map[string]string)
	}

	return nil
}
