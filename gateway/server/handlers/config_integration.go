package handlers

import (
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleGetIntegrations handles the get integration hook request
func HandleGetIntegrations(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: []interface{}{}})
	}
}
