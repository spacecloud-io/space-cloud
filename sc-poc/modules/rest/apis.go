package rest

import (
	"github.com/spacecloud-io/space-cloud/managers/apis"
)

// GetRoutes returns all the apis that are exposed by this app
func (a *App) GetAPIRoutes() apis.APIs {
	return a.apis
}
